package cmd

import (
	"HCPlatform/code/internal"
	"HCPlatform/code/pkg"
	"HCPlatform/code/protos/register"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	asDeviceClient     bool
	asControlServer    bool
	asUserClient       bool
	cfgPath            string
	needListAllDevices bool
	user               string

	rootCmd = &cobra.Command{
		Use:   "hcp",
		Short: "HCP was developed by AIoT of Tsinghua University",
		Long:  "HCP is a heterogeneous computing platform, which was developed by AIoT of Tsinghua University's AIR Institute",
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
			serverIP, serverPort := "127.0.0.1", 9520
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
				l := pkg.NewNetListener(serverIP, serverPort)
				s := grpc.NewServer()
				rs := protos.RegisterService{}
				protos.RegisterReisgterServer(s, &rs)
				//reflection.Register(s)
				err := s.Serve(l)
				if err != nil {
					log.Error(err.Error())
					return
				}
			}
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(connectCmd)
	rootCmd.PersistentFlags().BoolVarP(&needListAllDevices, "list", "l", false, "list all devices")
	rootCmd.PersistentFlags().BoolVarP(&asDeviceClient, "cli", "c", false, "run as device client")
	rootCmd.PersistentFlags().BoolVarP(&asUserClient, "user", "u", false, "run as user client")
	rootCmd.PersistentFlags().BoolVarP(&asControlServer, "server", "s", false, "run as server")
	rootCmd.PersistentFlags().StringVar(&cfgPath, "cfg", "", "devcie configuration")
}
