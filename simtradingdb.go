package tradingdb2

import (
	"bytes"
	"context"
	"crypto/sha1"
	"fmt"
	"sync"
	"time"

	ankadb "github.com/zhs007/ankadb"
	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const simtradingDBName = "simtrading"
const simtradingNodesKeyPrefix = "stn:"
const simtradingNodesKey2Prefix = "stn2:"

func makeSimTradingNodesDBKey(strategy string, market string, symbol string, tsStart int64, tsEnd int64) string {
	return fmt.Sprintf("%s%s:%s:%s:%v:%v", simtradingNodesKeyPrefix, strategy, market, symbol, tsStart, tsEnd)
}

func makeSimTradingNodesDBKey2(market string, symbol string, tsStart int64, tsEnd int64, hashHeader string) string {
	return fmt.Sprintf("%s%s:%s:%v:%v:%s", simtradingNodesKey2Prefix, market, symbol, tsStart, tsEnd, hashHeader)
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
	nparam, nbuf, nhash, err := rebuildSimTradingParams(params)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.UpdSimTrading:rebuildSimTradingParams",
			zap.Error(err))

		return err
	}

	cache, err := db.getSimTradingNodes2(ctx, nparam, nhash, nbuf)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.UpdSimTrading:getSimTradingNodes2",
			zap.Error(err))

		return err
	}

	if cache == nil {
		cache = &tradingpb.SimTradingCache{}
	}

	for _, v := range cache.Nodes {
		if bytes.Equal(nbuf, v.Buf) {
			// v.Buf = nil

			v.LastTs = time.Now().Unix()
			err = db.saveSimTradingNodes2Ex(ctx, cache, nparam, nhash)
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

		// v.Buf = nil

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
		// Params: params,
		Key:    nkey,
		LastTs: time.Now().Unix(),
		Buf:    nbuf,
	}

	cache.Nodes = append(cache.Nodes, node)

	err = db.saveSimTradingNodes2Ex(ctx, cache, nparam, nhash)
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
func (db *SimTradingDB) UpdSimTradingEx(ctx context.Context, params *tradingpb.SimTradingParams, pnldata *tradingpb.PNLData,
	dbcache *SimTradingDBCache) (*SimTradingDBCache, error) {

	nparams, nbuf, nhash, err := rebuildSimTradingParams(params)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.UpdSimTradingEx:rebuildSimTradingParams",
			zap.Error(err))

		return nil, err
	}

	isupdcache := false
	if dbcache != nil && dbcache.isMine(params) {
		haspnl, err := dbcache.hasSimTrading(ctx, db, nparams, nhash, nbuf)
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

	cache, err := db.getSimTradingNodes2(ctx, nparams, nhash, nbuf)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.UpdSimTradingEx:getSimTradingNodes",
			zap.Error(err))

		return dbcache, err
	}

	if cache == nil {
		cache = &tradingpb.SimTradingCache{}
	}

	for _, v := range cache.Nodes {
		if bytes.Equal(nbuf, v.Buf) {
			// v.Buf = nil

			v.LastTs = time.Now().Unix()
			err = db.saveSimTradingNodes2Ex(ctx, cache, nparams, nhash)
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

		// v.Buf = nil

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
		// Params: params,
		Key:    nkey,
		LastTs: time.Now().Unix(),
		Buf:    nbuf,
	}

	cache.Nodes = append(cache.Nodes, node)

	err = db.saveSimTradingNodes2Ex(ctx, cache, nparams, nhash)
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
		err = dbcache.addSimTrading(nparams, nhash, nbuf, node)
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

	nparams, nbuf, nhash, err := rebuildSimTradingParams(params)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.GetSimTrading:rebuildSimTradingParams",
			zap.Error(err))

		return nil, err
	}

	cache, err := db.getSimTradingNodes2(ctx, nparams, nhash, nbuf)
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

	// _, cbuf, err := rebuildSimTradingParams(params)
	// if err != nil {
	// 	tradingdb2utils.Warn("SimTradingDB.GetSimTrading:rebuildSimTradingParams params",
	// 		zap.Error(err))

	// 	return nil, err
	// }

	for _, v := range cache.Nodes {
		// if v.Buf == nil {
		// 	_, tbuf, err := rebuildSimTradingParams(v.Params)
		// 	if err != nil {
		// 		tradingdb2utils.Warn("SimTradingDB.GetSimTrading:rebuildSimTradingParams v.Params",
		// 			zap.Error(err))

		// 		return nil, err
		// 	}

		// 	// v.Params = cv
		// 	v.Buf = tbuf
		// }

		if bytes.Equal(nbuf, v.Buf) {
			// v.Buf = nil

			v.LastTs = time.Now().Unix()
			db.saveSimTradingNodes2Ex(ctx, cache, nparams, nhash)

			return db.getPNLData(ctx, v.Key)
		}

		// v.Buf = nil

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

	nparams, nbuf, nhash, err := rebuildSimTradingParams(params)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.GetSimTradingEx:rebuildSimTradingParams",
			zap.Error(err))

		return nil, nil, err
	}

	if dbcache == nil || !dbcache.isMine(params) {
		cache, err := newSimTradingDBCache(ctx, db, nparams, nhash, nbuf)
		if err != nil {
			tradingdb2utils.Warn("SimTradingDB.GetSimTradingEx:newSimTradingDBCache",
				zap.Error(err))

			return nil, nil, err
		}

		dbcache = cache
	}

	if dbcache != nil {
		pnldata, err := dbcache.getSimTrading(ctx, db, nparams, nhash, nbuf)
		if err != nil {
			tradingdb2utils.Warn("SimTradingDB.GetSimTradingEx:getSimTrading",
				zap.Error(err))

			return dbcache, nil, err
		}

		return dbcache, pnldata, nil
	}

	// tradingdb2utils.Debug("SimTradingDB.GetSimTradingEx",
	// 	tradingdb2utils.JSON("params", params))

	cache, err := db.getSimTradingNodes2(ctx, nparams, nhash, nbuf)
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

	// _, cbuf, err := rebuildSimTradingParams(params)
	// if err != nil {
	// 	tradingdb2utils.Warn("SimTradingDB.GetSimTradingEx:rebuildSimTradingParams params",
	// 		zap.Error(err))

	// 	return dbcache, nil, nil
	// }

	for _, v := range cache.Nodes {
		// if v.Buf == nil {
		// 	_, tbuf, err := rebuildSimTradingParams(v.Params)
		// 	if err != nil {
		// 		tradingdb2utils.Warn("SimTradingDB.GetSimTradingEx:rebuildSimTradingParams v.Params",
		// 			zap.Error(err))

		// 		return dbcache, nil, nil
		// 	}

		// 	// v.Params = cv
		// 	// v.Hash = hash
		// 	v.Buf = tbuf
		// }

		if bytes.Equal(nbuf, v.Buf) {
			// v.Buf = nil

			v.LastTs = time.Now().Unix()
			db.saveSimTradingNodes2Ex(ctx, cache, nparams, nhash)

			pnldata, err := db.getPNLData(ctx, v.Key)
			if err != nil {
				tradingdb2utils.Warn("SimTradingDB.GetSimTradingEx:getPNLData",
					zap.Error(err))

				return dbcache, nil, err
			}

			return dbcache, pnldata, err
		}

		// v.Buf = nil

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

// // getSimTradingNodes - get simtrading nodes
// func (db *SimTradingDB) getSimTradingNodes(ctx context.Context, params *tradingpb.SimTradingParams) (
// 	*tradingpb.SimTradingCache, error) {

// 	key := makeSimTradingNodesDBKey(params.Strategies[0].Name, params.Assets[0].Market, params.Assets[0].Code, params.StartTs, params.EndTs)

// 	db.mutexDB.Lock()
// 	buf, err := db.AnkaDB.Get(ctx, simtradingDBName, key)
// 	db.mutexDB.Unlock()
// 	if err != nil {
// 		if err == ankadb.ErrNotFoundKey {
// 			return nil, nil
// 		}

// 		tradingdb2utils.Warn("SimTradingDB.getSimTradingNodes:Get",
// 			zap.Error(err))

// 		return nil, err
// 	}

// 	cache := &tradingpb.SimTradingCache{}

// 	err = proto.Unmarshal(buf, cache)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return cache, nil
// }

// saveSimTradingNodes2 - save simtrading nodes v2
func (db *SimTradingDB) saveSimTradingNodes2(ctx context.Context, cache *tradingpb.SimTradingCache, key2 string) error {
	buf, err := proto.Marshal(cache)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.saveSimTradingNodes2:Marshal",
			zap.Error(err))

		return err
	}

	db.mutexDB.Lock()
	err = db.AnkaDB.Set(ctx, simtradingDBName, key2, buf)
	db.mutexDB.Unlock()
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.saveSimTradingNodes2:Set",
			zap.Error(err))

		return err
	}

	return nil
}

