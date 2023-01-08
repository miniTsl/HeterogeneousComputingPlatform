package pkg

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

type Cfg struct {
	jumpServerCfg ServerCfg `yaml:"ServerCfg"`
	deviceCfg     DeviceCfg `yaml:"DeviceCfg"`
}

type ServerCfg struct {
	ip           string `yaml:"ip"`
	registerPort int    `yaml:"registerPort"`
	terminalPort int    `yaml:"terminalPort"`
	profilePort  int    `yaml:"profilePort"`
}

type DeviceCfg struct {
	ip           string `yaml:"ip"`
	terminalPort int    `yaml:"terminalPort"`
	deviceName   string `yaml:"deviceName"`
	deviceId     uint64 `yaml:"deviceId"`
	_level       uint64 `yaml:"level"`
	_type        uint64 `yaml:"type"`
}

func (cfg *Cfg) GetServerCfg() ServerCfg {
	return cfg.jumpServerCfg
}

func (cfg *Cfg) GetDeviceCfg() DeviceCfg {
	return cfg.deviceCfg
}

func (cfg *ServerCfg) GetNetAddress() string {
	return cfg.ip
}
func (cfg *ServerCfg) GetRegisterPort() int {
	return cfg.registerPort
}
func (cfg *ServerCfg) GetTerminalPort() int {
	return cfg.terminalPort
}
func (cfg *ServerCfg) GetProfilePort() int {
	return cfg.profilePort
}

func (cfg *DeviceCfg) GetNetAddress() string {
	return cfg.ip
}

func (cfg *DeviceCfg) GetDeviceName() string {
	return cfg.deviceName
}

func (cfg *DeviceCfg) GetDeviceID() uint64 {
	return cfg.deviceId
}

func (cfg *DeviceCfg) GetDeviceLevel() uint64 {
	return cfg._level
}

func (cfg *DeviceCfg) GetDeviceType() uint64 {
	return cfg._type
}

func GetConfig(path string) (ServerCfg, DeviceCfg) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Fatal happend when reading cfg file")
		return ServerCfg{}, DeviceCfg{}
	}
	cfg := make(map[string]interface{})
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal(err)
		return ServerCfg{}, DeviceCfg{}
	}
	//compellent type cast
	deviceCfg := (cfg["DeviceCfg"]).(DeviceCfg)
	serverCfg := (cfg["ServerCfg"]).(ServerCfg)
	return serverCfg, deviceCfg
}
