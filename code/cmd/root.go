package cmd // 指定模块名

import (
	"HCPlatform/code/internal"
	"HCPlatform/code/pkg"
	"HCPlatform/code/protos/exec"
	"HCPlatform/code/protos/register"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
	"os"
)

// 定义全局变量，尤其是rootCmd
var (
	asDeviceClient     bool
	asControlServer    bool
	asUserClient       bool
	cfgPath            string
	needListAllDevices bool
	user               string

	rootCmd = &cobra.Command{
		Use:   "hcp",
		Short: "HCP was developed by AIoT of AIR",
		Long:  "HCP is a heterogeneous computing platform, which was developed by AIoT of Tsinghua University's AIR Institute",
		// field Run func(cmd *cobra.Command, args []string)
		Run: func(cmd *cobra.Command, args []string) {
			data, err := os.ReadFile(cfgPath)
			if err != nil {
				log.Fatal("Fatal happend when reading cfg file")
				return
			}
			cfg := internal.Cfg{}

			err = yaml.Unmarshal(data, &cfg)
			if err != nil {
				log.Fatal(err)
				return
			}
			//deviceCfg := cfg.GetDeviceCfg()
			//serverCfg := cfg.GetServerCfg()
			serverIP, serverPort := "0.0.0.0", 9520
			if asDeviceClient {

				conn, err := grpc.Dial(fmt.Sprintf("%s:%d", serverIP, serverPort), grpc.WithInsecure())
				if err != nil {
					log.Fatal(err)
				}
				defer conn.Close()
				c := protos.NewReisgterClient(conn)
				var devices []*protos.DeviceMessage
				req := &protos.RegisterRequest{Devices: devices}
				res, err := c.ResgisterDevice(context.Background(), req)
				if err != nil {
					log.Fatal(err)
				}

				log.Info(res.Msg)
			}
			if asControlServer {
				//launchRegisterService(serverIP,serverPort)
				launchProfileService(serverIP, serverPort)
			}
		},
	}
)

func launchRegisterService(ip string, port int) {
	l := pkg.NewNetListener(ip, port)
	s := grpc.NewServer()
	rs := protos.RegisterService{}
	protos.RegisterReisgterServer(s, &rs)
	err := s.Serve(l)
	if err != nil {
		log.Error(err.Error())
		return
	}

}

func launchProfileService(ip string, port int) {
	l := pkg.NewNetListener(ip, port)
	// 此处设置最佳发送文件大小512M
	var options = []grpc.ServerOption{
		grpc.MaxRecvMsgSize(1024 * 1024 * 512),
		grpc.MaxSendMsgSize(1024 * 1024 * 512),
	}
	s := grpc.NewServer(options...)
	rs := exec.ProfileService{}
	exec.RegisterProfileServer(s, &rs)
	err := s.Serve(l)
	if err != nil {
		log.Error(err.Error())
		return
	}
}

func Execute() error {
	return rootCmd.Execute()	// cobra提供的方法
}

func init() {
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(profileCmd)
	rootCmd.PersistentFlags().BoolVarP(&needListAllDevices, "list", "l", false, "list all devices")
	rootCmd.PersistentFlags().BoolVarP(&asDeviceClient, "cli", "c", false, "run as device client")
	rootCmd.PersistentFlags().BoolVarP(&asUserClient, "user", "u", false, "run as user client")
	rootCmd.PersistentFlags().BoolVarP(&asControlServer, "server", "s", false, "run as server")
	rootCmd.PersistentFlags().StringVar(&cfgPath, "cfg", "", "devcie configuration")
}
