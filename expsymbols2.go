package tradingdb2

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	jsoniter "github.com/json-iterator/go"
	excel "github.com/zhs007/adacore/excel"
	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
)

// fundsMember - fund member
var fundsMember2 = []string{
	"code",
	"name",
	"fullname",
	"tags",
	"createtime",
	"size",
	"company",
	"managers",
	"values",
	"type",
	"avgvolume30",
}

func getSymbol2(ctx context.Context, db2 *DB2, market string, symbol string) (*tradingpb.SymbolInfo, error) {
	curSymbol, err := db2.GetSymbol(ctx, market, symbol)
	if err != nil {
		return nil, err
	}

	if curSymbol == nil {
		curSymbol = &tradingpb.SymbolInfo{
			Market: market,
			Symbol: symbol,
		}
	}

	return curSymbol, nil
}

func calcAvgVolume30(candles *tradingpb.Candles) float32 {
	var totalvolume int64
	nums := 0

	for i := len(candles.Candles) - 1; i >= 0 && i > len(candles.Candles)-31; i-- {
		totalvolume += candles.Candles[i].Volume
		nums++
	}

	return float32(totalvolume) / float32(nums)
}

// getCandlesWithSymbolInfo - get candles
func getCandlesWithSymbolInfo(ctx context.Context, db2 *DB2, si *tradingpb.SymbolInfo) (*tradingpb.Candles, error) {
	if si.Market == "jqdata" {
		str := strings.ReplaceAll(si.Symbol, ".", "_")
		str += "|1d"

		return db2.GetCandles(ctx, si.Market, str, 0, 0, 0)
	}

	return db2.GetCandles(ctx, si.Market, si.Symbol, 0, 0, 0)
}

// ExpSymbols2 - export symbols
func ExpSymbols2(ctx context.Context, fn string, db2 *DB2, market string) error {
	f := excelize.NewFile()

	symbols, err := db2.GetMarketSymbols(ctx, market)
	if err != nil {
		tradingdb2utils.Warn("ExpSymbols2:GetMarketSymbols",
			zap.Error(err))

		return err
	}

	if len(symbols) <= 0 {
		tradingdb2utils.Warn("ExpSymbols2:GetMarketSymbols",
			zap.Error(ErrNoSymbols))

		return ErrNoSymbols
	}

	curSymbol, err := getSymbol2(ctx, db2, market, symbols[0])
	if err != nil {
		return err
	}

	// write head
	excel.SetSheet(f, "Sheet1", 1, 1, fundsMember2, len(symbols),
		func(i int, member string) string {
			return ""
		},
		func(i int, member string) (interface{}, error) {
			v := symbols[i]
			if v != curSymbol.Symbol {
				curSymbol, err = getSymbol2(ctx, db2, market, v)
				if err != nil {
					tradingdb2utils.Warn("ExpSymbols2:getSymbol Size",
						zap.Error(err))

					return nil, nil
				}
			}

			if member == "code" {
				return curSymbol.Symbol, nil
			} else if member == "name" {
				if curSymbol.Fund != nil {
					return curSymbol.Fund.Name, nil
				}

				return curSymbol.Name, nil
			} else if member == "fullname" {
				return curSymbol.Fullname, nil
			} else if member == "tags" {
				if curSymbol.Fund != nil {
					return strings.Join(curSymbol.Fund.Tags, "; "), nil
				}

				return "", nil
			} else if member == "createtime" {
				if curSymbol.Fund != nil {
					return time.Unix(curSymbol.Fund.CreateTime, 0).Format("2006-01-02"), nil
				}

				return "", nil
			} else if member == "size" {
				if curSymbol.Fund != nil {
					json := jsoniter.ConfigCompatibleWithStandardLibrary

					b, err := json.Marshal(curSymbol.Fund.Size)
					if err != nil {
						tradingdb2utils.Warn("ExpSymbols2:Marshal Size",
							zap.Error(err))

						return nil, nil
					}

					return string(b), nil
				}

				return "", nil
			} else if member == "company" {
				if curSymbol.Fund != nil {
					return curSymbol.Fund.Company, nil
				}

				return "", nil
			} else if member == "managers" {
				if curSymbol.Fund != nil {
					json := jsoniter.ConfigCompatibleWithStandardLibrary

					mgrs := curSymbol.Fund.Managers
					FixFundManagers(mgrs)

					b, err := json.Marshal(mgrs)
					if err != nil {
						tradingdb2utils.Warn("ExpSymbols2:Marshal Managers",
							zap.Error(err))

						return nil, nil
					}

					return string(b), nil
				}

				return "", nil
			} else if member == "values" {
				fv, err := getCandlesWithSymbolInfo(ctx, db2, curSymbol)
				if err != nil {
					tradingdb2utils.Warn("ExpSymbols2:getCandlesWithSymbolInfo values",
						zap.Error(err))

					return nil, nil
				}

				if fv != nil {
					return strconv.Itoa(len(fv.Candles)), nil
				}

				return "0", nil

			} else if member == "type" {
				return curSymbol.Type, nil
			} else if member == "avgvolume30" {
				fv, err := getCandlesWithSymbolInfo(ctx, db2, curSymbol)
				if err != nil {
					tradingdb2utils.Warn("ExpSymbols2:getCandlesWithSymbolInfo avgvolume30",
						zap.Error(err))

					return nil, nil
				}

				if fv != nil {
					return fmt.Sprintf("%f", calcAvgVolume30(fv)), nil
				}

				return "0", nil
			}

			return nil, nil
		})

	f.SaveAs(fn)

	return nil
}
