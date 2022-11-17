package cmd

import (
	"HCPlatform/code/network"
	"github.com/spf13/cobra"
)

var (
	asDeviceClient  bool
	asControlServer bool
	cfgPath         string
	listAllDevices  bool
	user            string

	rootCmd = &cobra.Command{
		Use:   "hcp",
		Short: "HCP was developed by AIoT of Tsinghua University",
		Long:  "HCP is a heterogeneous computing platform, which was developed by AIoT of Tsinghua University's AIR Institute",
		Run: func(cmd *cobra.Command, args []string) {
			if asDeviceClient {
				network.NewConn("127.0.0.1", 9521)
			}
			if asControlServer {
				network.RunService("127.0.0.1", 9521)
			}

		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(connectCmd)
	rootCmd.PersistentFlags().BoolVarP(&listAllDevices, "list", "l", false, "list all devices")
	rootCmd.PersistentFlags().BoolVarP(&asDeviceClient, "cli", "c", false, "run as device client")
	rootCmd.PersistentFlags().BoolVarP(&asControlServer, "server", "s", false, "run as server")

}
