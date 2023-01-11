package cmd

import (
	"HCPlatform/code/protos/register"
	"HCPlatform/code/protos/term"
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	listAllDevices    bool
	connectDeviceName string
	// TODO 暂时禁用
	interactive bool
	connectCmd  = &cobra.Command{
		Use:   "connect",
		Short: "Connect a device",
		Long:  "Conect a device which you have selected",
		Run: func(cmd *cobra.Command, args []string) {
			//TODO: 根据设备查询连接方式
			//
			res, _ := register.FastFreeDevices("127.0.0.1", 9520, []string{"ROG-LAPTOP"})
			res, _ = register.FastListAllDevices("127.0.0.1", 9520)
			devices, _ := register.FastAllocDevices("127.0.0.1", 9520, []string{"ROG-LAPTOP"})
			//log.Info(res)
			////TODO 根据设备号拿到目标Device的IP
			dstIP, dstPort := devices[0].GetDeviceAddress(), devices[0].GetTerminalPort()
			// We will access to device's command line,only for level1 device.You can launch ADBShell by your self
			if interactive {
				var inputContent string = ""
				shellId, _ := term.FastNewTerminal(dstIP, int(dstPort), 0, 0)
				for inputContent != "!exit" {
					//_, err := fmt.Scanln(&inputContent)
					reader := bufio.NewReader(os.Stdin)          // 标准输入输出
					inputContent, err := reader.ReadString('\n') // 回车结束
					if err != nil {
						log.Error(err)
						return
					}
					//shell, _ := pkg.NewPowerShell()
					//
					//sout, serr, _ := shell.Execute(inputContent)
					//fmt.Println(fmt.Sprintf("%s", sout))
					res, _ = term.FastExecCommand(dstIP, int(dstPort), 0, shellId, inputContent)
					fmt.Println(res)
					//log.Info(res)
				}
				shellId, _ = term.FastCloseTerminal(dstIP, int(dstPort), 0, shellId)
			}
			res, _ = register.FastFreeDevices("127.0.0.1", 9520, []string{"ROG-LAPTOP"})
			//log.Info(res)
		},
	}
)

func init() {
	connectCmd.PersistentFlags().StringVar(&connectDeviceName, "deviceName", "", "you have selected devicename")
	connectCmd.PersistentFlags().BoolVar(&listAllDevices, "listDevices", false, "list all devices registered to server")
	connectCmd.PersistentFlags().BoolVar(&interactive, "interactive", false, "Interactive access to the command line")
}

func preNUm(data byte) int {
	var mask byte = 0x80
	var num int = 0
	//8bit中首个0bit前有多少个1bits
	for i := 0; i < 8; i++ {
		if (data & mask) == mask {
			num++
			mask = mask >> 1
		} else {
			break
		}
	}
	return num
}
func isUtf8(data []byte) bool {
	i := 0
	for i < len(data) {
		if (data[i] & 0x80) == 0x00 {
			// 0XXX_XXXX
			i++
			continue
		} else if num := preNUm(data[i]); num > 2 {
			// 110X_XXXX 10XX_XXXX
			// 1110_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_0XXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_10XX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_110X 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// preNUm() 返回首个字节的8个bits中首个0bit前面1bit的个数，该数量也是该字符所使用的字节数
			i++
			for j := 0; j < num-1; j++ {
				//判断后面的 num - 1 个字节是不是都是10开头
				if (data[i] & 0xc0) != 0x80 {
					return false
				}
				i++
			}
		} else {
			//其他情况说明不是utf-8
			return false
		}
	}
	return true
}
