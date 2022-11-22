package network

import (
	"HCPlatform/code/pkg"
	"container/list"
	"fmt"
	"log"
	"net"
	"runtime"
)

type Service struct {
	ip       string
	port     int
	listener *net.Listener // A Listener is a generic network listener for stream-oriented protocols.
	handlers *list.List    // 链表中存放多个shell子进程(用Handler结构体表征)
	// List represents a doubly linked list双向链表. The zero value for List is an empty list ready to use.
}

func RunService(ip string, port int) *Service {
	service := Service{ip: ip, port: port, handlers: list.New()}
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", service.ip, service.port))

	if err != nil {
		log.Fatal(err)
	}
	// 在 defer 归属的函数即将返回时，将延迟处理的语句按 defer 的逆序进行执行，即close在RunService返回时执行
	defer l.Close()
	service.listener = &l
	for {
		// Wait for a connection.
		conn, err := l.Accept() // Accept waits for and returns the next connection to the listener.
		if err != nil {
			log.Fatal(err)
		}
		// 加入新建立的连接
		service.handlers.PushBack(processConn(conn))
	}
	return &service
}

type Handler struct {
	id         string
	isExited   bool
	remoteConn net.Conn
	shell      *pkg.Terminal
}

func processConn(conn net.Conn) *Handler {

	var _shell *pkg.Terminal
	var err error
	sysType := runtime.GOOS
	/* go tool dist list
	aix/ppc64
	android/386
	android/amd64
	android/arm
	android/arm64
	darwin/amd64
	darwin/arm64
	dragonfly/amd64
	freebsd/386
	freebsd/amd64
	freebsd/arm
	freebsd/arm64
	illumos/amd64
	ios/amd64
	ios/arm64
	js/wasm
	linux/386
	linux/amd64
	linux/arm
	linux/arm64
	linux/loong64
	linux/mips
	linux/mips64
	linux/mips64le
	linux/mipsle
	linux/ppc64
	linux/ppc64le
	linux/riscv64
	linux/s390x
	netbsd/386
	netbsd/amd64
	netbsd/arm
	netbsd/arm64
	openbsd/386
	openbsd/amd64
	openbsd/arm
	openbsd/arm64
	openbsd/mips64
	plan9/386
	plan9/amd64
	plan9/arm
	solaris/amd64
	windows/386
	windows/amd64
	windows/arm
	windows/arm64
	*/
	if sysType == "linux" {
		// LINUX系统
		_shell, err = pkg.NewBourneAgainShell()
		if err != nil {
			fmt.Printf("fail to create bash\n")
			return nil
		}

	} else if sysType == "windows" {
		_shell, err = pkg.NewPowerShell()
		if err != nil {
			fmt.Printf("fail to create powershell\n")
			return nil
		}
	} else if sysType == "darwin" {
		_shell, err = pkg.NewZShell()
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
		// 直到连接中断才退出
		if h.isExited {
			break
		}
		// 从连接Conn中获取请求
		request := RecvRequest(h.remoteConn)
		switch request.Pyload.(type) {
		case *(Request_CommandRequest):	// 命令类型请求
			// proto类型->string
			cmd := request.GetCommandRequest().GetCommand()

			// 执行命令
			sout, serr, err := h.shell.Execute(cmd)
			
			// 将执行结果返回client
			// string->proto格式->[]byte
			if err == nil {
				SendResponse(h.remoteConn, Response_success, sout)
			} else {
				full_out := fmt.Sprintf("%s\nerr:%s\n", sout, serr)
				SendResponse(h.remoteConn, Response_error, full_out)

				// SendResponse(h.remoteConn, Response_error, serr
			}
			break
		case *(Request_FileRequest):	// TODO 文件类型请求
			break
		}
	}
	h.shell.Exit()
}