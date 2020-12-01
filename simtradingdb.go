package tradingdb2

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	ankadb "github.com/zhs007/ankadb"
	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
)

const simtradingDBName = "simtrading"
const simtradingKeyPrefix = "st:"

func makeSimTradingDBKey(strategy string, market string, symbol string, tsStart int64, tsEnd int64) string {
	return fmt.Sprintf("%s%s:%s:%s:%v:%v", simtradingKeyPrefix, strategy, market, symbol, tsStart, tsEnd)
}

// SimTradingDB - database
type SimTradingDB struct {
	AnkaDB ankadb.AnkaDB
}

// NewSimTradingDB - new SimTradingDB
func NewSimTradingDB(dbpath string, httpAddr string, engine string) (*SimTradingDB, error) {
	cfg := ankadb.NewConfig()

	cfg.AddrHTTP = httpAddr
	cfg.PathDBRoot = dbpath
	cfg.ListDB = append(cfg.ListDB, ankadb.DBConfig{
		Name:   simtradingDBName,
		Engine: engine,
		PathDB: simtradingDBName,
	})

	ankaDB, err := ankadb.NewAnkaDB(cfg, nil)
	if ankaDB == nil {
		return nil, err
	}

	db := &SimTradingDB{
		AnkaDB: ankaDB,
	}

	return db, err
}

// UpdSimTrading - update simulation trading
func (db *DB) UpdSimTrading(ctx context.Context, params *tradingpb.SimTradingParams, pnldata *tradingpb.PNLData) error {
	buf, err := proto.Marshal(pnldata)
	if err != nil {
		return err
	}

	err = db.AnkaDB.Set(ctx, dbname,
		makeSimTradingDBKey(params.Strategies[0].Name, params.Assets[0].Market, params.Assets[0].Code, params.StartTs, params.EndTs), buf)
	if err != nil {
		return err
	}

	return nil
}

// GetSimTrading - get candles
func (db *DB) GetSimTrading(ctx context.Context, market string, symbol string, tags []string, tsStart int64, tsEnd int64) (
	*tradingpb.Candles, error) {

	if market == "" {
		return nil, ErrInvalidMarket
	}

	if symbol == "" {
		return nil, ErrInvalidSymbol
	}

	candles := &tradingpb.Candles{
		Market: market,
		Symbol: symbol,
	}

	if tsStart > 0 && tsEnd <= 0 {
		tsEnd = time.Now().Unix()
	}

	err := db.AnkaDB.ForEachWithPrefix(ctx, dbname, makeCandlesDBKeyPrefix(market, symbol), func(key string, buf []byte) error {
		cc := &tradingpb.Candles{}

		err := proto.Unmarshal(buf, cc)
		if err != nil {
			return err
		}

		if len(tags) == 0 || tradingdb2utils.IndexOfStringSlice(tags, cc.Tag, 0) >= 0 {
			if tsStart > 0 || tsEnd > 0 {
				for _, v := range cc.Candles {
					if v.Ts >= tsStart && v.Ts <= tsEnd {
						candles.Candles = append(candles.Candles, v)
					}
				}
			} else {
				candles.Candles = append(candles.Candles, cc.Candles...)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// if len(candles.Candles) == 0 {
	// 	return nil, nil
	// }

	return candles, nil
}
