package tradingdb2

import (
	"testing"

	"context"

	"github.com/stretchr/testify/assert"
	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
)

func Test_DB(t *testing.T) {
	db, err := NewDB("./unittestdata", "", "leveldb")
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
			&tradingpb.Candle{
				Open: 100,
			},
		},
	})
	assert.NoError(t, err)

	candles, err := db.GetCandles(context.Background(), "", "", []string{""}, 0, 0)
	assert.Error(t, err)

	candles, err = db.GetCandles(context.Background(), "bitmex", "", []string{""}, 0, 0)
	assert.Error(t, err)

	candles, err = db.GetCandles(context.Background(), "bitmex", "BTX", []string{""}, 0, 0)
	assert.NoError(t, err)
	assert.NotNil(t, candles)

	candles, err = db.GetCandles(context.Background(), "bitmex", "BTX", []string{"20200101"}, 0, 0)
	assert.NoError(t, err)

	assert.Equal(t, candles.Market, "bitmex")
	assert.Equal(t, candles.Symbol, "BTX")
	// assert.Equal(t, candles.Tag, "20200101")
	assert.Equal(t, len(candles.Candles), 1)
	assert.Equal(t, candles.Candles[0].Open, int64(100))

	candles, err = db.GetCandles(context.Background(), "bitmex", "BTX", []string{"20200102"}, 0, 0)
	assert.NoError(t, err)
	assert.NotNil(t, candles)

	t.Logf("Test_DB OK")
}
