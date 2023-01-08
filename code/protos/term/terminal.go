package term

import (
	"HCPlatform/code/pkg"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

var (
	// this variable is for client to store shellId-Terminal point
	shellIdsMap map[uint64]*pkg.Terminal
)

func init() {
	shellIdsMap = make(map[uint64]*pkg.Terminal)
}

type TermnialService struct {
}

func (s *TermnialService) NewTerminal(ctx context.Context, request *TerminalRequest) (*TerminalResponse, error) {
	//TODO 这个部分需要区分是要转发请求还是执行请求，目前仅执行,不转发,Server端不开Terminal Service
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
	resp.ShellId = shellId
	resp.Result = fmt.Sprintf("Exit Shell: %d", shellId)
	return resp, nil
}

func (s *TermnialService) ExecCommand(ctx context.Context, request *TerminalRequest) (*TerminalResponse, error) {
	resp := new(TerminalResponse)
	shellId := request.ShellId
	shell := shellIdsMap[shellId]

	sout, serr, err := shell.Execute(request.Command)
	//log.Info(fmt.Sprintf("Result:%s", request.Command, sout))
	resp.Result = sout
	if err != nil {
		log.Error(err)
		resp.Result = fmt.Sprintf("%s\nBut happen error:%s", resp.Result, serr)
	}
	return resp, nil
}

func (s *TermnialService) mustEmbedUnimplementedTerminalServer() {

}

func NewTerminalRequest(deviceId uint64, shellId uint64, command string) TerminalRequest {
	request := TerminalRequest{
		DstIP:    "",
		SrcIP:    "",
		ShellId:  shellId,
		Command:  command,
		DeviceId: deviceId,
	}
	return request
}

func NewTerminalRequestAllArgs(dstIP string, srcIP string, deviceId uint64, shellId uint64, command string) TerminalRequest {
	request := TerminalRequest{
		DstIP:    dstIP,
		SrcIP:    srcIP,
		ShellId:  shellId,
		Command:  command,
		DeviceId: deviceId,
	}
	return request
}

// This function is for client to fast call
func FastNewTerminal(serverIP string, serverPort int, deviceId uint64, shellId uint64) (uint64, error) {
	req := NewTerminalRequest(deviceId, shellId, "")
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", serverIP, serverPort), grpc.WithInsecure())
	if err != nil {
		log.Error(err)
		return 0, err
	}
	defer conn.Close()
	c := NewTerminalClient(conn)
	res, err := c.NewTerminal(context.Background(), &req)
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}
	err = conn.Close()
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}
	return res.ShellId, nil
}

// This function is for client to fast call
func FastCloseTerminal(serverIP string, serverPort int, deviceId uint64, shellId uint64) (uint64, error) {
	req := NewTerminalRequest(deviceId, shellId, "")
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", serverIP, serverPort), grpc.WithInsecure())
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}
	defer conn.Close()
	c := NewTerminalClient(conn)
	res, err := c.CloseTerminal(context.Background(), &req)
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}
	err = conn.Close()
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}
	return res.ShellId, nil
}

// This function is for client to fast call
func FastExecCommand(serverIP string, serverPort int, deviceId uint64, shellId uint64, command string) (string, error) {

	req := NewTerminalRequest(deviceId, shellId, command)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", serverIP, serverPort), grpc.WithInsecure())
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	defer conn.Close()
	c := NewTerminalClient(conn)
	res, err := c.ExecCommand(context.Background(), &req)

	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	err = conn.Close()
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	return res.Result, nil
}
