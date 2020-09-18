package tradingdb2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tradingdb2pb "github.com/zhs007/tradingdb2/pb"
)

func Test_SortCandles(t *testing.T) {
	candles := &tradingdb2pb.Candles{
		Candles: []*tradingdb2pb.Candle{
			&tradingdb2pb.Candle{
				Ts: 1,
			},
			&tradingdb2pb.Candle{
				Ts: 3,
			},
			&tradingdb2pb.Candle{
				Ts: 2,
			},
		},
	}

	SortCandles(candles)

	assert.Equal(t, len(candles.Candles), 3)
	assert.Equal(t, candles.Candles[0].Ts, int64(1))
	assert.Equal(t, candles.Candles[1].Ts, int64(2))
	assert.Equal(t, candles.Candles[2].Ts, int64(3))

	candles1 := &tradingdb2pb.Candles{}
	SortCandles(candles1)
	assert.Equal(t, len(candles1.Candles), 0)

	t.Logf("Test_SortCandles OK")
}

func Test_InsCandles(t *testing.T) {
	candles := &tradingdb2pb.Candles{
		Candles: []*tradingdb2pb.Candle{
			&tradingdb2pb.Candle{
				Ts: 10,
			},
			&tradingdb2pb.Candle{
				Ts: 3,
			},
			&tradingdb2pb.Candle{
				Ts: 2,
			},
		},
	}

	SortCandles(candles)

	assert.Equal(t, len(candles.Candles), 3)
	assert.Equal(t, candles.Candles[0].Ts, int64(2))
	assert.Equal(t, candles.Candles[1].Ts, int64(3))
	assert.Equal(t, candles.Candles[2].Ts, int64(10))

	InsCandles(candles, &tradingdb2pb.Candle{
		Ts: 10,
	})

	assert.Equal(t, len(candles.Candles), 3)
	assert.Equal(t, candles.Candles[0].Ts, int64(2))
	assert.Equal(t, candles.Candles[1].Ts, int64(3))
	assert.Equal(t, candles.Candles[2].Ts, int64(10))

	InsCandles(candles, &tradingdb2pb.Candle{
		Ts: 11,
	})

	assert.Equal(t, len(candles.Candles), 4)
	assert.Equal(t, candles.Candles[0].Ts, int64(2))
	assert.Equal(t, candles.Candles[1].Ts, int64(3))
	assert.Equal(t, candles.Candles[2].Ts, int64(10))
	assert.Equal(t, candles.Candles[3].Ts, int64(11))

	InsCandles(candles, &tradingdb2pb.Candle{
		Ts: 1,
	})

	assert.Equal(t, len(candles.Candles), 5)
	assert.Equal(t, candles.Candles[0].Ts, int64(1))
	assert.Equal(t, candles.Candles[1].Ts, int64(2))
	assert.Equal(t, candles.Candles[2].Ts, int64(3))
	assert.Equal(t, candles.Candles[3].Ts, int64(10))
	assert.Equal(t, candles.Candles[4].Ts, int64(11))

	t.Logf("Test_InsCandles OK")
}

func Test_MergeCandles(t *testing.T) {
	candles := &tradingdb2pb.Candles{
		Candles: []*tradingdb2pb.Candle{
			&tradingdb2pb.Candle{
				Ts: 10,
			},
			&tradingdb2pb.Candle{
				Ts: 3,
			},
			&tradingdb2pb.Candle{
				Ts: 2,
			},
		},
	}

	SortCandles(candles)

	assert.Equal(t, len(candles.Candles), 3)
	assert.Equal(t, candles.Candles[0].Ts, int64(2))
	assert.Equal(t, candles.Candles[1].Ts, int64(3))
	assert.Equal(t, candles.Candles[2].Ts, int64(10))

	candles1 := &tradingdb2pb.Candles{
		Candles: []*tradingdb2pb.Candle{
			&tradingdb2pb.Candle{
				Ts: 10,
			},
			&tradingdb2pb.Candle{
				Ts: 11,
			},
			&tradingdb2pb.Candle{
				Ts: 1,
			},
		},
	}

	MergeCandles(candles, candles1)

	assert.Equal(t, len(candles.Candles), 5)
	assert.Equal(t, candles.Candles[0].Ts, int64(1))
	assert.Equal(t, candles.Candles[1].Ts, int64(2))
	assert.Equal(t, candles.Candles[2].Ts, int64(3))
	assert.Equal(t, candles.Candles[3].Ts, int64(10))
	assert.Equal(t, candles.Candles[4].Ts, int64(11))

	t.Logf("Test_MergeCandles OK")
}

func Test_BatchCandles(t *testing.T) {
	candles := &tradingdb2pb.Candles{}

	for i := 0; i < 100; i++ {
		InsCandles(candles, &tradingdb2pb.Candle{
			Ts: int64(i),
		})
	}

	start := int64(0)
	lastnums := 100
	BatchCandles(candles, 30, func(lst *tradingdb2pb.Candles) error {
		if lastnums >= 30 {
			assert.Equal(t, len(lst.Candles), 30)
		} else {
			assert.Equal(t, len(lst.Candles), lastnums)
		}

		lastnums -= len(lst.Candles)

		for _, v := range lst.Candles {
			assert.Equal(t, v.Ts, start)
			start++
		}

		return nil
	})

	t.Logf("Test_BatchCandles OK")
}
