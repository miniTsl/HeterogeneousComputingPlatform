package network

import (
	"bufio"
	"fmt"
	"google.golang.org/protobuf/proto"
	"net"
)

func SendRequest(conn net.Conn, cmd string) (string, string, error) {
	msg := Request{
		From:  &ConnAddress{},
		To:    &ConnAddress{},
		XType: Request_command,
		Pyload: &Request_CommandRequest{
			CommandRequest: &CommandRequest{
				ShellId: 0,
				Command: cmd,
			},
		},
	}
	data, _ := proto.Marshal(&msg)
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
	proto.Unmarshal(data, request)
	return request
}

func RecvResponse(conn net.Conn) *Response {
	response := &Response{}
	data := Read(conn)
	proto.Unmarshal(data, response)
	return response
}

func Write(conn net.Conn, data []byte) {
	chunk_index := 0
	chunk_size := 4095
	buf := make([]byte, 4096)
	for {
		//检查data剩余多少字节没有传输
		if len(data)-chunk_index*chunk_size <= chunk_size {
			copy(buf[1:], data[chunk_index*chunk_size:])
			copy(buf[:1], []byte("0")) //标志0表示没有下一个4K数据了
			conn.Write(buf)
			break
		}
		copy(buf[1:], data[chunk_index*chunk_size:(chunk_index+1)*chunk_size])
		copy(buf[:1], []byte("1"))
		conn.Write(buf)
		chunk_index++
	}
}

func Read(conn net.Conn) []byte {
	reader := bufio.NewReader(conn)
	full_data := ""
	var buf [4096]byte
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
