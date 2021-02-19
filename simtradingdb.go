package tradingdb2

import (
	"context"
	"crypto/sha1"
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	ankadb "github.com/zhs007/ankadb"
	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
)

const simtradingDBName = "simtrading"
const simtradingNodesKeyPrefix = "stn:"

func makeSimTradingNodesDBKey(strategy string, market string, symbol string, tsStart int64, tsEnd int64) string {
	return fmt.Sprintf("%s%s:%s:%s:%v:%v", simtradingNodesKeyPrefix, strategy, market, symbol, tsStart, tsEnd)
}

// SimTradingDB - database
type SimTradingDB struct {
	AnkaDB      ankadb.AnkaDB
	mutexDB     sync.Mutex
	lastClearTs int64
	ticker      *time.Ticker
	chanDone    chan int
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
		AnkaDB:   ankaDB,
		ticker:   time.NewTicker(time.Hour),
		chanDone: make(chan int),
	}

	go db.onTimer()

	return db, err
}

// UpdSimTrading - update simulation trading
func (db *SimTradingDB) UpdSimTrading(ctx context.Context, params *tradingpb.SimTradingParams, pnldata *tradingpb.PNLData) error {
	cache, err := db.getSimTradingNodes(ctx, params)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.UpdSimTrading:getSimTradingNodes",
			zap.Error(err))

		return err
	}

	if cache == nil {
		cache = &tradingpb.SimTradingCache{}
	}

	_, chash, err := rebuildSimTradingParams(params)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.UpdSimTrading:rebuildSimTradingParams params",
			zap.Error(err))

		return err
	}

	for _, v := range cache.Nodes {
		if v.Hash == "" {
			cv, hash, err := rebuildSimTradingParams(v.Params)
			if err != nil {
				tradingdb2utils.Warn("SimTradingDB.UpdSimTrading:rebuildSimTradingParams v.Params",
					zap.Error(err))

				return err
			}

			v.Params = cv
			v.Hash = hash
		}

		if chash == v.Hash {
			v.LastTs = time.Now().Unix()
			err = db.updSimTradingNodes(ctx, params, cache)
			if err != nil {
				tradingdb2utils.Warn("SimTradingDB.UpdSimTrading:isSameHash updSimTradingNodes",
					zap.Error(err))

				return err
			}

			err = db.updPNLData(ctx, v.Key, pnldata)
			if err != nil {
				tradingdb2utils.Warn("SimTradingDB.UpdSimTrading:updPNLData",
					zap.Error(err))

				return err
			}

			return nil
		}

		// if db.isSameSimTradingParams(params, v.Params) {
		// 	v.LastTs = time.Now().Unix()
		// 	err = db.updSimTradingNodes(ctx, params, cache)
		// 	if err != nil {
		// 		tradingdb2utils.Warn("SimTradingDB.UpdSimTrading:isSameSimTradingParams updSimTradingNodes",
		// 			zap.Error(err))

		// 		return err
		// 	}

		// 	err = db.updPNLData(ctx, v.Key, pnldata)
		// 	if err != nil {
		// 		tradingdb2utils.Warn("SimTradingDB.UpdSimTrading:updPNLData",
		// 			zap.Error(err))

		// 		return err
		// 	}

		// 	return nil
		// }
	}

	nkey, err := db.hashParams(pnldata)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.UpdSimTrading:hashParams",
			zap.Error(err))

		return err
	}

	node := &tradingpb.SimTradingCacheNode{
		Params: params,
		Key:    nkey,
		LastTs: time.Now().Unix(),
	}

	cache.Nodes = append(cache.Nodes, node)

	err = db.updSimTradingNodes(ctx, params, cache)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.UpdSimTrading:updSimTradingNodes",
			zap.Error(err))

		return err
	}

	err = db.updPNLData(ctx, nkey, pnldata)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.UpdSimTrading:updPNLData",
			zap.Error(err))

		return err
	}

	return nil
}

