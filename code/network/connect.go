package network

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
)

func LoginByPassword(host string, port int, username string, password string) {
	config := &ssh.ClientConfig{
		Timeout:         0,
		User:            username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal("创建ssh失败", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("是啊比", err)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     //打开回显
		ssh.TTY_OP_ISPEED: 14400, //输入速率
		ssh.TTY_OP_OSPEED: 14400, //输出速率
		ssh.VSTATUS:       1,
	}

	//使用VT100终端来实现tab键提示，上下键查看历史命令，clear键清屏等操作
	//VT100 start
	//windows下不支持VT100
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer terminal.Restore(fd, oldState)

	termWidth, termHeight, err := terminal.GetSize(fd)

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	//此时打开终端
	err = session.RequestPty("xterm", termHeight, termWidth, modes)
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = session.Shell()
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = session.Wait()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
