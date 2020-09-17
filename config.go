package tradingdb2

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config -
type Config struct {
	DBPath   string   `yaml:"dbpath"`
	DBEngine string   `yaml:"dbengine"`
	BindAddr string   `yaml:"bindaddr"`
	LogLevel string   `yaml:"loglevel"`
	Tokens   []string `yaml:"tokens"`
}

// LoadConfig - load config
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