// UpdSimTradingEx - update simulation trading
func (db *SimTradingDB) UpdSimTradingEx(ctx context.Context, params *tradingpb.SimTradingParams, pnldata *tradingpb.PNLData, dbcache *SimTradingDBCache) (*SimTradingDBCache, error) {
	isupdcache := false
	if dbcache != nil && dbcache.isMine(params) {
		haspnl, err := dbcache.hasSimTrading(ctx, db, params)
		if err != nil {
			tradingdb2utils.Warn("SimTradingDB.UpdSimTradingEx:hasSimTrading",
				zap.Error(err))
		}

		if haspnl {
			dbcache = nil
		} else {
			isupdcache = true
		}
	}

	cache, err := db.getSimTradingNodes(ctx, params)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.UpdSimTradingEx:getSimTradingNodes",
			zap.Error(err))

		return dbcache, err
	}

	if cache == nil {
		cache = &tradingpb.SimTradingCache{}
	}

	_, chash, err := rebuildSimTradingParams(params)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.UpdSimTradingEx:rebuildSimTradingParams params",
			zap.Error(err))

		return dbcache, err
	}

	for _, v := range cache.Nodes {
		if v.Hash == "" {
			cv, hash, err := rebuildSimTradingParams(v.Params)
			if err != nil {
				tradingdb2utils.Warn("SimTradingDB.UpdSimTradingEx:rebuildSimTradingParams v.Params",
					zap.Error(err))

				return dbcache, err
			}

			v.Params = cv
			v.Hash = hash
		}

		if chash == v.Hash {
			v.LastTs = time.Now().Unix()
			err = db.updSimTradingNodes(ctx, params, cache)
			if err != nil {
				tradingdb2utils.Warn("SimTradingDB.UpdSimTradingEx:isSameSimTradingParams updSimTradingNodes",
					zap.Error(err))

				return dbcache, err
			}

			err = db.updPNLData(ctx, v.Key, pnldata)
			if err != nil {
				tradingdb2utils.Warn("SimTradingDB.UpdSimTradingEx:updPNLData",
					zap.Error(err))

				return dbcache, err
			}

			return dbcache, nil
		}

		// if db.isSameSimTradingParams(params, v.Params) {
		// 	v.LastTs = time.Now().Unix()
		// 	err = db.updSimTradingNodes(ctx, params, cache)
		// 	if err != nil {
		// 		tradingdb2utils.Warn("SimTradingDB.UpdSimTradingEx:isSameSimTradingParams updSimTradingNodes",
		// 			zap.Error(err))

		// 		return dbcache, err
		// 	}

		// 	err = db.updPNLData(ctx, v.Key, pnldata)
		// 	if err != nil {
		// 		tradingdb2utils.Warn("SimTradingDB.UpdSimTradingEx:updPNLData",
		// 			zap.Error(err))

		// 		return dbcache, err
		// 	}

		// 	return dbcache, nil
		// }
	}

	nkey, err := db.hashParams(pnldata)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.UpdSimTradingEx:hashParams",
			zap.Error(err))

		return dbcache, err
	}

	node := &tradingpb.SimTradingCacheNode{
		Params: params,
		Key:    nkey,
		LastTs: time.Now().Unix(),
	}

	cache.Nodes = append(cache.Nodes, node)

	err = db.updSimTradingNodes(ctx, params, cache)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.UpdSimTradingEx:updSimTradingNodes",
			zap.Error(err))

		return dbcache, err
	}

	err = db.updPNLData(ctx, nkey, pnldata)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.UpdSimTradingEx:updPNLData",
			zap.Error(err))

		return dbcache, err
	}

	if isupdcache {
		err = dbcache.addSimTrading(params, node)
		if err != nil {
			tradingdb2utils.Warn("SimTradingDB.UpdSimTradingEx:addSimTrading",
				zap.Error(err))
		}
	}

	return dbcache, nil
}

// GetSimTrading - get candles
func (db *SimTradingDB) GetSimTrading(ctx context.Context, params *tradingpb.SimTradingParams) (
	*tradingpb.PNLData, error) {
	// tradingdb2utils.Debug("SimTradingDB.GetSimTrading",
	// 	tradingdb2utils.JSON("params", params))

	cache, err := db.getSimTradingNodes(ctx, params)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.GetSimTrading:getSimTradingNodes",
			zap.Error(err))

		return nil, err
	}

	if cache == nil {
		tradingdb2utils.Debug("SimTradingDB.GetSimTrading:no cache",
			tradingdb2utils.JSON("params", params))

		return nil, nil
	}

	_, chash, err := rebuildSimTradingParams(params)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.GetSimTrading:rebuildSimTradingParams params",
			zap.Error(err))

		return nil, err
	}

	for _, v := range cache.Nodes {
		if v.Hash == "" {
			cv, hash, err := rebuildSimTradingParams(v.Params)
			if err != nil {
				tradingdb2utils.Warn("SimTradingDB.GetSimTrading:rebuildSimTradingParams v.Params",
					zap.Error(err))

				return nil, err
			}

			v.Params = cv
			v.Hash = hash
		}

		if chash == v.Hash {
			v.LastTs = time.Now().Unix()
			db.updSimTradingNodes(ctx, params, cache)

			return db.getPNLData(ctx, v.Key)
		}

		// if db.isSameSimTradingParams(params, v.Params) {
		// 	v.LastTs = time.Now().Unix()
		// 	db.updSimTradingNodes(ctx, params, cache)

		// 	return db.getPNLData(ctx, v.Key)
		// }

		// tradingdb2utils.Debug("SimTradingDB.GetSimTrading:isSameSimTradingParams",
		// 	tradingdb2utils.JSON("v", v))
	}

	return nil, nil
}

