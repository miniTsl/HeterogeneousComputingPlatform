package pkg

import (
	"crypto/rand"
	// Package rand implements a cryptographically secure random number generator.
	"encoding/hex"
	"fmt"
	"github.com/juju/errors"
	"io"
	// Package io provides basic interfaces to I/O primitives. 
	"os/exec"
	// Package exec runs external commands. It wraps os.StartProcess to make it easier to remap stdin and stdout, connect I/O with pipes, and do other adjustments.
	"strings"
	"sync"
)

type Terminal struct {
	shellName string
	newline   string
	fullFmt   string
	stdin     io.Writer
	stdout    io.Reader
	stderr    io.Reader
	handle    *exec.Cmd
}

func NewShell(cmd string, args ...string) (*exec.Cmd, io.Writer, io.Reader, io.Reader, error) {
	command := exec.Command(cmd, args...)	
	// Command returns the Cmd struct to execute the `named program` with the given arguments.

	stdin, err := command.StdinPipe()
// StdinPipe returns a pipe that will be connected to the command's standard input when the command starts. 
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, nil, nil, nil, errors.Annotate(err, "Could not get hold of the PowerShell's stdin stream")
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, nil, nil, nil, errors.Annotate(err, "Could not get hold of the PowerShell's stdout stream")
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, nil, nil, nil, errors.Annotate(err, "Could not get hold of the PowerShell's stderr stream")
	}

	err = command.Start()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, nil, nil, nil, errors.Annotate(err, "Could not spawn PowerShell process")
	}
	return command, stdin, stdout, stderr, nil
}

func NewPowerShell() (*Terminal, error) {
	handle, stdin, stdout, stderr, err := NewShell("powershell.exe", "-NoExit", "-Command", "-")
	if err != nil {
		return nil, err
	}
	t := Terminal{shellName: "powershell", newline: "\r\n", handle: handle, stdin: stdin, stdout: stdout, stderr: stderr, fullFmt: "%s; echo '%s'; [Console]::Error.WriteLine('%s')%s"}
	// https://learn.microsoft.com/zh-cn/dotnet/api/system.console.error?view=net-7.0
	return &t, nil
}

func NewZShell() (*Terminal, error) {

	handle, stdin, stdout, stderr, err := NewShell("/bin/zsh", "-i", "-s")
	if err != nil {
		return nil, err
	}
	t := Terminal{shellName: "zsh", newline: "\n", handle: handle, stdin: stdin, stdout: stdout, stderr: stderr, fullFmt: "%s; echo '%s'; echo '%s'>&2%s"}
	// >&2：重定向到标准错误输出
	return &t, nil
}

func NewBourneAgainShell() (*Terminal, error) {
	handle, stdin, stdout, stderr, err := NewShell("/bin/bash", "-i", "-s")
	if err != nil {
		return nil, err
	}
	t := Terminal{shellName: "bash", newline: "\n", handle: handle, stdin: stdin, stdout: stdout, stderr: stderr, fullFmt: "%s; echo '%s'; echo '%s'>&2%s"}

	return &t, nil
}

func (s *Terminal) Execute(cmd string) (string, string, error) {
	if s.handle == nil {
		return "", "", errors.Annotate(errors.New(cmd), "Cannot execute commands on closed shells.")
	}
	// Annotate is used to add extra context to an existing error. 

	// 创建分隔符
	outBoundary := createBoundary()
	errBoundary := createBoundary()

	
	// wrap the command in special markers so we know when to stop reading from the pipes
	// full执行的效果是：先执行cmd，然后在标准输出中继续输出outBoundary，然后在标准错误中输出errBoundary
	// todo 适配zsh
	full := fmt.Sprintf(s.fullFmt, cmd, outBoundary, errBoundary, s.newline)
	// 将wrap之后的总命令写入shell
	_, err := s.stdin.Write([]byte(full))

	if err != nil {
		return "", "", errors.Annotate(errors.Annotate(err, full), "Could not send command")
	}
	// read stdout and stderr
	sout := ""
	serr := ""

	waiter := &sync.WaitGroup{}
	// A WaitGroup waits for a collection of goroutines to finish. The main goroutine calls Add to set the number of goroutines to wait for. Then each of the goroutines runs and calls Done when finished. At the same time, Wait can be used to block until all goroutines have finished.
	waiter.Add(2)	// 因为我们要读取标准输出和标准错误

	go streamReader(s.stdout, outBoundary, &sout, waiter, s.newline)
	go streamReader(s.stderr, errBoundary, &serr, waiter, s.newline)

	waiter.Wait()

	if len(serr) > 0 {
		return sout, serr, errors.Annotate(errors.New(cmd), serr)
	}
	return sout, serr, nil
}

func (s *Terminal) Exit() {
	s.stdin.Write([]byte("exit" + s.newline))
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

func streamReader(stream io.Reader, boundary string, buffer *string, signal *sync.WaitGroup, newline string) error {
	// read all output until we have found our boundary token
	output := ""
	bufsize := 64
	marker := boundary + newline

	for {
		buf := make([]byte, bufsize)
		read, err := stream.Read(buf)	// 返回读取的字节个数以及错误类型
		if err != nil {
			fmt.Printf("err\n")
			return err
		}

		output = output + string(buf[:read])
		if strings.HasSuffix(output, marker) {	// 检查是否有后缀marker，如果有的话说明读取完毕
			break
		}
	}

	*buffer = strings.TrimSuffix(output, marker)	// 去掉后缀
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
