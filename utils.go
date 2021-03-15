package tradingdb2

import (
	"crypto/sha1"
	"fmt"
	"math"
	"sort"

	"google.golang.org/protobuf/proto"
	// proto "github.com/golang/protobuf/proto"
	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
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

// // IsSameStrategy - is same Strategy
// func IsSameStrategy(v0 *tradingpb.Strategy, v1 *tradingpb.Strategy) bool {
// 	if v0.Name != v1.Name {
// 		return false
// 	}

// 	if !proto.Equal(v0.Asset, v1.Asset) {
// 		return false
// 	}

// 	if len(v0.Buy) != len(v1.Buy) {
// 		return false
// 	}

// 	for bi, bv0 := range v0.Buy {
// 		bv1 := v1.Buy[bi]

// 		if !proto.Equal(bv0, bv1) {
// 			return false
// 		}
// 	}

// 	if len(v0.Sell) != len(v1.Sell) {
// 		return false
// 	}

// 	for si, sv0 := range v0.Sell {
// 		sv1 := v1.Sell[si]

// 		if !proto.Equal(sv0, sv1) {
// 			return false
// 		}
// 	}

// 	if len(v0.Stoploss) != len(v1.Stoploss) {
// 		return false
// 	}

// 	for si, sv0 := range v0.Stoploss {
// 		sv1 := v1.Stoploss[si]

// 		if !proto.Equal(sv0, sv1) {
// 			return false
// 		}
// 	}

// 	if len(v0.Takeprofit) != len(v1.Takeprofit) {
// 		return false
// 	}

// 	for ti, tv0 := range v0.Takeprofit {
// 		tv1 := v1.Takeprofit[ti]

// 		if !proto.Equal(tv0, tv1) {
// 			return false
// 		}
// 	}

// 	if !proto.Equal(v0.ParamsBuy, v1.ParamsBuy) {
// 		return false
// 	}

// 	if !proto.Equal(v0.ParamsSell, v1.ParamsSell) {
// 		return false
// 	}

// 	if !proto.Equal(v0.ParamsStopLoss, v1.ParamsStopLoss) {
// 		return false
// 	}

// 	if !proto.Equal(v0.ParamsTakeProfit, v1.ParamsTakeProfit) {
// 		return false
// 	}

// 	if !proto.Equal(v0.ParamsInit, v1.ParamsInit) {
// 		return false
// 	}

// 	if !proto.Equal(v0.ParamsAIP, v1.ParamsAIP) {
// 		return false
// 	}

// 	if len(v0.Indicators) != len(v1.Indicators) {
// 		return false
// 	}

// 	// 这里顺序不一致也没问题
// 	for _, iv0 := range v0.Indicators {
// 		if tradingdb2utils.IndexOfStringSlice(v1.Indicators, iv0, 0) < 0 {
// 			return false
// 		}
// 	}

// 	if !proto.Equal(v0.FeeBuy, v1.FeeBuy) {
// 		return false
// 	}

// 	if !proto.Equal(v0.FeeSell, v1.FeeSell) {
// 		return false
// 	}

// 	return true
// }

// cmpCtrlCondition - group name vals operators strvals combCondition
func cmpCtrlCondition(cc0 *tradingpb.CtrlCondition, cc1 *tradingpb.CtrlCondition) int {
	if cc0.Group != cc1.Group {
		if cc0.Group < cc1.Group {
			return -1
		}

		return 1
	}

	if cc0.Name != cc1.Name {
		if cc0.Name < cc1.Name {
			return -1
		}

		return 1
	}

	if len(cc0.Vals) != len(cc1.Vals) {
		if len(cc0.Vals) < len(cc1.Vals) {
			return -1
		}

		return 1
	}

	if len(cc0.Vals) > 0 {
		for ti := 0; ti < len(cc0.Vals); ti++ {
			if cc0.Vals[ti] != cc1.Vals[ti] {
				if cc0.Vals[ti] < cc1.Vals[ti] {
					return -1
				}

				return 1
			}
		}
	}

	if len(cc0.Operators) != len(cc1.Operators) {
		if len(cc0.Operators) < len(cc1.Operators) {
			return -1
		}

		return 1
	}

	if len(cc0.Operators) > 0 {
		for ti := 0; ti < len(cc0.Operators); ti++ {
			if cc0.Operators[ti] != cc1.Operators[ti] {
				if cc0.Operators[ti] < cc1.Operators[ti] {
					return -1
				}

				return 1
			}
		}
	}

	if len(cc0.StrVals) != len(cc1.StrVals) {
		if len(cc0.StrVals) < len(cc1.StrVals) {
			return -1
		}

		return 1
	}

	if len(cc0.StrVals) > 0 {
		for ti := 0; ti < len(cc0.StrVals); ti++ {
			if cc0.StrVals[ti] != cc1.StrVals[ti] {
				if cc0.StrVals[ti] < cc1.StrVals[ti] {
					return -1
				}

				return 1
			}
		}
	}

	if cc0.CombCondition < cc1.CombCondition {
		return -1
	}

	if cc0.CombCondition > cc1.CombCondition {
		return 1
	}

	return 0
}

// sortCtrlConditions - group name vals operators strvals combCondition
func sortCtrlConditions(ctrlConditions []*tradingpb.CtrlCondition) error {
	sort.Slice(ctrlConditions, func(i, j int) bool {
		return cmpCtrlCondition(ctrlConditions[i], ctrlConditions[j]) <= 0
	})

	return nil
}

// cmpAsset -
func cmpAsset(a0 *tradingpb.Asset, a1 *tradingpb.Asset) int {
	if a0.Market != a1.Market {
		if a0.Market < a1.Market {
			return -1
		}

		return 1
	}

	if a0.Code != a1.Code {
		if a0.Code < a1.Code {
			return -1
		}

		return 1
	}

	if len(a0.Tags) != len(a1.Tags) {
		if len(a0.Tags) < len(a1.Tags) {
			return -1
		}

		return 1
	}

	if len(a0.Tags) > 0 {
		for ti := 0; ti < len(a0.Tags); ti++ {
			if a0.Tags[ti] != a1.Tags[ti] {
				if a0.Tags[ti] < a1.Tags[ti] {
					return -1
				}

				return 1
			}
		}
	}

	return 0
}

// cmpStrategy - name asset buy sell stoploss takeprofit hash
func cmpStrategy(s0 *tradingpb.Strategy, s1 *tradingpb.Strategy) (int, error) {
	if s0.Name != s1.Name {
		if s0.Name < s1.Name {
			return -1, nil
		}

		return 1, nil
	}

	ca := cmpAsset(s0.Asset, s1.Asset)
	if ca != 0 {
		return ca, nil
	}

	if len(s0.Buy) != len(s1.Buy) {
		if len(s0.Buy) < len(s1.Buy) {
			return -1, nil
		}

		return 1, nil
	}

	if len(s0.Buy) > 0 {
		for ti := 0; ti < len(s0.Buy); ti++ {
			cb := cmpCtrlCondition(s0.Buy[ti], s1.Buy[ti])
			if cb != 0 {
				return cb, nil
			}
		}
	}

	if len(s0.Sell) != len(s1.Sell) {
		if len(s0.Sell) < len(s1.Sell) {
			return -1, nil
		}

		return 1, nil
	}

	if len(s0.Sell) > 0 {
		for ti := 0; ti < len(s0.Sell); ti++ {
			cs := cmpCtrlCondition(s0.Sell[ti], s1.Sell[ti])
			if cs != 0 {
				return cs, nil
			}
		}
	}

	if len(s0.Stoploss) != len(s1.Stoploss) {
		if len(s0.Stoploss) < len(s1.Stoploss) {
			return -1, nil
		}

		return 1, nil
	}

	if len(s0.Stoploss) > 0 {
		for ti := 0; ti < len(s0.Stoploss); ti++ {
			cs := cmpCtrlCondition(s0.Stoploss[ti], s1.Stoploss[ti])
			if cs != 0 {
				return cs, nil
			}
		}
	}

	if len(s0.Takeprofit) != len(s1.Takeprofit) {
		if len(s0.Takeprofit) < len(s1.Takeprofit) {
			return -1, nil
		}

		return 1, nil
	}

	if len(s0.Takeprofit) > 0 {
		for ti := 0; ti < len(s0.Takeprofit); ti++ {
			ct := cmpCtrlCondition(s0.Takeprofit[ti], s1.Takeprofit[ti])
			if ct != 0 {
				return ct, nil
			}
		}
	}

	buf0, err := proto.Marshal(s0)
	if err != nil {
		return -1, err
	}

	buf1, err := proto.Marshal(s1)
	if err != nil {
		return -1, err
	}

	if len(buf0) != len(buf1) {
		if len(buf0) < len(buf1) {
			return -1, nil
		}

		return 1, nil
	}

	for ti := 0; ti < len(buf0); ti++ {
		if buf0[ti] != buf1[ti] {
			if buf0[ti] < buf1[ti] {
				return -1, nil
			}

			return 0, nil
		}
	}

	return 0, nil
}

func sortStrategies(strategies []*tradingpb.Strategy) error {
	for _, v := range strategies {
		err := sortCtrlConditions(v.Buy)
		if err != nil {
			tradingdb2utils.Warn("sortStrategys:sortCtrlConditions Buy",
				zap.Error(err))

			return err
		}

		err = sortCtrlConditions(v.Sell)
		if err != nil {
			tradingdb2utils.Warn("sortStrategys:sortCtrlConditions Sell",
				zap.Error(err))

			return err
		}

		err = sortCtrlConditions(v.Stoploss)
		if err != nil {
			tradingdb2utils.Warn("sortStrategys:sortCtrlConditions Stoploss",
				zap.Error(err))

			return err
		}

		err = sortCtrlConditions(v.Takeprofit)
		if err != nil {
			tradingdb2utils.Warn("sortStrategys:sortCtrlConditions Takeprofit",
				zap.Error(err))

			return err
		}
	}

	sort.Slice(strategies, func(i, j int) bool {
		cs, err := cmpStrategy(strategies[i], strategies[j])
		if err != nil {
			tradingdb2utils.Warn("sortStrategys:cmpStrategy",
				zap.Error(err))
		}

		return cs <= 0
	})

	return nil
}

func rebuildSimTradingParams(params *tradingpb.SimTradingParams) (*tradingpb.SimTradingParams, []byte, string, error) {
	nmsg := proto.Clone(params)
	np, isok := nmsg.(*tradingpb.SimTradingParams)
	if !isok {
		tradingdb2utils.Warn("rebuildSimTradingParams:Clone",
			zap.Error(ErrInvalidSimTradingParams))

		return nil, nil, "", ErrInvalidSimTradingParams
	}

	np.Title = ""

	err := sortStrategies(np.Strategies)
	if err != nil {
		tradingdb2utils.Warn("rebuildSimTradingParams:sortStrategys",
			zap.Error(err))

		return nil, nil, "", err
	}

	buf, err := proto.Marshal(np)
	if err != nil {
		tradingdb2utils.Warn("rebuildSimTradingParams:Marshal",
			zap.Error(err))

		return nil, nil, "", err
	}

	h := sha1.New()
	h.Write(buf)
	bs := h.Sum(nil)

	return np, buf, fmt.Sprintf("%x", bs), nil
}
