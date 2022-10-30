package network

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/juju/errors"
	"io"
	"os/exec"
	"strings"
	"sync"
)

const newline = "\r\n"

type PowerShell struct {
	handle *exec.Cmd
	stdin  io.Writer
	stdout io.Reader
	stderr io.Reader
}

func New() (PowerShell, error) {
	return NewLocalPowerShell("powershell.exe", "-NoExit", "-Command", "-")
}

func NewLocalPowerShell(cmd string, args ...string) (PowerShell, error) {
	command := exec.Command(cmd, args...)
	stdin, err := command.StdinPipe()
	if err != nil {
		return PowerShell{}, errors.Annotate(err, "Could not get hold of the PowerShell's stdin stream")
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		return PowerShell{}, errors.Annotate(err, "Could not get hold of the PowerShell's stdout stream")
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		return PowerShell{}, errors.Annotate(err, "Could not get hold of the PowerShell's stderr stream")
	}

	err = command.Start()
	if err != nil {
		return PowerShell{}, errors.Annotate(err, "Could not spawn PowerShell process")
	}

	return PowerShell{command, stdin, stdout, stderr}, nil
}

func (s *PowerShell) Exec(cmd string) (string, error) {
	if s.handle == nil {
		return "", errors.Annotate(errors.New(cmd), "Cannot execute commands on closed shells.")
	}

	outBoundary := createBoundary()
	errBoundary := createBoundary()

	// wrap the command in special markers so we know when to stop reading from the pipes
	full := fmt.Sprintf("%s; echo '%s'; [Console]::Error.WriteLine('%s')%s", cmd, outBoundary, errBoundary, newline)

	_, err := s.stdin.Write([]byte(full))
	if err != nil {
		return "", errors.Annotate(errors.Annotate(err, cmd), "Could not send PowerShell command")
	}

	// read stdout and stderr
	sout := ""
	serr := ""

	waiter := &sync.WaitGroup{}
	waiter.Add(2)

	go streamReader(s.stdout, outBoundary, &sout, waiter)
	go streamReader(s.stderr, errBoundary, &serr, waiter)

	waiter.Wait()

	if len(serr) > 0 {
		return serr, errors.Annotate(errors.New(cmd), serr)
	}

	return sout, nil
}

func (s *PowerShell) Exit() {
	s.stdin.Write([]byte("exit" + newline))

	// if it's possible to close stdin, do so (some backends, like the local one,
	// do support it)
	closer, ok := s.stdin.(io.Closer)
	if ok {
		closer.Close()
	}

	s.handle.Wait()

	s.handle = nil
	s.stdin = nil
	s.stdout = nil
	s.stderr = nil
}

func streamReader(stream io.Reader, boundary string, buffer *string, signal *sync.WaitGroup) error {
	// read all output until we have found our boundary token
	output := ""
	bufsize := 64
	marker := boundary + newline

	for {
		buf := make([]byte, bufsize)
		read, err := stream.Read(buf)
		if err != nil {
			return err
		}

		output = output + string(buf[:read])

		if strings.HasSuffix(output, marker) {
			break
		}
	}

	*buffer = strings.TrimSuffix(output, marker)
	signal.Done()

	return nil
}

func createBoundary() string {
	return "$gorilla" + CreateRandomString(12) + "$"
}

func CreateRandomString(bytes int) string {
	c := bytes
	b := make([]byte, c)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}
