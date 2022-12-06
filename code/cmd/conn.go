package cmd

import (
	"github.com/spf13/cobra"
)

const (
	LocalConnection  int = 0
	RemoteConnection int = 1
)

var (
	connectType       int
	connectDeviceName string

	connectCmd = &cobra.Command{
		Use:   "connect",
		Short: "Connect a device",
		Long:  "Conect a device which you have selected",
		Run: func(cmd *cobra.Command, args []string) {
			//TODO: 根据设备查询连接方式
			//pkg.LoginByPassword("192.168.13.189", 22, "yang", "274085")
		},
	}
)

func init() {
	connectCmd.PersistentFlags().StringVar(&connectDeviceName, "deviceName", "", "you have selected device name")
	connectCmd.MarkPersistentFlagRequired("device")
	connectCmd.PersistentFlags().IntVar(&connectType, "connectType", 0, "Connect remote deivce or local device")
	connectCmd.MarkPersistentFlagRequired("connectType")

}