// saveSimTradingNodes2Ex - save simtrading nodes v2
func (db *SimTradingDB) saveSimTradingNodes2Ex(ctx context.Context, cache *tradingpb.SimTradingCache, params *tradingpb.SimTradingParams, hash string) error {
	hashHeader := hash[:2]
	key2 := makeSimTradingNodesDBKey2(params.Assets[0].Market, params.Assets[0].Code, params.StartTs, params.EndTs, hashHeader)

	return db.saveSimTradingNodes2(ctx, cache, key2)
}

// upgradeSimTradingNodes2 - simtrading nodes v1 => simtrading nodes v2
func (db *SimTradingDB) upgradeSimTradingNodes2(ctx context.Context, params *tradingpb.SimTradingParams, hash string, nbuf []byte) (
	*tradingpb.SimTradingCache, error) {

	cache2 := &tradingpb.SimTradingCache{}

	hashHeader := hash[:2]
	key2 := makeSimTradingNodesDBKey2(params.Assets[0].Market, params.Assets[0].Code, params.StartTs, params.EndTs, hashHeader)

	key := makeSimTradingNodesDBKey(params.Strategies[0].Name, params.Assets[0].Market, params.Assets[0].Code, params.StartTs, params.EndTs)

	db.mutexDB.Lock()
	buf, err := db.AnkaDB.Get(ctx, simtradingDBName, key)
	db.mutexDB.Unlock()
	if err != nil {
		if err == ankadb.ErrNotFoundKey {
			err = db.saveSimTradingNodes2(ctx, cache2, key2)
			if err != nil {
				tradingdb2utils.Warn("SimTradingDB.upgradeSimTradingNodes2:saveSimTradingNodes2:nil",
					zap.Error(err))
			}

			return nil, nil
		}

		tradingdb2utils.Warn("SimTradingDB.upgradeSimTradingNodes2:Get",
			zap.Error(err))

		return nil, err
	}

	cache := &tradingpb.SimTradingCache{}

	err = proto.Unmarshal(buf, cache)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.upgradeSimTradingNodes2:Unmarshal",
			zap.Error(err))

		return nil, err
	}

	for _, v := range cache.Nodes {
		_, cb, ch, err := rebuildSimTradingParams(v.Params)
		if err != nil {
			tradingdb2utils.Warn("SimTradingDB.upgradeSimTradingNodes2:rebuildSimTradingParams",
				zap.Error(err))
		}

		if ch[:2] == hashHeader {
			v.Buf = cb
			v.Params = nil

			cache2.Nodes = append(cache2.Nodes, v)
		}
	}

	err = db.saveSimTradingNodes2(ctx, cache2, key2)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.upgradeSimTradingNodes2:saveSimTradingNodes2",
			zap.Error(err))
	}

	return cache2, nil
}

