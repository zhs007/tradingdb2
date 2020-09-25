package tradingdb2

import (
	"testing"

	"context"

	"github.com/stretchr/testify/assert"
	tradingdb2pb "github.com/zhs007/tradingdb2/tradingdb2pb"
)

func Test_DB(t *testing.T) {
	db, err := NewDB("./unittestdata", "", "leveldb")
	assert.NoError(t, err)

	err = db.UpdCandles(context.Background(), &tradingdb2pb.Candles{})
	assert.Error(t, err)

	err = db.UpdCandles(context.Background(), &tradingdb2pb.Candles{Market: "bitmex"})
	assert.Error(t, err)

	err = db.UpdCandles(context.Background(), &tradingdb2pb.Candles{Market: "bitmex", Symbol: "BTX"})
	assert.Error(t, err)

	err = db.UpdCandles(context.Background(), &tradingdb2pb.Candles{
		Market: "bitmex",
		Symbol: "BTX",
		Tag:    "20200101",
		Candles: []*tradingdb2pb.Candle{
			&tradingdb2pb.Candle{
				Open: 100,
			},
		},
	})
	assert.NoError(t, err)

	candles, err := db.GetCandles(context.Background(), "", "", "")
	assert.Error(t, err)

	candles, err = db.GetCandles(context.Background(), "bitmex", "", "")
	assert.Error(t, err)

	candles, err = db.GetCandles(context.Background(), "bitmex", "BTX", "")
	assert.Error(t, err)

	candles, err = db.GetCandles(context.Background(), "bitmex", "BTX", "20200101")
	assert.NoError(t, err)

	assert.Equal(t, candles.Market, "bitmex")
	assert.Equal(t, candles.Symbol, "BTX")
	assert.Equal(t, candles.Tag, "20200101")
	assert.Equal(t, len(candles.Candles), 1)
	assert.Equal(t, candles.Candles[0].Open, int64(100))

	candles, err = db.GetCandles(context.Background(), "bitmex", "BTX", "20200102")
	assert.NoError(t, err)
	assert.Nil(t, candles)

	t.Logf("Test_DB OK")
}
