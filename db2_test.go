package tradingdb2

import (
	"testing"

	"context"

	"github.com/stretchr/testify/assert"
	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
)

func Test_DB2(t *testing.T) {
	cfg := &Config{
		DBPath:   "./unittestdata",
		DBEngine: "leveldb",
		DB2Markets: []string{
			"bitmex",
		},
	}
	db, err := NewDB2(cfg)
	assert.NoError(t, err)

	err = db.UpdCandles(context.Background(), &tradingpb.Candles{})
	assert.Error(t, err)

	err = db.UpdCandles(context.Background(), &tradingpb.Candles{Market: "bitmex"})
	assert.Error(t, err)

	err = db.UpdCandles(context.Background(), &tradingpb.Candles{Market: "bitmex", Symbol: "BTX"})
	assert.Error(t, err)

	err = db.UpdCandles(context.Background(), &tradingpb.Candles{
		Market: "bitmex",
		Symbol: "BTX",
		Tag:    "20200101",
		Candles: []*tradingpb.Candle{
			{
				Ts:   1,
				Open: 100,
			},
		},
	})
	assert.NoError(t, err)

	candles, err := db.GetCandles(context.Background(), "", "", 0, 0, 0)
	assert.Error(t, err)

	candles, err = db.GetCandles(context.Background(), "bitmex", "", 0, 0, 0)
	assert.Error(t, err)

	candles, err = db.GetCandles(context.Background(), "bitmex", "BTX", 0, 0, 0)
	assert.NoError(t, err)
	assert.NotNil(t, candles)

	candles, err = db.GetCandles(context.Background(), "bitmex", "BTX", 0, 0, 0)
	assert.NoError(t, err)

	assert.Equal(t, candles.Market, "bitmex")
	assert.Equal(t, candles.Symbol, "BTX")
	// assert.Equal(t, candles.Tag, "20200101")
	assert.Equal(t, len(candles.Candles), 1)
	assert.Equal(t, candles.Candles[0].Open, int64(100))

	candles, err = db.GetCandles(context.Background(), "bitmex", "BTX", 0, 0, 0)
	assert.NoError(t, err)
	assert.NotNil(t, candles)

	lst := []*tradingpb.Candle{}

	for i := 0; i < 100; i++ {
		cv := &tradingpb.Candle{
			Ts:   int64(100 + i),
			Open: 100,
		}

		lst = append(lst, cv)
	}

	err = db.UpdCandles(context.Background(), &tradingpb.Candles{
		Market:  "bitmex",
		Symbol:  "BTX",
		Tag:     "20200101",
		Candles: lst,
	})
	assert.NoError(t, err)

	candles, err = db.GetCandles(context.Background(), "bitmex", "BTX", 150, 0, 10)
	assert.NoError(t, err)

	assert.Equal(t, candles.Market, "bitmex")
	assert.Equal(t, candles.Symbol, "BTX")
	// assert.Equal(t, candles.Tag, "20200101")
	assert.Equal(t, len(candles.Candles), 60)
	assert.Equal(t, candles.Candles[0].Open, int64(100))

	t.Logf("Test_DB OK")
}
