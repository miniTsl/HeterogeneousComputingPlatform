package term

import (
	"HCPlatform/code/pkg"
	"context"
	"fmt"
	"time"
)

var (
	// this variable is for client to store shellId-Terminal point
	shellIdsMap map[uint64]*pkg.Terminal
)

type TermnialService struct {
}

func (s *TermnialService) NewTerminal(ctx context.Context, request *TerminalRequest) (*TerminalResponse, error) {

	shellId := uint64(time.Now().Unix())
	//TODO We should launch shell according to current OS.
	shell, err := pkg.NewPowerShell()
	resp := new(TerminalResponse)
	if err != nil {
		resp.Result = fmt.Sprint("Error happened\n%s", err.Error())
	}
	shellIdsMap[shellId] = shell
	resp.ShellId = shellId
	return resp, nil
}

func (s *TermnialService) CloseTerminal(ctx context.Context, request *TerminalRequest) (*TerminalResponse, error) {
	resp := new(TerminalResponse)
	shellId := request.ShellId
	shell := shellIdsMap[shellId]
	shell.Exit()
	resp.Result = fmt.Sprintf("Exit Shell: %d", shellId)
	return resp, nil
}

func (s *TermnialService) ExecCommand(ctx context.Context, request *TerminalRequest) (*TerminalResponse, error) {
	resp := new(TerminalResponse)
	shellId := request.ShellId
	shell := shellIdsMap[shellId]
	sout, serr, err := shell.Execute(request.Command)
	if err != nil {
		resp.Result = fmt.Sprintf("Error happened\n%s", serr)
	}
	resp.Result = sout
	return resp, nil
}

func (s *TermnialService) mustEmbedUnimplementedTerminalServer(ctx context.Context, request *TerminalRequest) (*TerminalResponse, error) {
	resp := new(TerminalResponse)
	return resp, nil
}

func NewTerminalRequest(deviceId string, shellId string, command string) {

}

func NewTerminalRequestAllArgs(dstIP string, srcIP string, deviceId string, shellId string, command string) {

}

// This function is for client to fast call
func FastNewTerminal(deviceName string) uint64 {
	return 0
}

// This function is for client to fast call
func FastCloseTerminal(deviceName string, deviceId uint64) {

}

// This function is for client to fast call
func FastExecCommand(deviceName string, command string, shellId uint64) (string, error) {

	return "", nil
}
