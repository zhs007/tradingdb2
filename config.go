package tradingdb2

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// TrNode2Cfg -
type TrNode2Cfg struct {
	Host  string `yaml:"host"`
	Token string `yaml:"token"`
}

// Config -
type Config struct {
	BatchCandleNums int          `yaml:"batchcandlenums"`
	DBPath          string       `yaml:"dbpath"`
	DBEngine        string       `yaml:"dbengine"`
	BindAddr        string       `yaml:"bindaddr"`
	LogLevel        string       `yaml:"loglevel"`
	LogPath         string       `yaml:"logpath"`
	DB2Markets      []string     `yaml:"db2markets"`
	Tokens          []string     `yaml:"tokens"`
	DataPath        string       `yaml:"datapath"`
	DataURL         string       `yaml:"dataurl"`
	Nodes           []TrNode2Cfg `yaml:"trnode2"`
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

	if cfg.LogPath == "" {
		cfg.LogPath = "./logs"
	}

	if cfg.BatchCandleNums <= 0 {
		cfg.BatchCandleNums = BatchCandleNums
	}

	if cfg.DataPath == "" {
		cfg.DataPath = "./output"
	}

	if cfg.DataURL == "" {
		cfg.DataURL = "http://127.0.0.1/"
	}

	return cfg, nil
}
