package tradingdb2

import (
	"context"
	"strconv"
	"time"

	ankadb "github.com/zhs007/ankadb"
	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const dbname = "tradingdb2"
const candlesKeyPrefix = "candles:"
const symbolKeyPrefix = "symbol:"

func makeCandlesDBKey(market string, symbol string, tag string) string {
	return tradingdb2utils.AppendString(candlesKeyPrefix, market, ":", symbol, ":", tag)
}

func makeCandlesDBKeyPrefix(market string, symbol string) string {
	return tradingdb2utils.AppendString(candlesKeyPrefix, market, ":", symbol, ":")
}

func makeSymbolDBKey(market string, symbol string) string {
	return tradingdb2utils.AppendString(symbolKeyPrefix, market, ":", symbol)
}

func makeSymbolDBKeyPrefix(market string) string {
	return tradingdb2utils.AppendString(symbolKeyPrefix, market, ":")
}

// DB - database
type DB struct {
	AnkaDB ankadb.AnkaDB
}

// NewDB - new DB
func NewDB(dbpath string, httpAddr string, engine string) (*DB, error) {
	cfg := ankadb.NewConfig()

	cfg.AddrHTTP = httpAddr
	cfg.PathDBRoot = dbpath
	cfg.ListDB = append(cfg.ListDB, ankadb.DBConfig{
		Name:   dbname,
		Engine: engine,
		PathDB: dbname,
	})

	ankaDB, err := ankadb.NewAnkaDB(cfg, nil)
	if ankaDB == nil {
		return nil, err
	}

	db := &DB{
		AnkaDB: ankaDB,
	}

	return db, err
}

// UpdCandles - update candles
func (db *DB) UpdCandles(ctx context.Context, candles *tradingpb.Candles) error {
	if candles.Market == "" {
		return ErrInvalidMarket
	}

	if candles.Symbol == "" {
		return ErrInvalidSymbol
	}

	if candles.Tag == "" {
		return ErrInvalidTag
	}

	buf, err := proto.Marshal(candles)
	if err != nil {
		return err
	}

	err = db.AnkaDB.Set(ctx, dbname, makeCandlesDBKey(candles.Market, candles.Symbol, candles.Tag), buf)
	if err != nil {
		return err
	}

	return nil
}

// GetCandles - get candles
func (db *DB) GetCandles(ctx context.Context, market string, symbol string, tags []string, tsStart int64, tsEnd int64) (
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

// GetAllData - get all data
func (db *DB) GetAllData(ctx context.Context) (*TreeMapNode, error) {
	root := NewTreeMapNode("root")

	err := db.AnkaDB.ForEachWithPrefix(ctx, dbname, candlesKeyPrefix, func(key string, buf []byte) error {
		candles := &tradingpb.Candles{}

		err := proto.Unmarshal(buf, candles)
		if err != nil {
			return err
		}

		market := root.GetChildEx(candles.Market)
		symbol := market.GetChildEx(candles.Symbol)
		tag := symbol.GetChildEx(candles.Tag)
		tag.GetChildEx(strconv.Itoa(len(candles.Candles)))

		return nil
	})
	if err != nil {
		return nil, err
	}

	return root, nil
}

// UpdSymbol - update symbol
func (db *DB) UpdSymbol(ctx context.Context, si *tradingpb.SymbolInfo) error {
	if si.Market == "" {
		return ErrInvalidMarket
	}

	if si.Symbol == "" {
		return ErrInvalidSymbol
	}

	buf, err := proto.Marshal(si)
	if err != nil {
		return err
	}

	err = db.AnkaDB.Set(ctx, dbname, makeSymbolDBKey(si.Market, si.Symbol), buf)
	if err != nil {
		return err
	}

	return nil
}

// GetSymbol - get symbol
func (db *DB) GetSymbol(ctx context.Context, market string, symbol string) (
	*tradingpb.SymbolInfo, error) {

	if market == "" {
		return nil, ErrInvalidMarket
	}

	if symbol == "" {
		return nil, ErrInvalidSymbol
	}

	buf, err := db.AnkaDB.Get(ctx, dbname, makeSymbolDBKey(market, symbol))
	if err != nil {
		if err == ankadb.ErrNotFoundKey {
			return nil, nil
		}

		return nil, err
	}

	si := &tradingpb.SymbolInfo{}

	err = proto.Unmarshal(buf, si)
	if err != nil {
		return nil, err
	}

	return si, nil
}

// GetMarketSymbols - get symbols in market
func (db *DB) GetMarketSymbols(ctx context.Context, market string) ([]string, error) {
	if market == "" {
		return nil, ErrInvalidMarket
	}

	symbols := []string{}
	err := db.AnkaDB.ForEachWithPrefix(ctx, dbname, makeSymbolDBKeyPrefix(market), func(key string, buf []byte) error {
		si := &tradingpb.SymbolInfo{}

		err := proto.Unmarshal(buf, si)
		if err != nil {
			return err
		}

		symbols = append(symbols, si.Symbol)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return symbols, nil
}

// GetAssetTimestamp - get tsStart & tsEnd for asset
func (db *DB) GetAssetTimestamp(ctx context.Context, market string, symbol string, tags []string, tsStart int64, tsEnd int64) (
	int64, int64, error) {

	if market == "" {
		return 0, 0, ErrInvalidMarket
	}

	if symbol == "" {
		return 0, 0, ErrInvalidSymbol
	}

	if tsStart > 0 && tsEnd <= 0 {
		tsEnd = time.Now().Unix()
	}

	var mints int64
	var maxts int64

	mints = tsEnd
	maxts = 0

	err := db.AnkaDB.ForEachWithPrefix(ctx, dbname, makeCandlesDBKeyPrefix(market, symbol), func(key string, buf []byte) error {
		cc := &tradingpb.Candles{}

		err := proto.Unmarshal(buf, cc)
		if err != nil {
			return err
		}

		if len(tags) == 0 || tradingdb2utils.IndexOfStringSlice(tags, cc.Tag, 0) >= 0 {
			for _, v := range cc.Candles {
				if v.Ts >= tsStart && v.Ts <= tsEnd {
					if mints > v.Ts {
						mints = v.Ts
					}

					if maxts < v.Ts {
						maxts = v.Ts
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		return 0, 0, err
	}

	return mints, maxts, nil
}

// FixSimTradingParams - reset SimTradingParams
func (db *DB) FixSimTradingParams(ctx context.Context, params *tradingpb.SimTradingParams) (*tradingpb.SimTradingParams, error) {
	var mints int64
	var maxts int64

	mints = time.Now().Unix()
	maxts = 0

	for _, v := range params.Assets {
		cmin, cmax, err := db.GetAssetTimestamp(ctx, v.Market, v.Code, nil, params.StartTs, params.EndTs)
		if err != nil {
			tradingdb2utils.Warn("DB.FixSimTradingParams:GetAssetTimestamp",
				zap.Error(err))

			return nil, err
		}

		if mints > cmin {
			mints = cmin
		}

		if maxts < cmax {
			maxts = cmax
		}
	}

	params.StartTs = mints
	params.EndTs = maxts

	return params, nil
}
