package tradingdb2

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	ankadb "github.com/zhs007/ankadb"
	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
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
func (db *SimTradingDB) UpdSimTrading(ctx context.Context, params *tradingpb.SimTradingParams, pnldata *tradingpb.PNLData) error {
	pnldata.Lastts = time.Now().Unix()

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
func (db *SimTradingDB) GetSimTrading(ctx context.Context, params *tradingpb.SimTradingParams) (
	*tradingpb.PNLData, error) {

	key := makeSimTradingDBKey(params.Strategies[0].Name, params.Assets[0].Market, params.Assets[0].Code, params.StartTs, params.EndTs)

	buf, err := db.AnkaDB.Get(ctx, dbname, key)
	if err != nil {
		if err == ankadb.ErrNotFoundKey {
			return nil, nil
		}

		tradingdb2utils.Warn("SimTradingDB.GetSimTrading:Get",
			zap.Error(err))

		return nil, err
	}

	pnl := &tradingpb.PNLData{}

	err = proto.Unmarshal(buf, pnl)
	if err != nil {
		return nil, err
	}

	err = db.UpdSimTrading(ctx, params, pnl)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.GetSimTrading:UpdSimTrading",
			zap.Error(err))

		return nil, err
	}

	return pnl, nil
}
