package pkg

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	serverIP   string
	serverPort int
	conn       net.Conn
}

func NewConn(ip string, port int) *Client {

	client := &Client{serverIP: ip, serverPort: port}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Printf("conn server failed, err:%v\n", err)
		return nil
	}
	client.conn = conn
	return client
}

func (c *Client) loop() {
	input := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("$:")
		s, _ := input.ReadString('\n')
		s = strings.TrimSpace(s)
		if strings.ToUpper(s) == "Q" {
			return
		}
		SendShellRequest(c.conn, s)
		response := RecvResponse(c.conn)
		fmt.Printf("%s\n", response.GetMsg())
		// 从服务端接收回复消息
		//var buf [4096]byte
		//n, err := c.conn.Read(buf[:])
		//if err != nil {
		//	fmt.Printf("read failed:%v\n", err)
		//	return
		//}
		//fmt.Printf("%v\n", string(buf[:n]))
	}
}
