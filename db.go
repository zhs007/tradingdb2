package tradingdb2

import (
	"context"

	"github.com/golang/protobuf/proto"
	ankadb "github.com/zhs007/ankadb"
	tradingdb2pb "github.com/zhs007/tradingdb2/pb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
)

const dbname = "tradingdb2"

func makeCandlesDBKey(market string, symbol string, tag string) string {
	return tradingdb2utils.AppendString(market, ":", symbol, ":", tag)
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
