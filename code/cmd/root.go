package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

var (
	asDeviceClient  bool
	asControlServer bool
	listAllDevices  bool
	user            string

	rootCmd = &cobra.Command{
		Use:   "hcp",
		Short: "HCP was developed by AIoT of Tsinghua University",
		Long:  "HCP is a heterogeneous computing platform, which was developed by AIoT of Tsinghua University's AIR Institute",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(time.Now())
			fmt.Println(user)
			if listAllDevices {
				fmt.Println("We Need List all online devices")
			} else {
			}

		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {

	rootCmd.PersistentFlags().BoolVarP(&listAllDevices, "list", "l", false, "list all devices")
	rootCmd.AddCommand(connectCmd)
}
