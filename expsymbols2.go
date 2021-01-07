package tradingdb2

import (
	"context"
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

// ExpSymbols2 - export symbols
func ExpSymbols2(ctx context.Context, fn string, db2 *DB2, market string) error {
	f := excelize.NewFile()

	symbols, err := db2.GetMarketSymbols(ctx, market)
	if err != nil {
		return err
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
					tradingdb2utils.Warn("ExpSymbols:getSymbol Size",
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
						tradingdb2utils.Warn("ExpSymbols:Marshal Size",
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
						tradingdb2utils.Warn("ExpSymbols:Marshal Managers",
							zap.Error(err))

						return nil, nil
					}

					return string(b), nil
				}

				return "", nil
			} else if member == "values" {
				if curSymbol.Fund != nil {
					fv, err := db2.GetCandles(ctx, curSymbol.Market, curSymbol.Symbol, 0, 0)
					if err != nil {
						tradingdb2utils.Warn("ExpSymbols:Marshal GetCandles",
							zap.Error(err))

						return nil, nil
					}

					if fv != nil {
						return strconv.Itoa(len(fv.Candles)), nil
					}

					return "0", nil
				}

				return "0", nil
			}

			return nil, nil
		})

	f.SaveAs(fn)

	return nil
}
