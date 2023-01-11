package pkg

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

type Cfg struct {
	ServerList  []ServerCfg `yaml:"ServersList"`
	DevicesList []DeviceCfg `yaml:"DevicesList"`
}

type ServerCfg struct {
	ServerName   string `yaml:"serverName"`
	Ip           string `yaml:"ip"`
	RegisterPort int    `yaml:"registerPort"`
	TerminalPort int    `yaml:"terminalPort"`
	ProfilePort  int    `yaml:"profilePort"`
}

type DeviceCfg struct {
	Ip           string `yaml:"ip"`
	TerminalPort int    `yaml:"terminalPort"`
	DeviceName   string `yaml:"deviceName"`
	DeviceId     uint64 `yaml:"deviceId"`
	Level        uint64 `yaml:"level"`
	Type         uint64 `yaml:"type"`
}

//	func (cfg *Cfg) GetServerCfg() []ServerCfg {
//		return cfg.serverList
//	}
//
//	func (cfg *Cfg) GetDeviceCfg() []DeviceCfg {
//		return cfg.devicesList
//	}
func (cfg *ServerCfg) GetNetAddress() string {
	return cfg.Ip
}

func (cfg *ServerCfg) GetRegisterPort() int {
	return cfg.RegisterPort
}
func (cfg *ServerCfg) GetTerminalPort() int {
	return cfg.TerminalPort
}
func (cfg *ServerCfg) GetProfilePort() int {
	return cfg.ProfilePort
}

func (cfg *DeviceCfg) GetNetAddress() string {
	return cfg.Ip
}

func (cfg *DeviceCfg) GetDeviceName() string {
	return cfg.DeviceName
}

func (cfg *DeviceCfg) GetDeviceID() uint64 {
	return cfg.DeviceId
}

func (cfg *DeviceCfg) GetDeviceLevel() uint64 {
	return cfg.Level
}

func (cfg *DeviceCfg) GetDeviceType() uint64 {
	return cfg.Type
}

func GetConfig(path string) ([]ServerCfg, []DeviceCfg) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Fatal happend when reading cfg file")
		return []ServerCfg{}, []DeviceCfg{}
	}
	cfg := Cfg{}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal(err)
		return []ServerCfg{}, []DeviceCfg{}
	}
	return cfg.ServerList, cfg.DevicesList
}
