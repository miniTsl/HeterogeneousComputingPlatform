package main

import (
	"HCPlatform/code/cmd"
	"os"
	log "github.com/sirupsen/logrus"
)

func main() {
	// cmd是一个自定义的模块
	// := 是声明并赋值，并且系统自动推断类型，不需要var关键字
	// 普通的=赋值运算符需要先用var声明变量名
	err := cmd.Execute()
	// nil is a predeclared identifier representing the zero value for a pointer, channel, func, interface, map, or slice type. cmd.Execute()返回的是一个interface。如果有错则直接退出
	// error is the conventional interface for representing an error condition
	if err != nil {
		return
	}

}

// 初始化log的格式
func init() {
	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
	// Stdin, Stdout, and Stderr are open Files pointing to the standard input, standard output, and standard error file descriptors. Note that the Go runtime writes to standard error for panics and crashes; closing Stderr may cause those messages to go elsewhere, perhaps to a file opened later.
	log.SetLevel(log.InfoLevel) // info级别之上的log信息会输出到stdout中
	// log信息的format
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	// 打印详细信息，返回calling method的文件位置和行数
	log.SetReportCaller(true)
}
