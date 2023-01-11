package cmd

import (
	"HCPlatform/code/protos/register"
	"HCPlatform/code/protos/term"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	listAllDevices    bool
	connectDeviceName string
	connectCmd        = &cobra.Command{
		Use:   "connect",
		Short: "Connect a device",
		Long:  "Conect a device which you have selected",
		Run: func(cmd *cobra.Command, args []string) {
			//TODO: 根据设备查询连接方式

			res, _ := register.FastFreeDevices("127.0.0.1", 9520, []string{"ROG-LAPTOP"})
			res, _ = register.FastListAllDevices("127.0.0.1", 9520)
			log.Info(res)
			devices, _ := register.FastAllocDevices("127.0.0.1", 9520, []string{"ROG-LAPTOP"})
			log.Info(res)
			//TODO 根据设备号拿到目标Device的IP
			dstIP, dstPort := devices[0].GetDeviceAddress(), devices[0].GetTerminalPort()
			log.Info(fmt.Sprintf("Get,%s,%d", dstIP, dstPort))
			shellId, _ := term.FastNewTerminal(dstIP, int(dstPort), 0, 0)
			if shellId != 0 {
				res, _ = term.FastExecCommand(dstIP, int(dstPort), 0, shellId, "ls")
				log.Info(res)
				command := "cd .."
				fmt.Println(isUtf8([]byte(command)))
				res, _ = term.FastExecCommand(dstIP, int(dstPort), 0, shellId, command)
				log.Info(res)
				command = "pwd"
				fmt.Println(isUtf8([]byte(command)))
				res, _ = term.FastExecCommand(dstIP, int(dstPort), 0, shellId, command)
				log.Info(res)
				shellId, _ = term.FastCloseTerminal(dstIP, int(dstPort), 0, shellId)
			}
			res, _ = register.FastFreeDevices("127.0.0.1", 9520, []string{"ROG-LAPTOP"})
			log.Info(res)

		},
	}
)

func init() {
	connectCmd.PersistentFlags().StringVar(&connectDeviceName, "deviceName", "", "you have selected devicename")
	connectCmd.PersistentFlags().BoolVar(&listAllDevices, "listDevices", false, "list all devices registered to server")

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
