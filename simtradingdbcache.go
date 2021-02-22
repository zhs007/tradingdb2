package tradingdb2

import (
	"context"
	"fmt"
	"sync"
	"time"

	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
)

// // SimTradingDBCacheNode -
// type SimTradingDBCacheNode struct {
// 	node *tradingpb.SimTradingCacheNode
// }

// SimTradingDBCache -
type SimTradingDBCache struct {
	mapCache      map[string]*tradingpb.SimTradingCacheNode
	AssetMarket   string
	AssetCode     string
	StartTs       int64
	EndTs         int64
	mapHashHeader map[string]*tradingpb.SimTradingCache
	mutexCache    sync.Mutex
}

// SimTradingDBCacheObj -
type SimTradingDBCacheObj struct {
	Cache *SimTradingDBCache
}

func newSimTradingDBCache(ctx context.Context, db *SimTradingDB, params *tradingpb.SimTradingParams, hash string, buf []byte) (*SimTradingDBCache, error) {
	cache := &SimTradingDBCache{
		mapCache:      make(map[string]*tradingpb.SimTradingCacheNode),
		AssetMarket:   params.Assets[0].Market,
		AssetCode:     params.Assets[0].Code,
		StartTs:       params.StartTs,
		EndTs:         params.EndTs,
		mapHashHeader: make(map[string]*tradingpb.SimTradingCache),
	}

	// err := cache.init(ctx, db, params)
	err := cache.load(ctx, db, params, hash, buf)
	if err != nil {
		tradingdb2utils.Warn("newSimTradingDBCache:load",
			zap.Error(err))

		return nil, err
	}

	return cache, nil
}

// isMine - is mine
func (cache *SimTradingDBCache) isMine(params *tradingpb.SimTradingParams) bool {
	if cache.AssetMarket == params.Assets[0].Market &&
		cache.AssetCode == params.Assets[0].Code &&
		cache.StartTs == params.StartTs &&
		cache.EndTs == params.EndTs {

		return true
	}

	return false
}

// init - initial
func (cache *SimTradingDBCache) load(ctx context.Context, db *SimTradingDB, params *tradingpb.SimTradingParams, hash string, buf []byte) error {
	hashHeader := hash[:2]

	cache.mutexCache.Lock()
	_, isok := cache.mapHashHeader[hashHeader]
	cache.mutexCache.Unlock()
	if isok {
		return nil
	}

	pbcache, err := db.getSimTradingNodes2(ctx, params, hash, buf)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDBCache.load:getSimTradingNodes",
			zap.Error(err))

		return err
	}

	if pbcache == nil {
		return nil
	}

	cache.mutexCache.Lock()
	defer cache.mutexCache.Unlock()

	for _, v := range pbcache.Nodes {
		hash, err := cache.hashNode(v)
		if err != nil {
			tradingdb2utils.Warn("SimTradingDBCache.load:hashNode",
				zap.Error(err))

			return err
		}

		cache.mapCache[hash] = v
	}

	cache.mapHashHeader[hashHeader] = pbcache

	return nil
}

// // init - initial
// func (cache *SimTradingDBCache) init(ctx context.Context, db *SimTradingDB, params *tradingpb.SimTradingParams) error {
// 	pbcache, err := db.getSimTradingNodes(ctx, params)
// 	if err != nil {
// 		tradingdb2utils.Warn("SimTradingDBCache.init:getSimTradingNodes",
// 			zap.Error(err))

// 		return err
// 	}

// 	if pbcache == nil {
// 		return nil
// 	}

// 	cache.mutexCache.Lock()
// 	defer cache.mutexCache.Unlock()

// 	for _, v := range pbcache.Nodes {
// 		hash, err := cache.hashParams(v.Params)
// 		if err != nil {
// 			tradingdb2utils.Warn("SimTradingDBCache.init:hashParams",
// 				zap.Error(err))

// 			return err
// 		}

// 		cache.mapCache[hash] = &SimTradingDBCacheNode{
// 			node: v,
// 		}
// 	}

// 	cache.pbcache = pbcache

// 	return nil
// }

// hashNode - hash node
func (cache *SimTradingDBCache) hashNode(node *tradingpb.SimTradingCacheNode) (string, error) {
	// msg := proto.Clone(params)
	// msgParams, isok := msg.(*tradingpb.SimTradingParams)
	// if !isok {
	// 	tradingdb2utils.Warn("SimTradingDBCache.hashParams:Clone",
	// 		zap.Error(ErrInvalidSimTradingParams))

	// 	return "", ErrInvalidSimTradingParams
	// }

	// msgParams.Title = ""

	// buf, err := proto.Marshal(msgParams)
	// if err != nil {
	// 	tradingdb2utils.Warn("SimTradingDBCache.hashParams:Marshal",
	// 		zap.Error(err))

	// 	return "", err
	// }

	// h := sha1.New()
	// h.Write(buf)
	// bs := h.Sum(nil)

	return fmt.Sprintf("%x", node.Buf), nil
}

// // hashParams - hash parameters
// func (cache *SimTradingDBCache) hashParams(params *tradingpb.SimTradingParams) (string, error) {
// 	msg := proto.Clone(params)
// 	msgParams, isok := msg.(*tradingpb.SimTradingParams)
// 	if !isok {
// 		tradingdb2utils.Warn("SimTradingDBCache.hashParams:Clone",
// 			zap.Error(ErrInvalidSimTradingParams))

