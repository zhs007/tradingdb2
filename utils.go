package tradingdb2

import (
	"math"
	"sort"

	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
)

// SortCandles - sort Candles
func SortCandles(candles *tradingpb.Candles) {
	sort.Slice(candles.Candles, func(i, j int) bool {
		return candles.Candles[i].Ts < candles.Candles[j].Ts
	})
}

// InsCandles - insert a candle into candles, sort by ts
func InsCandles(candles *tradingpb.Candles, candle *tradingpb.Candle) {
	for i, v := range candles.Candles {
		if v.Ts == candle.Ts {
			return
		}

		if candle.Ts < v.Ts {
			candles.Candles = append(candles.Candles[:i], append([]*tradingpb.Candle{candle}, candles.Candles[i:]...)...)

			return
		}
	}

	candles.Candles = append(candles.Candles, candle)
}

// MergeCandles - merge Candles, candles is a sorted candles
func MergeCandles(candles *tradingpb.Candles, src *tradingpb.Candles) {
	for _, v := range src.Candles {
		InsCandles(candles, v)
	}
}

// FuncOnBatchCandles - used in BatchCandles
type FuncOnBatchCandles func(candles *tradingpb.Candles) error

// BatchCandles - batch candles
func BatchCandles(candles *tradingpb.Candles, nums int, onBatch FuncOnBatchCandles) error {
	if len(candles.Candles) == 0 {
		err := onBatch(&tradingpb.Candles{})
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

		err := onBatch(&tradingpb.Candles{
			Candles: candles.Candles[i : i+curlen],
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func insFundSize(fund *tradingpb.Fund, fs *tradingpb.FundSize) {
	if len(fund.Size) == 0 {
		fund.Size = append(fund.Size, fs)

		return
	}

	cv := fund.Size[len(fund.Size)-1]
	if fs.Time > cv.Time && math.Abs(float64(cv.Size-fs.Size)) >= 0.1 {
		fund.Size = append(fund.Size, fs)
	}
}

func insFundResult(fund *tradingpb.Fund, fr *tradingpb.FundResult) {
	for _, v := range fund.Results {
		if v.Name == fr.Name && v.StartTime == fr.StartTime && v.EndTime == fr.EndTime {
			return
		}
	}

	fund.Results = append(fund.Results, fr)
}

// MergeFund - merge fund
func MergeFund(fund0 *tradingpb.Fund, fund1 *tradingpb.Fund) *tradingpb.Fund {
	if fund1.Code != "" {
		fund0.Code = fund1.Code
	}

	if fund1.Name != "" {
		fund0.Name = fund1.Name
	}

	if fund1.Tags != nil {
		fund0.Tags = fund1.Tags
	}

	if fund1.CreateTime > 0 {
		fund0.CreateTime = fund1.CreateTime
	}

	for _, fs := range fund1.Size {
		insFundSize(fund0, fs)
	}

	if fund1.Company != "" {
		fund0.Company = fund1.Company
	}

	if fund1.Managers != nil {
		fund0.Managers = fund1.Managers
	}

	for _, fr := range fund1.Results {
		insFundResult(fund0, fr)
	}

	return fund0
}

// FixFundResult - fix
func FixFundResult(fr *tradingpb.FundResult) {
	if math.IsNaN(float64(fr.MaxDrawdown)) {
		fr.MaxDrawdown = 0
	}

	if math.IsNaN(float64(fr.Sharpe)) {
		fr.Sharpe = 0
	}

	if math.IsNaN(float64(fr.AnnualizedReturns)) {
		fr.AnnualizedReturns = 0
	}

	if math.IsNaN(float64(fr.AnnualizedVolatility)) {
		fr.AnnualizedVolatility = 0
	}

	if math.IsNaN(float64(fr.TotalReturns)) {
		fr.TotalReturns = 0
	}
}

// FixFundManager - fix
func FixFundManager(fm *tradingpb.FundManager) {
	for _, v := range fm.Results {
		FixFundResult(v)
	}
}

// FixFundManagers - fix
func FixFundManagers(fms []*tradingpb.FundManager) {
	for _, v := range fms {
		FixFundManager(v)
	}
}
