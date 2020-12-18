package tradingdb2

import (
	"context"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	ankadb "github.com/zhs007/ankadb"
	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
)

const candlesDB2KeyPrefix = "c:"
const symbolDB2KeyPrefix = "s:"

func makeCandlesDB2Key(market string, symbol string, tag string) string {
	return tradingdb2utils.AppendString(candlesDB2KeyPrefix, symbol, ":", tag)
}

func makeCandlesDB2KeyPrefix(market string, symbol string) string {
	return tradingdb2utils.AppendString(candlesDB2KeyPrefix, symbol, ":")
}

func makeSymbolDB2Key(market string, symbol string) string {
	return tradingdb2utils.AppendString(symbolDB2KeyPrefix, symbol)
}

func makeSymbolDB2KeyPrefix(market string) string {
	return tradingdb2utils.AppendString(symbolDB2KeyPrefix)
}

// DB2 - database v2
type DB2 struct {
	AnkaDB ankadb.AnkaDB
	cfg    *Config
}

// NewDB2 - new DB2
func NewDB2(cfg *Config) (*DB2, error) {
	dbcfg := ankadb.NewConfig()

	dbcfg.PathDBRoot = cfg.DBPath

	for _, v := range cfg.DB2Markets {
		dbcfg.ListDB = append(dbcfg.ListDB, ankadb.DBConfig{
			Name:   v,
			Engine: cfg.DBEngine,
			PathDB: v,
		})
	}

	ankaDB, err := ankadb.NewAnkaDB(dbcfg, nil)
	if ankaDB == nil {
		return nil, err
	}

	db2 := &DB2{
		AnkaDB: ankaDB,
		cfg:    cfg,
	}

	return db2, err
}

// UpdCandles - update candles
func (db2 *DB2) UpdCandles(ctx context.Context, candles *tradingpb.Candles) error {
	if candles.Market == "" || tradingdb2utils.IndexOfStringSlice(db2.cfg.DB2Markets, candles.Market, 0) < 0 {
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

	err = db2.AnkaDB.Set(ctx, candles.Market, makeCandlesDB2Key(candles.Market, candles.Symbol, candles.Tag), buf)
	if err != nil {
		return err
	}

	return nil
}

// GetCandles - get candles
func (db2 *DB2) GetCandles(ctx context.Context, market string, symbol string, tsStart int64, tsEnd int64) (*tradingpb.Candles, error) {
	if market == "" || tradingdb2utils.IndexOfStringSlice(db2.cfg.DB2Markets, market, 0) < 0 {
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

	err := db2.AnkaDB.ForEachWithPrefix(ctx, market, makeCandlesDB2KeyPrefix(market, symbol), func(key string, buf []byte) error {
		cc := &tradingpb.Candles{}

		err := proto.Unmarshal(buf, cc)
		if err != nil {
			return err
		}

		if tsStart > 0 || tsEnd > 0 {
			for _, v := range cc.Candles {
				if v.Ts >= tsStart && v.Ts <= tsEnd {
					candles.Candles = append(candles.Candles, v)
				}
			}
		} else {
			candles.Candles = append(candles.Candles, cc.Candles...)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return candles, nil
}

// HasCandles - has candles
func (db2 *DB2) HasCandles(ctx context.Context, market string, symbol string) bool {
	if market == "" || tradingdb2utils.IndexOfStringSlice(db2.cfg.DB2Markets, market, 0) < 0 {
		return false
	}

	if symbol == "" {
		return false
	}

	hascandles := false
	err := db2.AnkaDB.ForEachWithPrefix(ctx, market, makeCandlesDB2KeyPrefix(market, symbol), func(key string, buf []byte) error {
		hascandles = true

		return nil
	})
	if err != nil {
		return false
	}

	return hascandles
}

// GetAllData - get all data
func (db2 *DB2) GetAllData(ctx context.Context) (*TreeMapNode, error) {
	root := NewTreeMapNode("root")

	for _, v := range db2.cfg.DB2Markets {
		err := db2.AnkaDB.ForEachWithPrefix(ctx, v, candlesDB2KeyPrefix, func(key string, buf []byte) error {
			candles := &tradingpb.Candles{}

			err := proto.Unmarshal(buf, candles)
			if err != nil {
				return err
			}

			market := root.GetChildEx(v)
			symbol := market.GetChildEx(candles.Symbol)
			tag := symbol.GetChildEx(candles.Tag)
			tag.GetChildEx(strconv.Itoa(len(candles.Candles)))

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return root, nil
}

// UpdSymbol - update symbol
func (db2 *DB2) UpdSymbol(ctx context.Context, si *tradingpb.SymbolInfo) error {
	if si.Market == "" || tradingdb2utils.IndexOfStringSlice(db2.cfg.DB2Markets, si.Market, 0) < 0 {
		return ErrInvalidMarket
	}

	if si.Symbol == "" {
		return ErrInvalidSymbol
	}

	buf, err := proto.Marshal(si)
	if err != nil {
		return err
	}

	err = db2.AnkaDB.Set(ctx, si.Market, makeSymbolDB2Key(si.Market, si.Symbol), buf)
	if err != nil {
		return err
	}

	return nil
}

// GetSymbol - get symbol
func (db2 *DB2) GetSymbol(ctx context.Context, market string, symbol string) (*tradingpb.SymbolInfo, error) {

	if market == "" || tradingdb2utils.IndexOfStringSlice(db2.cfg.DB2Markets, market, 0) < 0 {
		return nil, ErrInvalidMarket
	}

	if symbol == "" {
		return nil, ErrInvalidSymbol
	}

	buf, err := db2.AnkaDB.Get(ctx, market, makeSymbolDB2Key(market, symbol))
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
func (db2 *DB2) GetMarketSymbols(ctx context.Context, market string) ([]string, error) {
	if market == "" || tradingdb2utils.IndexOfStringSlice(db2.cfg.DB2Markets, market, 0) < 0 {
		return nil, ErrInvalidMarket
	}

	symbols := []string{}
	err := db2.AnkaDB.ForEachWithPrefix(ctx, market, makeSymbolDB2KeyPrefix(market), func(key string, buf []byte) error {
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
func (db2 *DB2) GetAssetTimestamp(ctx context.Context, market string, symbol string, tsStart int64, tsEnd int64) (
	int64, int64, error) {

	if market == "" || tradingdb2utils.IndexOfStringSlice(db2.cfg.DB2Markets, market, 0) < 0 {
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

	err := db2.AnkaDB.ForEachWithPrefix(ctx, market, makeCandlesDB2KeyPrefix(market, symbol), func(key string, buf []byte) error {
		cc := &tradingpb.Candles{}

		err := proto.Unmarshal(buf, cc)
		if err != nil {
			return err
		}

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

		return nil
	})
	if err != nil {
		return 0, 0, err
	}

	return mints, maxts, nil
}

// FixSimTradingParams - reset SimTradingParams
func (db2 *DB2) FixSimTradingParams(ctx context.Context, params *tradingpb.SimTradingParams) (*tradingpb.SimTradingParams, error) {
	var mints int64
	var maxts int64

	mints = time.Now().Unix()
	maxts = 0

	for _, v := range params.Assets {
		cmin, cmax, err := db2.GetAssetTimestamp(ctx, v.Market, v.Code, params.StartTs, params.EndTs)
		if err != nil {
			tradingdb2utils.Warn("DB2.FixSimTradingParams:GetAssetTimestamp",
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