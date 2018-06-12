package ba

import (
	"gopkg.in/yaml.v2"
)

type Config struct {
	OuterIf     string   `yaml:"outerIf"`
	BridgeName  string   `yaml:"brName"`
	LanNames    []string `yaml:"lan"`
	BridgeAddr  string   `yaml:"brAddr"`
	DnsmasqArgs []string `yaml:"dnsmasqArgs"`
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
	return &cfg, nil
}