// GetSimTradingEx - get candles
func (db *SimTradingDB) GetSimTradingEx(ctx context.Context, params *tradingpb.SimTradingParams, dbcache *SimTradingDBCache) (
	*SimTradingDBCache, *tradingpb.PNLData, error) {
	if dbcache == nil || !dbcache.isMine(params) {
		cache, err := newSimTradingDBCache(ctx, db, params)
		if err != nil {
			tradingdb2utils.Warn("SimTradingDB.GetSimTradingEx:newSimTradingDBCache",
				zap.Error(err))

			return nil, nil, err
		}

		dbcache = cache
	}

	if dbcache != nil {
		pnldata, err := dbcache.getSimTrading(ctx, db, params)
		if err != nil {
			tradingdb2utils.Warn("SimTradingDB.GetSimTradingEx:getSimTrading",
				zap.Error(err))

			return dbcache, nil, err
		}

		return dbcache, pnldata, nil
	}

	// tradingdb2utils.Debug("SimTradingDB.GetSimTradingEx",
	// 	tradingdb2utils.JSON("params", params))

	cache, err := db.getSimTradingNodes(ctx, params)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.GetSimTradingEx:getSimTradingNodes",
			zap.Error(err))

		return dbcache, nil, err
	}

	if cache == nil {
		tradingdb2utils.Debug("SimTradingDB.GetSimTradingEx:no cache",
			tradingdb2utils.JSON("params", params))

		return dbcache, nil, nil
	}

	_, chash, err := rebuildSimTradingParams(params)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.GetSimTradingEx:rebuildSimTradingParams params",
			zap.Error(err))

		return dbcache, nil, nil
	}

	for _, v := range cache.Nodes {
		if v.Hash == "" {
			cv, hash, err := rebuildSimTradingParams(v.Params)
			if err != nil {
				tradingdb2utils.Warn("SimTradingDB.GetSimTradingEx:rebuildSimTradingParams v.Params",
					zap.Error(err))

				return dbcache, nil, nil
			}

			v.Params = cv
			v.Hash = hash
		}

		if chash == v.Hash {
			v.LastTs = time.Now().Unix()
			db.updSimTradingNodes(ctx, params, cache)

			pnldata, err := db.getPNLData(ctx, v.Key)
			if err != nil {
				tradingdb2utils.Warn("SimTradingDB.GetSimTradingEx:getPNLData",
					zap.Error(err))

				return dbcache, nil, err
			}

			return dbcache, pnldata, err
		}

		// if db.isSameSimTradingParams(params, v.Params) {
		// 	v.LastTs = time.Now().Unix()
		// 	db.updSimTradingNodes(ctx, params, cache)

		// 	pnldata, err := db.getPNLData(ctx, v.Key)
		// 	if err != nil {
		// 		tradingdb2utils.Warn("SimTradingDB.GetSimTradingEx:getPNLData",
		// 			zap.Error(err))

		// 		return dbcache, nil, err
		// 	}

		// 	return dbcache, pnldata, err
		// }

		// tradingdb2utils.Debug("SimTradingDB.GetSimTradingEx:isSameSimTradingParams",
		// 	tradingdb2utils.JSON("v", v))
	}

	return dbcache, nil, nil
}

// getSimTradingNodes - get simtrading nodes
func (db *SimTradingDB) getSimTradingNodes(ctx context.Context, params *tradingpb.SimTradingParams) (
	*tradingpb.SimTradingCache, error) {

	key := makeSimTradingNodesDBKey(params.Strategies[0].Name, params.Assets[0].Market, params.Assets[0].Code, params.StartTs, params.EndTs)

	db.mutexDB.Lock()
	buf, err := db.AnkaDB.Get(ctx, simtradingDBName, key)
	db.mutexDB.Unlock()
	if err != nil {
		if err == ankadb.ErrNotFoundKey {
			return nil, nil
		}

		tradingdb2utils.Warn("SimTradingDB.getSimTradingNodes:Get",
			zap.Error(err))

		return nil, err
	}

	cache := &tradingpb.SimTradingCache{}

	err = proto.Unmarshal(buf, cache)
	if err != nil {
		return nil, err
	}

	return cache, nil
}

// updSimTradingNodes - update simtrading nodes
func (db *SimTradingDB) updSimTradingNodes(ctx context.Context, params *tradingpb.SimTradingParams, cache *tradingpb.SimTradingCache) error {
	key := makeSimTradingNodesDBKey(params.Strategies[0].Name, params.Assets[0].Market, params.Assets[0].Code, params.StartTs, params.EndTs)

	buf, err := proto.Marshal(cache)
	if err != nil {
		return err
	}

	db.mutexDB.Lock()
	err = db.AnkaDB.Set(ctx, simtradingDBName, key, buf)
	db.mutexDB.Unlock()
	if err != nil {
		return err
	}

	return nil
}

