package pkg

import (
	"HCPlatform/code/protos/term"
	"context"
	"fmt"
	"time"
)

var (
	shellIdsMap map[uint64]*Terminal
)

type TermnialService struct {
}

func (s *TermnialService) NewTerminal(ctx context.Context, request *term.TerminalRequest) (*term.TerminalResponse, error) {
	//TODO implement me
	shellId := uint64(time.Now().Unix())
	//TODO 根据当前系统检查
	shell, err := NewPowerShell()
	resp := new(term.TerminalResponse)
	if err != nil {
		resp.Result = fmt.Sprint("Error happened\n%s", err.Error())
	}
	shellIdsMap[shellId] = shell
	resp.ShellId = shellId
	return resp, nil
}

func (s *TermnialService) CloseTerminal(ctx context.Context, request *term.TerminalRequest) (*term.TerminalResponse, error) {
	//TODO implement me
	resp := new(term.TerminalResponse)

	shellId := request.ShellId
	shell := shellIdsMap[shellId]
	shell.Exit()
	resp.Result = fmt.Sprintf("Exit Shell: %d", shellId)
	return resp, nil
}

func (s *TermnialService) ExecCommand(ctx context.Context, request *term.TerminalRequest) (*term.TerminalResponse, error) {
	//TODO implement me
	resp := new(term.TerminalResponse)
	shellId := request.ShellId
	shell := shellIdsMap[shellId]
	sout, serr, err := shell.Execute(request.Command)
	if err != nil {
		resp.Result = fmt.Sprintf("Error happened\n%s", serr)
	}
	resp.Result = sout
	return resp, nil
}

func (s *TermnialService) mustEmbedUnimplementedTerminalServer(ctx context.Context, request *term.TerminalRequest) (*term.TerminalResponse, error) {
	//TODO implement me
	resp := new(term.TerminalResponse)
	return resp, nil
}
