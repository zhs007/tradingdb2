package tradingdb2

import (
	"context"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	ankadb "github.com/zhs007/ankadb"
	tradingdb2pb "github.com/zhs007/tradingdb2/tradingdb2pb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
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
func (db *DB) UpdCandles(ctx context.Context, candles *tradingdb2pb.Candles) error {
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
	*tradingdb2pb.Candles, error) {

	if market == "" {
		return nil, ErrInvalidMarket
	}

	if symbol == "" {
		return nil, ErrInvalidSymbol
	}

	candles := &tradingdb2pb.Candles{
		Market: market,
		Symbol: symbol,
	}

	if tsStart > 0 && tsEnd <= 0 {
		tsEnd = time.Now().Unix()
	}

	err := db.AnkaDB.ForEachWithPrefix(ctx, dbname, makeCandlesDBKeyPrefix(market, symbol), func(key string, buf []byte) error {
		cc := &tradingdb2pb.Candles{}

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
		candles := &tradingdb2pb.Candles{}

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
func (db *DB) UpdSymbol(ctx context.Context, si *tradingdb2pb.SymbolInfo) error {
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
	*tradingdb2pb.SymbolInfo, error) {

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

	si := &tradingdb2pb.SymbolInfo{}

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
		si := &tradingdb2pb.SymbolInfo{}

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
