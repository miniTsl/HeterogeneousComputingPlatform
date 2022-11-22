package cmd // 指定模块名

import (
	"HCPlatform/code/network"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// 定义全局变量
var (
	asDeviceClient  bool
	asControlServer bool
	cfgPath         string
	listAllDevices  bool
	user            string
	// rootCmd是struct类型
	rootCmd = &cobra.Command{
		Use:   "hcp",
		Short: "HCP was developed by AIoT of Tsinghua University",
		Long:  "HCP is a heterogeneous computing platform, which was developed by AIoT of Tsinghua University's AIR Institute",
		// field Run func(cmd *cobra.Command, args []string)
		Run: func(cmd *cobra.Command, args []string) {
			// 根据参数确定做为client机还是server机，client负责发送、接受请求，server负责建立shell、处理请求并返回结果
			if asDeviceClient {
				log.Info("Starting ...,As Device Client")
				// Info logs a message at level Info on the standard logger.
				network.NewConn("183.172.197.36", 9521) // 数据发往地址
			}
			if asControlServer {
				log.Info("Starting ...,As Control Server")
				network.RunService("183.172.197.36", 9521) // 监听来自地址的数据
			}

		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(connectCmd) // AddCommand adds one or more commands to this parent command.
	// BoolVarP()支持-来简化输入
	rootCmd.PersistentFlags().BoolVarP(&listAllDevices, "list", "l", false, "list all devices")
	rootCmd.PersistentFlags().BoolVarP(&asDeviceClient, "cli", "c", false, "run as device client")
	rootCmd.PersistentFlags().BoolVarP(&asControlServer, "server", "s", false, "run as server")

}
