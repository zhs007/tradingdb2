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

	replygetcandles, err := client0.GetCandles(context.Background(), "bitmex", "BTX", []string{"20200101"}, 0, 0, nil)
	assert.Error(t, err)
	assert.Nil(t, replygetcandles)

	client1, err := NewClient("127.0.0.1:5002", "wzDkh9h2fhfUVuS9jZ8uVbhV3vC5AWX3")
	assert.NoError(t, err)

	ret, err = client1.UpdCandles(context.Background(), candles, 30, nil)
	assert.NoError(t, err)
	assert.Equal(t, ret.LengthOK, int32(100))

	dbcandles, err := serv.DB.GetCandles(context.Background(), "bitmex", "BTX", []string{"20200101"}, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, dbcandles.Market, "bitmex")
	assert.Equal(t, dbcandles.Symbol, "BTX")
	// assert.Equal(t, dbcandles.Tag, "20200101")
	assert.Equal(t, len(dbcandles.Candles), 100)
	for i := 0; i < 100; i++ {
		assert.Equal(t, dbcandles.Candles[i].Ts, int64(i))
		assert.Equal(t, dbcandles.Candles[i].Open, int64(100+i))
	}

	replygetcandles, err = client1.GetCandles(context.Background(), "bitmex", "BTX", []string{"20200101"}, 0, 0, nil)
	assert.NoError(t, err)
	assert.NotNil(t, replygetcandles)
	assert.Equal(t, replygetcandles.Market, "bitmex")
	assert.Equal(t, replygetcandles.Symbol, "BTX")
	// assert.Equal(t, replygetcandles.Tag, "20200101")
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

	replygetcandles, err = client1.GetCandles(context.Background(), "bitmex", "BTX", []string{"20200102"}, 0, 0, nil)
	assert.NoError(t, err)
	assert.NotNil(t, replygetcandles)
	assert.Equal(t, replygetcandles.Market, "bitmex")
	assert.Equal(t, replygetcandles.Symbol, "BTX")
	// assert.Equal(t, replygetcandles.Tag, "20200102")
	assert.Equal(t, len(replygetcandles.Candles), 100)

	for i := 0; i < 100; i++ {
		assert.Equal(t, replygetcandles.Candles[i].Ts, int64(i))
		assert.Equal(t, replygetcandles.Candles[i].Open, int64(100+i))
	}

	replygetcandles, err = client1.GetCandles(context.Background(), "bitmex", "BTX", []string{"20200102"}, 0, 90, nil)
	assert.NoError(t, err)
	assert.NotNil(t, replygetcandles)
	assert.Equal(t, replygetcandles.Market, "bitmex")
	assert.Equal(t, replygetcandles.Symbol, "BTX")
	// assert.Equal(t, replygetcandles.Tag, "20200102")
	assert.Equal(t, len(replygetcandles.Candles), 91)

	for i := 0; i < 91; i++ {
		assert.Equal(t, replygetcandles.Candles[i].Ts, int64(i))
		assert.Equal(t, replygetcandles.Candles[i].Open, int64(100+i))
	}

	replygetcandles, err = client1.GetCandles(context.Background(), "bitmex", "BTX", []string{"20200102"}, 10, 11, nil)
	assert.NoError(t, err)
	assert.NotNil(t, replygetcandles)
	assert.Equal(t, replygetcandles.Market, "bitmex")
	assert.Equal(t, replygetcandles.Symbol, "BTX")
	// assert.Equal(t, replygetcandles.Tag, "20200102")
	assert.Equal(t, len(replygetcandles.Candles), 2)

	for i := 0; i < 2; i++ {
		assert.Equal(t, replygetcandles.Candles[i].Ts, int64(10+i))
		assert.Equal(t, replygetcandles.Candles[i].Open, int64(110+i))
	}

	replygetcandles, err = client1.GetCandles(context.Background(), "bitmex", "BTX", []string{"20200102"}, 10, 10, nil)
	assert.NoError(t, err)
	assert.NotNil(t, replygetcandles)
	assert.Equal(t, replygetcandles.Market, "bitmex")
	assert.Equal(t, replygetcandles.Symbol, "BTX")
	// assert.Equal(t, replygetcandles.Tag, "20200102")
	assert.Equal(t, len(replygetcandles.Candles), 1)

	for i := 0; i < 1; i++ {
		assert.Equal(t, replygetcandles.Candles[i].Ts, int64(10+i))
		assert.Equal(t, replygetcandles.Candles[i].Open, int64(110+i))
	}

	replygetcandles, err = client1.GetCandles(context.Background(), "bitmex", "BTX", []string{"20200102"}, 10, 9, nil)
	assert.NoError(t, err)
	assert.NotNil(t, replygetcandles)
	assert.Equal(t, replygetcandles.Market, "bitmex")
	assert.Equal(t, replygetcandles.Symbol, "BTX")
	// assert.Equal(t, replygetcandles.Tag, "20200102")
	assert.Equal(t, len(replygetcandles.Candles), 0)

	// 0 < 200
	candles2 := &tradingdb2pb.Candles{
		Market: "bitmex",
		Symbol: "BTX",
		Tag:    "20200103",
	}

	ret, err = client1.UpdCandles(context.Background(), candles2, 200, nil)
	assert.NoError(t, err)
	assert.Equal(t, ret.LengthOK, int32(0))

	replygetcandles, err = client1.GetCandles(context.Background(), "bitmex", "BTX", []string{"20200103"}, 0, 0, nil)
	assert.NoError(t, err)
	assert.NotNil(t, replygetcandles)
	assert.Equal(t, replygetcandles.Market, "bitmex")
	assert.Equal(t, replygetcandles.Symbol, "BTX")
	// assert.Equal(t, replygetcandles.Tag, "20200103")
	assert.Equal(t, len(replygetcandles.Candles), 0)

	si := &tradingdb2pb.SymbolInfo{
		Market: "cnfund",
		Symbol: "240001",
		Fund: &tradingdb2pb.Fund{
			Code:       "240001",
			Name:       "华宝宝康消费品",
			Tags:       []string{"开放式", "混合型", "高风险"},
			CreateTime: int64(123),
			Size: []*tradingdb2pb.FundSize{&tradingdb2pb.FundSize{
				Size: 1.23,
				Time: 123,
			}},
			Company: "华宝基金管理有限公司",
			// Managers   []*FundManager `protobuf:"bytes,7,rep,name=managers,proto3" json:"managers,omitempty"`
			// Results    []*FundResult  `protobuf:"bytes,8,rep,name=results,proto3" json:"results,omitempty"`
		},
	}

	replyupdsymbol, err := client1.UpdSymbol(context.Background(), si, nil)
	assert.NoError(t, err)
	assert.NotNil(t, replyupdsymbol)
	assert.Equal(t, replyupdsymbol.IsOK, true)

	replygetsymbol, err := client1.GetSymbol(context.Background(), "cnfund", "240001", nil)
	assert.NoError(t, err)
	assert.NotNil(t, replygetsymbol)
	assert.Equal(t, replygetsymbol.Market, "cnfund")
	assert.Equal(t, replygetsymbol.Symbol, "240001")
	assert.NotNil(t, replygetsymbol.Fund)
	assert.Equal(t, replygetsymbol.Fund.Code, "240001")
	assert.Equal(t, replygetsymbol.Fund.Name, "华宝宝康消费品")
	assert.Equal(t, len(replygetsymbol.Fund.Tags), 3)
	assert.Equal(t, replygetsymbol.Fund.Tags[0], "开放式")
	assert.Equal(t, replygetsymbol.Fund.Tags[1], "混合型")
	assert.Equal(t, replygetsymbol.Fund.Tags[2], "高风险")
	assert.Equal(t, replygetsymbol.Fund.CreateTime, int64(123))
	assert.Equal(t, len(replygetsymbol.Fund.Size), 1)
	assert.Equal(t, replygetsymbol.Fund.Size[0].Size, float32(1.23))
	assert.Equal(t, replygetsymbol.Fund.Size[0].Time, int64(123))
	assert.Equal(t, replygetsymbol.Fund.Company, "华宝基金管理有限公司")

	nums := 0
	client1.GetSymbols(context.Background(), "cnfund", nil, func(si *tradingdb2pb.SymbolInfo) {
		nums++

		assert.Equal(t, si.Fund.Code, "240001")
		assert.Equal(t, si.Fund.Name, "华宝宝康消费品")
		assert.Equal(t, len(si.Fund.Tags), 3)
		assert.Equal(t, si.Fund.Tags[0], "开放式")
		assert.Equal(t, si.Fund.Tags[1], "混合型")
		assert.Equal(t, si.Fund.Tags[2], "高风险")
		assert.Equal(t, si.Fund.CreateTime, int64(123))
		assert.Equal(t, len(si.Fund.Size), 1)
		assert.Equal(t, si.Fund.Size[0].Size, float32(1.23))
		assert.Equal(t, si.Fund.Size[0].Time, int64(123))
		assert.Equal(t, si.Fund.Company, "华宝基金管理有限公司")
	}, nil)

	assert.Equal(t, nums, 1)

	serv.Stop()

	t.Logf("Test_Serv OK")
}
