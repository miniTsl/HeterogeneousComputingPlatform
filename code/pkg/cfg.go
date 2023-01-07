package pkg

type Cfg struct {
	jumpServerCfg ServerCfg `yaml:"jumpServer"`
	deviceCfg     DeviceCfg `yaml:"device"`
}

type ServerCfg struct {
	ip   string `yaml:"ip"`
	port int    `yaml:"port"`
}

type DeviceCfg struct {
	ip         string `yaml:"ip"`
	port       int    `yaml:"port"`
	deviceName string `yaml:"deviceName"`
	deviceId   uint64 `yaml:"deviceId"`
	_level     uint64 `yaml:"level"`
	_type      uint64 `yaml:"type"`
}

func (cfg *Cfg) GetServerCfg() ServerCfg {
	return cfg.jumpServerCfg
}

func (cfg *Cfg) GetDeviceCfg() DeviceCfg {
	return cfg.deviceCfg
}

func (cfg *ServerCfg) GetNetAddress() (string, int) {
	return cfg.ip, cfg.port
}

func (cfg *DeviceCfg) GetNetAddress() (string, int) {
	return cfg.ip, cfg.port
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
