package network

import (
	"bufio"
	"fmt"
	"google.golang.org/protobuf/proto"
	"net"
)
// client端发送请求
func SendRequest(conn net.Conn, cmd string) (string, string, error) {
	// 将要发送的数据包装为protobuf要求的request类型结构体
	msg := Request{
		From:  &ConnAddress{},
		To:    &ConnAddress{},
		XType: Request_command,
		Pyload: &Request_CommandRequest{
			CommandRequest: &CommandRequest{
				ShellId: 0,	// TODO 定制id
				Command: cmd,
			},
		},
	}
	// 对发送数据序列化，结果是[]byte
	data, _ := proto.Marshal(&msg)
	// 将序列化后的数据（[]byte类型）写入网络通道
	Write(conn, data)
	return "", "", nil
}

func SendResponse(conn net.Conn, ret_state ResponseState, ret_data string) (string, string, error) {
	msg := Response{
		From:      &ConnAddress{},
		To:        &ConnAddress{},
		StateCode: ret_state,
		Msg:       ret_data,
	}
	data, _ := proto.Marshal(&msg)
	Write(conn, data)
	return "", "", nil
}

func RecvRequest(conn net.Conn) *Request {
	request := &Request{}
	data := Read(conn)
	// 反序列化
	proto.Unmarshal(data, request)
	return request
}

func RecvResponse(conn net.Conn) *Response {
	response := &Response{}
	data := Read(conn)
	// 反序列化
	proto.Unmarshal(data, response)
	return response
}

func Write(conn net.Conn, data []byte) {
	chunk_index := 0
	chunk_size := 4095
	buf := make([]byte, 4096)
	// The make built-in function allocates and initializes an object of type slice, map, or chan (only). Like new, the first argument is a type, not a value. Unlike new, make's return type is the same as the type of its argument, not a pointer to it.
	for {
		// 检查data剩余多少字节够不够一个chunk（一个chunk=4095字节，首字节是1或者0）
		if len(data) - chunk_index * chunk_size <= chunk_size {
			copy(buf[1:], data[chunk_index*chunk_size:])
			copy(buf[:1], []byte("0")) // 首位=标志0表示没有下一端4K数据了
			conn.Write(buf)
			break
		}
		copy(buf[1:], data[chunk_index*chunk_size:(chunk_index+1)*chunk_size])
		copy(buf[:1], []byte("1"))	// 剩余字节还够一个chunk，首字节=1
		conn.Write(buf)
		chunk_index++
	}
}

func Read(conn net.Conn) []byte {
	reader := bufio.NewReader(conn)
	full_data := ""
	var buf [4096]byte	// 每次读取一个4kib，按照首字节是否是1确定是否退出
	for {
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Printf("read from conn failed, err:%v\n", err)
			break
		}
		if buf[0] == '0' {
			recv := string(buf[1:n])
			full_data += recv
			break
		} else {
			recv := string(buf[1:n])
			full_data += recv
		}
	}
	return []byte(full_data)
}
