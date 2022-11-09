package cmd

import (
	"github.com/spf13/cobra"
)

var (
	connectDeviceName string
	connectCmd        = &cobra.Command{
		Use:   "connect",
		Short: "Connect a device",
		Long:  "Conect a device which you have selected",
		Run: func(cmd *cobra.Command, args []string) {
			//TODO: 根据设备查询连接方式
			//network.LoginByPassword("192.168.13.189", 22, "yang", "274085")
		},
	}
)

func init() {
	connectCmd.PersistentFlags().StringVar(&connectDeviceName, "device", "", "you have selected device name")
	connectCmd.MarkPersistentFlagRequired("device")

}
