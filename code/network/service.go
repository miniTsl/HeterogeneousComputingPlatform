package network

import (
	"HCPlatform/code/shell"
	"bufio"
	"container/list"
	"fmt"
	"log"
	"net"
)

type Service struct {
	ip       string
	port     int
	listener *net.Listener
	handlers *list.List
}

func RunService(ip string, port int) *Service {
	service := Service{ip: ip, port: port, handlers: list.New()}
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", service.ip, service.port))

	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	service.listener = &l
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		service.handlers.PushBack(processConn(conn))
	}
	return &service
}

type Handler struct {
	id         string
	isExited   bool
	remoteConn net.Conn
	shell      *shell.Terminal
}

func processConn(conn net.Conn) *Handler {
	shell, _ := shell.NewPowerShell()
	handler := Handler{isExited: false, remoteConn: conn, shell: shell}
	// 使用go关键字实现goroutines协程执行函数
	go handler.loop()
	return &handler
}

func (h *Handler) loop() {
	for {
		if h.isExited {
			break
		}
		reader := bufio.NewReader(h.remoteConn)
		var buf [1024]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Printf("read from conn failed, err:%v\n", err)
			break
		}

		recv := string(buf[:n])
		fmt.Printf("exec: %s\n", recv)
		sout, serr, err := h.shell.Execute(recv)
		if err != nil {
			_, err = h.remoteConn.Write([]byte(serr))
			fmt.Printf("err: %s\n", serr)
		} else {
			_, err = h.remoteConn.Write([]byte(sout))
			fmt.Printf("out: %s\n", sout)
		}

		if err != nil {
			fmt.Printf("write from conn failed, err:%v\n", err)
			break
		}
	}
	h.shell.Exit()
}
