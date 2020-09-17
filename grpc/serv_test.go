package tradingdb2grpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	tradingdb2 "github.com/zhs007/tradingdb2"
	tradingdb2pb "github.com/zhs007/tradingdb2/pb"
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

	ret, err := client0.UpdCandles(context.Background(), &tradingdb2pb.Candles{
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
	assert.Equal(t, ret.LengthOK, int32(0))
	assert.Equal(t, ret.Error, "invalid token")

	client1, err := NewClient("127.0.0.1:5002", "wzDkh9h2fhfUVuS9jZ8uVbhV3vC5AWX3")
	assert.NoError(t, err)

	ret, err = client1.UpdCandles(context.Background(), &tradingdb2pb.Candles{
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
	assert.Equal(t, ret.LengthOK, int32(1))
	assert.Equal(t, ret.Error, "")

	candles, err := serv.db.GetCandles(context.Background(), "bitmex", "BTX", "20200101")
	assert.NoError(t, err)
	assert.Equal(t, candles.Market, "bitmex")
	assert.Equal(t, candles.Symbol, "BTX")
	assert.Equal(t, candles.Tag, "20200101")
	assert.Equal(t, len(candles.Candles), 1)
	assert.Equal(t, candles.Candles[0].Open, int64(100))

	serv.Stop()

	t.Logf("Test_Serv OK")
}
