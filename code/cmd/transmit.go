package cmd

import "github.com/spf13/cobra"

var (
	transmitDeviceName     string
	transmitRemoteFilePath string
	transmitLocalFilePath  string

	uploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "upload your local data",
		Long:  "upload your local data",
		Run: func(cmd *cobra.Command, args []string) {
			//调用ssh去传输
		},
	}

	downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "download you remote data",
		Long:  "download you remote data",
		Run: func(cmd *cobra.Command, args []string) {
			//调用ssh去传输
		},
	}
)

func init() {
	uploadCmd.PersistentFlags().StringVar(&transmitDeviceName, "device", "", "")
	uploadCmd.PersistentFlags().StringVar(&transmitLocalFilePath, "local", "", "")
	uploadCmd.PersistentFlags().StringVar(&transmitRemoteFilePath, "remote", "", "")
	uploadCmd.MarkPersistentFlagRequired("device")
	uploadCmd.MarkPersistentFlagRequired("local")
	uploadCmd.MarkPersistentFlagRequired("remote")
	downloadCmd.PersistentFlags().StringVar(&transmitDeviceName, "device", "", "")
	downloadCmd.PersistentFlags().StringVar(&transmitLocalFilePath, "local", "", "")
	downloadCmd.PersistentFlags().StringVar(&transmitRemoteFilePath, "remote", "", "")
	downloadCmd.MarkPersistentFlagRequired("device")
	downloadCmd.MarkPersistentFlagRequired("local")
	downloadCmd.MarkPersistentFlagRequired("remote")
}