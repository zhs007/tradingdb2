package tradingdb2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadConfig(t *testing.T) {

	cfg, err := LoadConfig("./cfg/config.yaml.default")
	assert.NoError(t, err)

	assert.Equal(t, cfg.BatchCandleNums, BatchCandleNums)
	assert.Equal(t, cfg.DBPath, "./data")
	assert.Equal(t, cfg.DBEngine, "leveldb")
	assert.Equal(t, cfg.BindAddr, "0.0.0.0:5002")
	assert.Equal(t, cfg.LogLevel, "debug")
	assert.Equal(t, len(cfg.DB2Markets), 2)
	assert.Equal(t, cfg.DB2Markets[0], "jrj")
	assert.Equal(t, cfg.DB2Markets[1], "bitmex")
	assert.Equal(t, len(cfg.Tokens), 1)
	assert.Equal(t, cfg.Tokens[0], "wzDkh9h2fhfUVuS9jZ8uVbhV3vC5AWX3")
	assert.Equal(t, len(cfg.Nodes), 1)
	assert.Equal(t, cfg.Nodes[0].Host, "127.0.0.1:123")
	assert.Equal(t, cfg.Nodes[0].Token, "123456")

	cfg, err = LoadConfig("./cfg/config.yaml.default0")
	assert.Error(t, err)

	cfg, err = LoadConfig("./unittestdata/err.yaml")
	assert.Error(t, err)

	t.Logf("Test_string2LogLevel OK")
}
