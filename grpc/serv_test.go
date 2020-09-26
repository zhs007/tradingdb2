package tradingdb2grpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	tradingdb2 "github.com/zhs007/tradingdb2"
	tradingdb2pb "github.com/zhs007/tradingdb2/tradingdb2pb"
)

func Test_Serv(t *testing.T) {
	cfg, err := tradingdb2.LoadConfig("../unittestdata/serv.config.yaml")
	assert.NoError(t, err)

	serv, err := NewServ(cfg)
	assert.NoError(t, err)

	go func() {
		serv.Start(context.Background())
	}()

	time.Sleep(time.Second * 3)

	client0, err := NewClient("127.0.0.1:5002", "123")
	assert.NoError(t, err)

	candles := &tradingdb2pb.Candles{
		Market: "bitmex",
		Symbol: "BTX",
		Tag:    "20200101",
	}

	for i := 0; i < 100; i++ {
		candles.Candles = append(candles.Candles, &tradingdb2pb.Candle{
			Ts:   int64(i),
			Open: 100 + int64(i),
		})
	}

	ret, err := client0.UpdCandles(context.Background(), candles, 30, nil)
	assert.Error(t, err)
	assert.Nil(t, ret)

	replygetcandles, err := client0.GetCandles(context.Background(), "bitmex", "BTX", "20200101", nil)
	assert.Error(t, err)
	assert.Nil(t, replygetcandles)

	client1, err := NewClient("127.0.0.1:5002", "wzDkh9h2fhfUVuS9jZ8uVbhV3vC5AWX3")
	assert.NoError(t, err)

	ret, err = client1.UpdCandles(context.Background(), candles, 30, nil)
	assert.NoError(t, err)
	assert.Equal(t, ret.LengthOK, int32(100))

	dbcandles, err := serv.DB.GetCandles(context.Background(), "bitmex", "BTX", "20200101")
	assert.NoError(t, err)
	assert.Equal(t, dbcandles.Market, "bitmex")
	assert.Equal(t, dbcandles.Symbol, "BTX")
	assert.Equal(t, dbcandles.Tag, "20200101")
	assert.Equal(t, len(dbcandles.Candles), 100)
	for i := 0; i < 100; i++ {
		assert.Equal(t, dbcandles.Candles[i].Ts, int64(i))
		assert.Equal(t, dbcandles.Candles[i].Open, int64(100+i))
	}

	replygetcandles, err = client1.GetCandles(context.Background(), "bitmex", "BTX", "20200101", nil)
	assert.NoError(t, err)
	assert.NotNil(t, replygetcandles)
	assert.Equal(t, replygetcandles.Market, "bitmex")
	assert.Equal(t, replygetcandles.Symbol, "BTX")
	assert.Equal(t, replygetcandles.Tag, "20200101")
	assert.Equal(t, len(replygetcandles.Candles), 100)

	for i := 0; i < 100; i++ {
		assert.Equal(t, replygetcandles.Candles[i].Ts, int64(i))
		assert.Equal(t, replygetcandles.Candles[i].Open, int64(100+i))
	}

	// 100 < 200
	candles1 := &tradingdb2pb.Candles{
		Market: "bitmex",
		Symbol: "BTX",
		Tag:    "20200102",
	}

	for i := 0; i < 100; i++ {
		candles1.Candles = append(candles1.Candles, &tradingdb2pb.Candle{
			Ts:   int64(i),
			Open: 100 + int64(i),
		})
	}

	ret, err = client1.UpdCandles(context.Background(), candles1, 200, nil)
	assert.NoError(t, err)
	assert.Equal(t, ret.LengthOK, int32(100))

	replygetcandles, err = client1.GetCandles(context.Background(), "bitmex", "BTX", "20200102", nil)
	assert.NoError(t, err)
	assert.NotNil(t, replygetcandles)
	assert.Equal(t, replygetcandles.Market, "bitmex")
	assert.Equal(t, replygetcandles.Symbol, "BTX")
	assert.Equal(t, replygetcandles.Tag, "20200102")
	assert.Equal(t, len(replygetcandles.Candles), 100)

	for i := 0; i < 100; i++ {
		assert.Equal(t, replygetcandles.Candles[i].Ts, int64(i))
		assert.Equal(t, replygetcandles.Candles[i].Open, int64(100+i))
	}

	// 0 < 200
	candles2 := &tradingdb2pb.Candles{
		Market: "bitmex",
		Symbol: "BTX",
		Tag:    "20200103",
	}

	ret, err = client1.UpdCandles(context.Background(), candles2, 200, nil)
	assert.NoError(t, err)
	assert.Equal(t, ret.LengthOK, int32(0))

	replygetcandles, err = client1.GetCandles(context.Background(), "bitmex", "BTX", "20200103", nil)
	assert.NoError(t, err)
	assert.NotNil(t, replygetcandles)
	assert.Equal(t, replygetcandles.Market, "bitmex")
	assert.Equal(t, replygetcandles.Symbol, "BTX")
	assert.Equal(t, replygetcandles.Tag, "20200103")
	assert.Equal(t, len(replygetcandles.Candles), 0)

	serv.Stop()

	t.Logf("Test_Serv OK")
}
