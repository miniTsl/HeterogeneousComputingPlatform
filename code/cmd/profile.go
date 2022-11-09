package cmd

import (
	"github.com/spf13/cobra"
)

var (
	modelPath  string
	framework  string
	deviceName string

	profileCmd = &cobra.Command{
		Use:   "profile",
		Short: "profile model",
		Long:  "profile model,support tflite,paddlelite,onnxruntime",
		Run: func(cmd *cobra.Command, args []string) {
			//shell, err := pkg.NewPowerShell(1)
			//if err == nil {
			//	out, err := shell.Exec("pwd")
			//	if err == nil {
			//		fmt.Println(out)
			//	}
			//	out, err = shell.Exec("ls")
			//	if err == nil {
			//		fmt.Println(out)
			//	}
			//	shell.Exit()
			//}

		},
	}
)

func init() {
	profileCmd.PersistentFlags().StringVar(&modelPath, "modelPath", "", "")
	profileCmd.MarkPersistentFlagRequired("modelPath")

	profileCmd.PersistentFlags().StringVar(&framework, "framework", "", "")
	profileCmd.MarkPersistentFlagRequired("framework")

	profileCmd.PersistentFlags().StringVar(&deviceName, "deviceName", "", "")
	profileCmd.MarkPersistentFlagRequired("deviceName")

}
