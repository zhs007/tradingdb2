package tradingdb2

import (
	"sort"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	excel "github.com/zhs007/adacore/excel"
)

// alldataMember - alldata
var alldataMember = []string{
	"market",
	"symbol",
	"tag",
	"candles",
}

type alldata struct {
	market  string
	symbol  string
	tag     string
	candles string
}

func genAllData(root *TreeMapNode) []alldata {
	lst := []alldata{}

	for _, market := range root.Children {
		for _, symbol := range market.Children {
			for _, tag := range symbol.Children {
				for k := range tag.Children {
					lst = append(lst, alldata{
						market:  market.Name,
						symbol:  symbol.Name,
						tag:     tag.Name,
						candles: k,
					})
				}
			}
		}
	}

	sort.Slice(lst, func(i, j int) bool {
		if lst[i].market != lst[j].market {
			return lst[i].market < lst[j].market
		}

		if lst[i].symbol != lst[j].symbol {
			return lst[i].symbol < lst[j].symbol
		}

		return lst[i].tag < lst[j].tag
	})

	return lst
}

// ExpAllData - export all data
func ExpAllData(fn string, root *TreeMapNode) error {
	f := excelize.NewFile()

	lst := genAllData(root)

	// write head
	excel.SetSheet(f, "Sheet1", 1, 1, alldataMember, len(lst),
		func(i int, member string) string {
			return ""
		},
		func(i int, member string) (interface{}, error) {
			v := lst[i]
			if member == "market" {
				return v.market, nil
			} else if member == "symbol" {
				return v.symbol, nil
			} else if member == "tag" {
				return v.tag, nil
			} else if member == "candles" {
				return v.candles, nil
			}

			return nil, nil
		})

	f.SaveAs(fn)

	return nil
}
