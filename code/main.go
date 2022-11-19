package main

import (
	"HCPlatform/code/cmd"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {

	err := cmd.Execute()
	if err != nil {
		return
	}

}

// 初始化log格式
func init() {
	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	// 打印详细信息，返回方法文件位置和行数
	log.SetReportCaller(true)
}
