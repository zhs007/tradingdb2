package tradingdb2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
)

func Test_SortCandles(t *testing.T) {
	candles := &tradingpb.Candles{
		Candles: []*tradingpb.Candle{
			&tradingpb.Candle{
				Ts: 1,
			},
			&tradingpb.Candle{
				Ts: 3,
			},
			&tradingpb.Candle{
				Ts: 2,
			},
		},
	}

	SortCandles(candles)

	assert.Equal(t, len(candles.Candles), 3)
	assert.Equal(t, candles.Candles[0].Ts, int64(1))
	assert.Equal(t, candles.Candles[1].Ts, int64(2))
	assert.Equal(t, candles.Candles[2].Ts, int64(3))

	candles1 := &tradingpb.Candles{}
	SortCandles(candles1)
	assert.Equal(t, len(candles1.Candles), 0)

	t.Logf("Test_SortCandles OK")
}

func Test_InsCandles(t *testing.T) {
	candles := &tradingpb.Candles{
		Candles: []*tradingpb.Candle{
			&tradingpb.Candle{
				Ts: 10,
			},
			&tradingpb.Candle{
				Ts: 3,
			},
			&tradingpb.Candle{
				Ts: 2,
			},
		},
	}

	SortCandles(candles)

	assert.Equal(t, len(candles.Candles), 3)
	assert.Equal(t, candles.Candles[0].Ts, int64(2))
	assert.Equal(t, candles.Candles[1].Ts, int64(3))
	assert.Equal(t, candles.Candles[2].Ts, int64(10))

	InsCandles(candles, &tradingpb.Candle{
		Ts: 10,
	})

	assert.Equal(t, len(candles.Candles), 3)
	assert.Equal(t, candles.Candles[0].Ts, int64(2))
	assert.Equal(t, candles.Candles[1].Ts, int64(3))
	assert.Equal(t, candles.Candles[2].Ts, int64(10))

	InsCandles(candles, &tradingpb.Candle{
		Ts: 11,
	})

	assert.Equal(t, len(candles.Candles), 4)
	assert.Equal(t, candles.Candles[0].Ts, int64(2))
	assert.Equal(t, candles.Candles[1].Ts, int64(3))
	assert.Equal(t, candles.Candles[2].Ts, int64(10))
	assert.Equal(t, candles.Candles[3].Ts, int64(11))

	InsCandles(candles, &tradingpb.Candle{
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
	candles := &tradingpb.Candles{
		Candles: []*tradingpb.Candle{
			&tradingpb.Candle{
				Ts: 10,
			},
			&tradingpb.Candle{
				Ts: 3,
			},
			&tradingpb.Candle{
				Ts: 2,
			},
		},
	}

	SortCandles(candles)

	assert.Equal(t, len(candles.Candles), 3)
	assert.Equal(t, candles.Candles[0].Ts, int64(2))
	assert.Equal(t, candles.Candles[1].Ts, int64(3))
	assert.Equal(t, candles.Candles[2].Ts, int64(10))

	candles1 := &tradingpb.Candles{
		Candles: []*tradingpb.Candle{
			&tradingpb.Candle{
				Ts: 10,
			},
			&tradingpb.Candle{
				Ts: 11,
			},
			&tradingpb.Candle{
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
	candles := &tradingpb.Candles{}

	for i := 0; i < 100; i++ {
		InsCandles(candles, &tradingpb.Candle{
			Ts: int64(i),
		})
	}

	start := int64(0)
	lastnums := 100
	BatchCandles(candles, 30, func(lst *tradingpb.Candles) error {
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
