package tradingdb2

import (
	"context"
	"sync"
	"time"

	ankadb "github.com/zhs007/ankadb"
	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const simtradingDB2Name = "simtrading2"

// SimTradingDB2 - database
type SimTradingDB2 struct {
	AnkaDB      ankadb.AnkaDB
	mutexDB     sync.Mutex
	lastClearTs int64
	ticker      *time.Ticker
	chanDone    chan int
}

// NewSimTradingDB2 - new SimTradingDB2
func NewSimTradingDB2(dbpath string, httpAddr string, engine string) (*SimTradingDB2, error) {
	cfg := ankadb.NewConfig()

	cfg.AddrHTTP = httpAddr
	cfg.PathDBRoot = dbpath
	cfg.ListDB = append(cfg.ListDB, ankadb.DBConfig{
		Name:   simtradingDB2Name,
		Engine: engine,
		PathDB: simtradingDB2Name,
	})

	ankaDB, err := ankadb.NewAnkaDB(cfg, nil)
	if ankaDB == nil {
		return nil, err
	}

	db := &SimTradingDB2{
		AnkaDB:   ankaDB,
		ticker:   time.NewTicker(time.Hour),
		chanDone: make(chan int),
	}

	go db.onTimer()

	return db, err
}

// UpdSimTrading - update simulation trading
func (db *SimTradingDB2) UpdSimTrading(ctx context.Context, params *tradingpb.SimTradingParams, pnldata *tradingpb.PNLData) error {
	_, nbuf, err := rebuildSimTradingParams2(params)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB2.UpdSimTrading:rebuildSimTradingParams2",
			zap.Error(err))

		return err
	}

	pnldata.CacheTs = time.Now().Unix()

	buf, err := proto.Marshal(pnldata)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB2.UpdSimTrading:Marshal",
			zap.Error(err))

		return err
	}

	db.mutexDB.Lock()
	err = db.AnkaDB.SetEx(ctx, simtradingDB2Name, nbuf, buf)
	db.mutexDB.Unlock()
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB2.UpdSimTrading:Set",
			zap.Error(err))

		return err
	}

	return nil
}

// GetSimTrading - get candles
func (db *SimTradingDB2) GetSimTrading(ctx context.Context, params *tradingpb.SimTradingParams) (
	*tradingpb.PNLData, error) {
	_, nbuf, err := rebuildSimTradingParams2(params)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB2.GetSimTrading:rebuildSimTradingParams2",
			zap.Error(err))

		return nil, err
	}

	db.mutexDB.Lock()
	buf, err := db.AnkaDB.GetEx(ctx, simtradingDB2Name, nbuf)
	db.mutexDB.Unlock()
	if err != nil {
		if err == ankadb.ErrNotFoundKey {
			return nil, nil
		}

		tradingdb2utils.Warn("SimTradingDB2.getSimTradingNodes:Get",
			zap.Error(err))

		return nil, err
	}

	pnl := &tradingpb.PNLData{}

	err = proto.Unmarshal(buf, pnl)
	if err != nil {
		return nil, err
	}

	return pnl, nil
}

// ClearCache - clear cache
func (db *SimTradingDB2) ClearCache(ctx context.Context) error {

	if time.Now().Unix()-db.lastClearTs > SimTradingCacheTimerTs {
		curts := time.Now().Unix()

		db.mutexDB.Lock()
		defer db.mutexDB.Unlock()

		err := db.AnkaDB.ForEachWithPrefix(ctx, simtradingDB2Name, simtradingNodesKeyPrefix, func(key string, value []byte) error {
			pnl := &tradingpb.PNLData{}

			err := proto.Unmarshal(value, pnl)
			if err != nil {
				tradingdb2utils.Warn("SimTradingDB2.ClearCache:FuncAnkaDBForEach:Unmarshal",
					zap.Error(err),
					zap.String("key", key))

				return err
			}

			if curts-pnl.CacheTs > SimTradingCacheTimeOut {
				err = db.AnkaDB.Delete(ctx, simtradingDB2Name, key)
				if err != nil {
					tradingdb2utils.Warn("SimTradingDB2.ClearCache:FuncAnkaDBForEach:Delete",
						zap.Error(err),
						zap.String("key", key))
				}
			}

			return nil
		})

		if err != nil {
			tradingdb2utils.Warn("SimTradingDB2.ClearCache:FuncAnkaDBForEach",
				zap.Error(err))

			return err
		}
	}

	return nil
}

// onTimer - on timer
func (db *SimTradingDB2) onTimer() {
	for {
		select {
		case <-db.chanDone:
			return
		case <-db.ticker.C:
			db.ClearCache(context.Background())
		}
	}
}

// Stop - stop
func (db *SimTradingDB2) Stop() {
	db.ticker.Stop()
	db.chanDone <- 0
}
