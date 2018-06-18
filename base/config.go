package ba

import (
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	OuterIf     string   `yaml:"outerIf"`
	BridgeName  string   `yaml:"brName"`
	LanNames    []string `yaml:"lan"`
	BridgeAddr  string   `yaml:"brAddr"`
	DnsmasqArgs []string `yaml:"dnsmasqArgs"`
	RPCHost     string   `yaml:"rpcHost"`
	RPCPort     int      `yaml:"rpcPort"`
}

type ConfigReaderFn func() ([]byte, error)

func ReadConfig(readConfig ConfigReaderFn) (*Config, error) {
	content, err := readConfig()
	if err != nil {
		return nil, err
	}
	cfg := Config{}
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.RPCHost == "" {
		cfg.RPCHost = "0.0.0.0"
	}
	if cfg.RPCPort == 0 {
		cfg.RPCPort = 2018
	}
	log.Printf("config readed: %+v", cfg)
	return &cfg, nil
}
