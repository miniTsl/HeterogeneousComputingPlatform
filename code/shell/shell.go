package shell

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

type Terminal struct {
	shellName string
	newline   string
	stdin     io.Writer
	stdout    io.Reader
	stderr    io.Reader
	handle    *exec.Cmd
}

func NewShell(cmd string, args ...string) (*exec.Cmd, io.Writer, io.Reader, io.Reader, error) {
	command := exec.Command(cmd, args...)

	stdin, err := command.StdinPipe()
	if err != nil {
		return nil, nil, nil, nil, errors.Annotate(err, "Could not get hold of the PowerShell's stdin stream")
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, nil, nil, nil, errors.Annotate(err, "Could not get hold of the PowerShell's stdout stream")
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		return nil, nil, nil, nil, errors.Annotate(err, "Could not get hold of the PowerShell's stderr stream")
	}

	err = command.Start()
	if err != nil {
		return nil, nil, nil, nil, errors.Annotate(err, "Could not spawn PowerShell process")
	}
	return command, stdin, stdout, stderr, nil
}

func NewPowerShell() (*Terminal, error) {
	handle, stdin, stdout, stderr, err := NewShell("powershell.exe", "-NoExit", "-Command", "-")
	if err != nil {
		return nil, err
	}
	t := Terminal{shellName: "powershell", newline: "\r\n", handle: handle, stdin: stdin, stdout: stdout, stderr: stderr}

	return &t, nil
}

func (s *Terminal) Execute(cmd string) (string, string, error) {

	if s.handle == nil {
		return "", "", errors.Annotate(errors.New(cmd), "Cannot execute commands on closed shells.")
	}

	outBoundary := createBoundary()
	errBoundary := createBoundary()

	// wrap the command in special markers so we know when to stop reading from the pipes
	full := fmt.Sprintf("%s; echo '%s'; [Console]::Error.WriteLine('%s')%s", cmd, outBoundary, errBoundary, s.newline)

	_, err := s.stdin.Write([]byte(full))
	if err != nil {
		return "", "", errors.Annotate(errors.Annotate(err, cmd), "Could not send PowerShell command")
	}

	// read stdout and stderr
	sout := ""
	serr := ""

	waiter := &sync.WaitGroup{}
	waiter.Add(2)

	go streamReader(s.stdout, outBoundary, &sout, waiter, s.newline)
	go streamReader(s.stderr, errBoundary, &serr, waiter, s.newline)

	waiter.Wait()

	if len(serr) > 0 {
		return sout, serr, errors.Annotate(errors.New(cmd), serr)
	}

	return sout, serr, nil
}

func (s *Terminal) Exit() {

}

func streamReader(stream io.Reader, boundary string, buffer *string, signal *sync.WaitGroup, newline string) error {
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

func CreateRandomString(bytes int) string {
	c := bytes
	b := make([]byte, c)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}

func createBoundary() string {
	return "$gorilla" + CreateRandomString(12) + "$"
}
