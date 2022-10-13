package cmd

import(
	"fmt"
	"github.com/spf13/cobra"
    "github.com/spf13/viper"
)


var (
	rootCmd = &cobra.Command{
		Use: "main",
		Short: "HCPlatrom devloped by AIoT",
		Long: "xxxx",
		Run : func(cmd *cobra.Command, args [] string){
			fmt.Println(time.Now())
		}
	}
	
)


func execute() error{
	return rootCmd.Execute();
}


func init(){
	cobra.OnInitialize(initConfig)

    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
    rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
    rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
    rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
    viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
    viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
    viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
    viper.SetDefault("license", "apache")
}