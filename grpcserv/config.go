package tradingdb2serv

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config - config
type Config struct {
	BindAddr string `yaml:"bindaddr"`
}

// LoadConfig - load configuration
func LoadConfig(fn string) (*Config, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