// updSimTradingNodesEx - update simtrading nodes
func (db *SimTradingDB) updSimTradingNodesEx(ctx context.Context, name string, market string, code string, startTs int64, endTs int64, cache *tradingpb.SimTradingCache) error {
	key := makeSimTradingNodesDBKey(name, market, code, startTs, endTs)

	buf, err := proto.Marshal(cache)
	if err != nil {
		return err
	}

	db.mutexDB.Lock()
	err = db.AnkaDB.Set(ctx, simtradingDBName, key, buf)
	db.mutexDB.Unlock()
	if err != nil {
		return err
	}

	return nil
}

// // isSameSimTradingParams - is same SimTradingParams
// func (db *SimTradingDB) isSameSimTradingParams(v0 *tradingpb.SimTradingParams, v1 *tradingpb.SimTradingParams) bool {
// 	if len(v0.Strategies) == len(v1.Strategies) {
// 		for si, sv0 := range v0.Strategies {
// 			sv1 := v1.Strategies[si]
// 			if !IsSameStrategy(sv0, sv1) {
// 				return false
// 			}
// 		}
// 	}

// 	return true
// }

// getPNLData - get candles
func (db *SimTradingDB) getPNLData(ctx context.Context, key string) (
	*tradingpb.PNLData, error) {

	db.mutexDB.Lock()
	buf, err := db.AnkaDB.Get(ctx, simtradingDBName, key)
	db.mutexDB.Unlock()
	if err != nil {
		if err == ankadb.ErrNotFoundKey {
			return nil, nil
		}

		tradingdb2utils.Warn("SimTradingDB.getPNLData:Get",
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

// hashParams - hash parameters
func (db *SimTradingDB) hashParams(pnldata *tradingpb.PNLData) (string, error) {
	buf, err := proto.Marshal(pnldata)
	if err != nil {
		return "", err
	}

	h := sha1.New()
	h.Write(buf)
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs), nil
}

// updPNLData - update PNLData
func (db *SimTradingDB) updPNLData(ctx context.Context, key string, pnldata *tradingpb.PNLData) error {
	buf, err := proto.Marshal(pnldata)
	if err != nil {
		return err
	}

	db.mutexDB.Lock()
	err = db.AnkaDB.Set(ctx, simtradingDBName, key, buf)
	db.mutexDB.Unlock()
	if err != nil {
		return err
	}

	return nil
}

// ClearCache - clear cache
func (db *SimTradingDB) ClearCache(ctx context.Context) error {

	if time.Now().Unix()-db.lastClearTs > SimTradingCacheTimerTs {
		curts := time.Now().Unix()

		db.mutexDB.Lock()
		defer db.mutexDB.Unlock()

		err := db.AnkaDB.ForEachWithPrefix(ctx, simtradingDBName, simtradingNodesKeyPrefix, func(key string, value []byte) error {
			cache := &tradingpb.SimTradingCache{}

			err := proto.Unmarshal(value, cache)
			if err != nil {
				tradingdb2utils.Warn("SimTradingDB.ClearCache:FuncAnkaDBForEach:Unmarshal",
					zap.Error(err),
					zap.String("key", key))

				return err
			}

			delnums := 0
			newcache := &tradingpb.SimTradingCache{}
			for _, v := range cache.Nodes {
				if curts-v.LastTs > SimTradingCacheTimeOut {
					delnums++

					err = db.AnkaDB.Delete(ctx, simtradingDBName, v.Key)

					tradingdb2utils.Warn("SimTradingDB.ClearCache:FuncAnkaDBForEach:Delete",
						zap.Error(err),
						zap.String("key", v.Key))
				} else {
					newcache.Nodes = append(newcache.Nodes, v)
				}
			}

			if delnums > 0 {
				buf, err := proto.Marshal(newcache)
				if err != nil {
					tradingdb2utils.Warn("SimTradingDB.ClearCache:FuncAnkaDBForEach:Marshal",
						zap.Error(err),
						zap.String("key", key))

					return err
				}

				err = db.AnkaDB.Set(ctx, simtradingDBName, key, buf)
				if err != nil {
					return err
				}

				tradingdb2utils.Info("SimTradingDB.ClearCache",
					zap.String("key", key),
					zap.Int("nums", delnums),
					zap.Int("firstnums", len(cache.Nodes)),
					zap.Int("lastnums", len(newcache.Nodes)))
			}

			return nil
		})

		if err != nil {
			tradingdb2utils.Warn("SimTradingDB.ClearCache:FuncAnkaDBForEach",
				zap.Error(err))

			return err
		}
	}

	return nil
}

// onTimer - on timer
func (db *SimTradingDB) onTimer() {
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
func (db *SimTradingDB) Stop() {
	db.ticker.Stop()
	db.chanDone <- 0
}
