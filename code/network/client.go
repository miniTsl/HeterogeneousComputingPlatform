package network

import (
	"bufio"
	// Package bufio implements buffered I/O.
	"fmt"
	// Package fmt implements formatted I/O with functions analogous to C's printf and scanf. The format 'verbs' are derived from C's but are simpler.
	"net"
	// Package net provides a portable interface for network I/O, including TCP/IP, UDP, domain name resolution, and Unix domain sockets.
	"os"
	// Package os provides a platform-independent interface to operating system functionality.
	"strings"
)

type Client struct {
	serverIP   string
	serverPort int
	conn       net.Conn
	//  tcp连接的句柄，Conn is a generic stream-oriented network connection. 一个interface，支持read和write等功能
}

func NewConn(ip string, port int) *Client {
	// 创建客户端对象
	client := &Client{serverIP: ip, serverPort: port}
	// 建立指向特定地址的网络连接
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Printf("conn server failed, err:%v\n", err)
		return nil
	}
	client.conn = conn
	// 客户端对象开始接受用户输入并发送
	client.loop()
	return client
}

// 自动成为成员函数？？？？
func (c *Client) loop() {
	input := bufio.NewReader(os.Stdin) // NewReader returns a new Reader whose buffer has the default size.
	for {
		fmt.Printf("$:")
		s, _ := input.ReadString('\n') // 命令以\n结尾
		s = strings.TrimSpace(s)
		// TrimSpace returns a slice of the string s, with all leading and trailing white space removed
		// q/Q退出
		if strings.ToUpper(s) == "Q" {
			return
		}
		// 调用发送函数发送数据
		// string->proto格式->[]byte
		SendRequest(c.conn, s)
		// 接收返回的数据
		// []byte->proto格式
		response := RecvResponse(c.conn)
		fmt.Printf("%s\n", response.GetMsg())
	}
}