// getSimTradingNodes2 - get simtrading nodes
func (db *SimTradingDB) getSimTradingNodes2(ctx context.Context, params *tradingpb.SimTradingParams, hash string, nbuf []byte) (
	*tradingpb.SimTradingCache, error) {

	key2 := makeSimTradingNodesDBKey2(params.Assets[0].Market, params.Assets[0].Code, params.StartTs, params.EndTs, hash[:2])

	db.mutexDB.Lock()
	buf, err := db.AnkaDB.Get(ctx, simtradingDBName, key2)
	db.mutexDB.Unlock()
	if err != nil {
		if err == ankadb.ErrNotFoundKey {
			return db.upgradeSimTradingNodes2(ctx, params, hash, buf)
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

// // updSimTradingNodes - update simtrading nodes
// func (db *SimTradingDB) updSimTradingNodes(ctx context.Context, params *tradingpb.SimTradingParams, cache *tradingpb.SimTradingCache) error {
// 	key := makeSimTradingNodesDBKey(params.Strategies[0].Name, params.Assets[0].Market, params.Assets[0].Code, params.StartTs, params.EndTs)

// 	buf, err := proto.Marshal(cache)
// 	if err != nil {
// 		return err
// 	}

// 	db.mutexDB.Lock()
// 	err = db.AnkaDB.Set(ctx, simtradingDBName, key, buf)
// 	db.mutexDB.Unlock()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // updSimTradingNodesEx - update simtrading nodes
// func (db *SimTradingDB) updSimTradingNodesEx(ctx context.Context, name string, market string, code string, startTs int64, endTs int64, cache *tradingpb.SimTradingCache) error {
// 	key := makeSimTradingNodesDBKey(name, market, code, startTs, endTs)

// 	buf, err := proto.Marshal(cache)
// 	if err != nil {
// 		return err
// 	}

// 	db.mutexDB.Lock()
// 	err = db.AnkaDB.Set(ctx, simtradingDBName, key, buf)
// 	db.mutexDB.Unlock()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

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

		err = db.AnkaDB.ForEachWithPrefix(ctx, simtradingDBName, simtradingNodesKey2Prefix, func(key string, value []byte) error {
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

// Upgrade2SimTradingDB2 - upgrade to SimTradfingDB2
func (db *SimTradingDB) Upgrade2SimTradingDB2(ctx context.Context, params *tradingpb.SimTradingParams, db2 *SimTradingDB2) (bool, error) {
	nparams, nbuf, nhash, err := rebuildSimTradingParams(params)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.Upgrade2SimTradingDB2:rebuildSimTradingParams",
			zap.Error(err))

		return false, err
	}

	cache, err := db.getSimTradingNodes2(ctx, nparams, nhash, nbuf)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.Upgrade2SimTradingDB2:getSimTradingNodes",
			zap.Error(err))

		return false, err
	}

	if cache == nil {
		tradingdb2utils.Debug("SimTradingDB.Upgrade2SimTradingDB2:no cache",
			tradingdb2utils.JSON("params", params))

		return false, nil
	}

	for _, v := range cache.Nodes {
		pnldata, err := db.getPNLData(ctx, v.Key)
		if err != nil {
			tradingdb2utils.Warn("SimTradingDB.Upgrade2SimTradingDB2:getPNLData",
				zap.Error(err))

			return false, err
		}

		err = db2.UpdSimTrading(ctx, v.Params, pnldata)
		if err != nil {
			tradingdb2utils.Warn("SimTradingDB.Upgrade2SimTradingDB2:UpdSimTrading",
				zap.Error(err))

			return false, err
		}
	}

	key2 := makeSimTradingNodesDBKey2(params.Assets[0].Market, params.Assets[0].Code, params.StartTs, params.EndTs, nhash[:2])

	db.mutexDB.Lock()
	err = db.AnkaDB.Delete(ctx, simtradingDBName, key2)
	db.mutexDB.Unlock()

	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.Upgrade2SimTradingDB2:Delete",
			zap.Error(err))

		return false, err
	}

	return true, nil
}
