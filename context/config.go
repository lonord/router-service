package context

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

func ReadConfig(readConfig ConfigReaderFn, c *MainContext) error {
	content, err := readConfig()
	if err != nil {
		return err
	}
	cfg := Config{}
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return err
	}
	c.Cfg = cfg
	return nil
}
