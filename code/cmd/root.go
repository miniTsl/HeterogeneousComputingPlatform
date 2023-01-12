package cmd

import (
	"HCPlatform/code/internal"
	"HCPlatform/code/pkg"
	"HCPlatform/code/protos/exec"
	"HCPlatform/code/protos/register"
	"HCPlatform/code/protos/term"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	currentDevices     []*register.DeviceMessage
	asDeviceClient     bool
	asControlServer    bool
	asUserClient       bool
	cfgPath            string
	needListAllDevices bool

	rootCmd = &cobra.Command{
		Use:   "hcp",
		Short: "HCP was developed by AIoT of AIR",
		Long:  "HCP is a heterogeneous computing platform, which was developed by AIoT of Tsinghua University's AIR Institute",
		Run: func(cmd *cobra.Command, args []string) {
			serverList, deviceList := pkg.GetConfig(cfgPath)
			serverCfg := serverList[0]
			serverIP := serverCfg.GetNetAddress()
			registerPort, terminalPort, profilePort := serverCfg.GetRegisterPort(), serverCfg.GetTerminalPort(), serverCfg.GetProfilePort()
			if asControlServer {
				go launchRegisterService(serverIP, registerPort)
				go launchProfileService(serverIP, profilePort)
				//launchTerminalService(serverIP, terminalPort)
				for {

				}
			} else if asDeviceClient {
				deviceMsg := make([]*register.DeviceMessage, len(deviceList))
				for i, deviceCfg := range deviceList {
					deviceMsg[i] = &register.DeviceMessage{
						DeviceName:    deviceCfg.GetDeviceName(),
						DeviceAddress: deviceCfg.GetNetAddress(),
						XLevel:        register.DeviceMessage_Level(deviceCfg.GetDeviceLevel()),
						XType:         register.DeviceMessage_Type(deviceCfg.GetDeviceType()),
						TerminalPort:  int32(deviceCfg.TerminalPort),
					}
				}
				currentDevices = deviceMsg
				register.FastRegisterDevices(serverIP, registerPort, deviceMsg)
				for _, deviceCfg := range deviceList {
					launchTerminalService(deviceCfg.GetNetAddress(), deviceCfg.TerminalPort)
				}
			}
			fmt.Println(registerPort, terminalPort, profilePort)
		},
	}
)

func launchRegisterService(ip string, port int) {
	l := internal.NewNetListener(ip, port)
	s := grpc.NewServer()
	rs := register.RegisterService{}
	register.RegisterReisgterServer(s, &rs)
	err := s.Serve(l)
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Info(fmt.Sprintf("launchRegisterService at %s:%d", ip, port))
}

func launchProfileService(ip string, port int) {
	l := internal.NewNetListener(ip, port)
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
	log.Info(fmt.Sprintf("launchProfileService at %s:%d", ip, port))
}

func launchTerminalService(ip string, port int) {
	l := internal.NewNetListener(ip, port)
	// 此处设置最佳发送文件大小512M
	var options = []grpc.ServerOption{
		grpc.MaxRecvMsgSize(1024 * 1024 * 512),
		grpc.MaxSendMsgSize(1024 * 1024 * 512),
	}
	s := grpc.NewServer(options...)
	rs := term.TermnialService{}
	term.RegisterTerminalServer(s, &rs)
	err := s.Serve(l)
	if err != nil {
		log.Error(err.Error())
		return
	}
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(profileCmd)
	rootCmd.PersistentFlags().BoolVarP(&needListAllDevices, "list", "l", false, "list all devices")
	rootCmd.PersistentFlags().BoolVarP(&asDeviceClient, "cli", "c", false, "run as device client")
	rootCmd.PersistentFlags().BoolVarP(&asControlServer, "server", "s", false, "run as server")
	rootCmd.PersistentFlags().StringVar(&cfgPath, "cfg", "", "devcie configuration")
}
