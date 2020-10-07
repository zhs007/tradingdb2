package tradingdb2

import (
	"context"
	"strconv"

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

func makeSymbolDBKey(market string, symbol string) string {
	return tradingdb2utils.AppendString(symbolKeyPrefix, market, ":", symbol)
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
func (db *DB) GetCandles(ctx context.Context, market string, symbol string, tag string) (
	*tradingdb2pb.Candles, error) {

	if market == "" {
		return nil, ErrInvalidMarket
	}

	if symbol == "" {
		return nil, ErrInvalidSymbol
	}

	if tag == "" {
		return nil, ErrInvalidTag
	}

	buf, err := db.AnkaDB.Get(ctx, dbname, makeCandlesDBKey(market, symbol, tag))
	if err != nil {
		if err == ankadb.ErrNotFoundKey {
			return nil, nil
		}

		return nil, err
	}

	res := &tradingdb2pb.Candles{}

	err = proto.Unmarshal(buf, res)
	if err != nil {
		return nil, err
	}

	return res, nil
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
