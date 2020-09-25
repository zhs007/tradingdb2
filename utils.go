package tradingdb2

import (
	"sort"

	tradingdb2pb "github.com/zhs007/tradingdb2/tradingdb2pb"
)

// SortCandles - sort Candles
func SortCandles(candles *tradingdb2pb.Candles) {
	sort.Slice(candles.Candles, func(i, j int) bool {
		return candles.Candles[i].Ts < candles.Candles[j].Ts
	})
}

// InsCandles - insert a candle into candles, sort by ts
func InsCandles(candles *tradingdb2pb.Candles, candle *tradingdb2pb.Candle) {
	for i, v := range candles.Candles {
		if v.Ts == candle.Ts {
			return
		}

		if candle.Ts < v.Ts {
			candles.Candles = append(candles.Candles[:i], append([]*tradingdb2pb.Candle{candle}, candles.Candles[i:]...)...)

			return
		}
	}

	candles.Candles = append(candles.Candles, candle)
}

// MergeCandles - merge Candles, candles is a sorted candles
func MergeCandles(candles *tradingdb2pb.Candles, src *tradingdb2pb.Candles) {
	for _, v := range src.Candles {
		InsCandles(candles, v)
	}
}

// FuncOnBatchCandles - used in BatchCandles
type FuncOnBatchCandles func(candles *tradingdb2pb.Candles) error

// BatchCandles - batch candles
func BatchCandles(candles *tradingdb2pb.Candles, nums int, onBatch FuncOnBatchCandles) error {
	if len(candles.Candles) == 0 {
		err := onBatch(&tradingdb2pb.Candles{})
		if err != nil {
			return err
		}

		return nil
	}

	curlen := nums
	for i := 0; i < len(candles.Candles); i += curlen {
		if i+nums > len(candles.Candles) {
			curlen = len(candles.Candles) - i
		}

		err := onBatch(&tradingdb2pb.Candles{
			Candles: candles.Candles[i : i+curlen],
		})
		if err != nil {
			return err
		}
	}

	return nil
}