// 		return "", ErrInvalidSimTradingParams
// 	}

// 	msgParams.Title = ""

// 	buf, err := proto.Marshal(msgParams)
// 	if err != nil {
// 		tradingdb2utils.Warn("SimTradingDBCache.hashParams:Marshal",
// 			zap.Error(err))

// 		return "", err
// 	}

// 	// h := sha1.New()
// 	// h.Write(buf)
// 	// bs := h.Sum(nil)

// 	return fmt.Sprintf("%x", buf), nil
// }

// hashParams2 - hash parameters
func (cache *SimTradingDBCache) hashParams2(params *tradingpb.SimTradingParams, hash string, buf []byte) (string, error) {
	// msg := proto.Clone(params)
	// msgParams, isok := msg.(*tradingpb.SimTradingParams)
	// if !isok {
	// 	tradingdb2utils.Warn("SimTradingDBCache.hashParams2:Clone",
	// 		zap.Error(ErrInvalidSimTradingParams))

	// 	return "", ErrInvalidSimTradingParams
	// }

	// msgParams.Title = ""

	// buf, err := proto.Marshal(msgParams)
	// if err != nil {
	// 	tradingdb2utils.Warn("SimTradingDBCache.hashParams2:Marshal",
	// 		zap.Error(err))

	// 	return "", err
	// }

	// h := sha1.New()
	// h.Write(buf)
	// bs := h.Sum(nil)

	return fmt.Sprintf("%x", buf), nil
}

// getSimTrading - get PNLData
func (cache *SimTradingDBCache) getSimTrading(ctx context.Context, db *SimTradingDB, params *tradingpb.SimTradingParams, hash string, buf []byte) (
	*tradingpb.PNLData, error) {

	err := cache.load(ctx, db, params, hash, buf)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDBCache.getSimTrading:load",
			zap.Error(err))

		return nil, err
	}

	hashcache, err := cache.hashParams2(params, hash, buf)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDBCache.getSimTrading:hashParams2",
			zap.Error(err))

		return nil, err
	}

	cache.mutexCache.Lock()

	v, isok := cache.mapCache[hashcache]
	if !isok {
		cache.mutexCache.Unlock()

		return nil, nil
	}

	v.LastTs = time.Now().Unix()

	cache.mutexCache.Unlock()

	// err = db.updSimTradingNodes(ctx, params, cache.pbcache)
	// if err != nil {
	// 	tradingdb2utils.Warn("SimTradingDBCache.getSimTrading:updSimTradingNodes",
	// 		zap.Error(err))

	// 	return nil, err
	// }

	return db.getPNLData(ctx, v.Key)
}

// SaveCache - save cache to db on ending
func (cache *SimTradingDBCache) SaveCache(ctx context.Context, db *SimTradingDB) error {
	for k, v := range cache.mapHashHeader {
		key2 := makeSimTradingNodesDBKey2(cache.AssetMarket, cache.AssetCode, cache.StartTs, cache.EndTs, k)
		err := db.saveSimTradingNodes2(ctx, v, key2)
		if err != nil {
			tradingdb2utils.Warn("SimTradingDBCache.SaveCache:saveSimTradingNodes2",
				zap.Error(err))

			return err
		}
	}

	// err := db.updSimTradingNodesEx(ctx, cache.StrategyName, cache.AssetMarket,
	// 	cache.AssetCode, cache.StartTs, cache.EndTs, cache.pbcache)
	// if err != nil {
	// 	tradingdb2utils.Warn("SimTradingDBCache.getSimTrading:updSimTradingNodesEx",
	// 		zap.Error(err))

	// 	return err
	// }

	return nil
}

// hasSimTrading - has PNLData
func (cache *SimTradingDBCache) hasSimTrading(ctx context.Context, db *SimTradingDB, params *tradingpb.SimTradingParams, hash string, buf []byte) (bool, error) {
	err := cache.load(ctx, db, params, hash, buf)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDBCache.hasSimTrading:load",
			zap.Error(err))

		return false, err
	}

	hashcache, err := cache.hashParams2(params, hash, buf)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDBCache.hasSimTrading:hashParams2",
			zap.Error(err))

		return false, err
	}

	cache.mutexCache.Lock()

	v, isok := cache.mapCache[hashcache]
	if !isok {
		cache.mutexCache.Unlock()

		return false, nil
	}

	v.LastTs = time.Now().Unix()

	cache.mutexCache.Unlock()

	return true, nil
}

// addSimTrading - add simulation trading
func (cache *SimTradingDBCache) addSimTrading(params *tradingpb.SimTradingParams, hash string, buf []byte, node *tradingpb.SimTradingCacheNode) error {
	hashHeader := hash[:2]
	v, isok := cache.mapHashHeader[hashHeader]
	if !isok {
		return nil
	}

	hashCache, err := cache.hashParams2(params, hash, buf)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDBCache.addSimTrading:hashParams2",
			zap.Error(err))

		return err
	}

	cache.mutexCache.Lock()

	cache.mapCache[hashCache] = node

	v.Nodes = append(v.Nodes, node)
	cache.mapHashHeader[hashHeader] = v

	cache.mutexCache.Unlock()

	return nil
}
