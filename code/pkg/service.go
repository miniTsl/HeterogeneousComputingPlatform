package pkg

import (
	"HCPlatform/code/internal"
	"container/list"
	"fmt"
	"net"

	"log"
	"runtime"
)

type Service struct {
	ip       string
	port     int
	listener *net.Listener
	handlers *list.List
}

func NewNetListener(ip string, port int) net.Listener {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Fatal(err)
	}
	return l
}

func NewService(ip string, port int) *Service {
	service := Service{ip: ip, port: port, handlers: list.New()}
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", service.ip, service.port))

	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	service.listener = &l
	return &service
}

func RunService(service *Service) {

	for {
		// Wait for a connection.
		conn, err := (*service.listener).Accept()
		if err != nil {
			log.Fatal(err)
		}
		service.handlers.PushBack(processConn(conn))
	}

}

func StopService(service *Service) {

}

type Handler struct {
	id         string
	isExited   bool
	remoteConn net.Conn
	shell      *internal.Terminal
}

func processConn(conn net.Conn) *Handler {
	//_shell, err := shell.NewPowerShell()
	var _shell *internal.Terminal
	var err error
	sysType := runtime.GOOS
	if sysType == "linux" {
		// LINUX系统
		_shell, err = internal.NewBourneAgainShell()
		if err != nil {
			fmt.Printf("fail to create bash\n")
			return nil
		}

	} else if sysType == "windows" {
		_shell, err = internal.NewPowerShell()
		if err != nil {
			fmt.Printf("fail to create powershell\n")
			return nil
		}
	} else if sysType == "darwin" {
		_shell, err = internal.NewZShell()
		if err != nil {
			fmt.Printf("fail to create zsh\n")
			return nil
		}
	} else {
		return nil
	}
	handler := Handler{isExited: false, remoteConn: conn, shell: _shell}
	// 使用go关键字实现goroutines协程执行函数
	go handler.loop()
	return &handler
}

func (h *Handler) loop() {
	for {
		if h.isExited {
			break
		}

		request := RecvRequest(h.remoteConn)
		switch request.Pyload.(type) {
		case *(Request_ShellRequest):
			cmd := request.GetShellRequest().GetCommand()
			sout, serr, err := h.shell.Execute(cmd)
			if err == nil {
				SendResponse(h.remoteConn, Response_success, sout)
			} else {
				SendResponse(h.remoteConn, Response_error, serr)
			}
			break
		case *(Request_FileRequest):
			break
		}
	}
	h.shell.Exit()
}
