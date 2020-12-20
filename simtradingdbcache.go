package tradingdb2

import (
	"context"
	"crypto/sha1"
	"fmt"

	"github.com/golang/protobuf/proto"
	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
)

// SimTradingDBCacheNode -
type SimTradingDBCacheNode struct {
	node *tradingpb.SimTradingCacheNode
}

// SimTradingDBCache -
type SimTradingDBCache struct {
	mapCache     map[string]*SimTradingDBCacheNode
	StrategyName string
	AssetMarket  string
	AssetCode    string
	StartTs      int64
	EndTs        int64
	pbcache      *tradingpb.SimTradingCache
}

func newSimTradingDBCache(ctx context.Context, db *SimTradingDB, params *tradingpb.SimTradingParams) (*SimTradingDBCache, error) {
	cache := &SimTradingDBCache{
		mapCache:     make(map[string]*SimTradingDBCacheNode),
		StrategyName: params.Strategies[0].Name,
		AssetMarket:  params.Assets[0].Market,
		AssetCode:    params.Assets[0].Code,
		StartTs:      params.StartTs,
		EndTs:        params.EndTs,
	}

	err := cache.init(ctx, db, params)
	if err != nil {
		tradingdb2utils.Warn("newSimTradingDBCache:init",
			zap.Error(err))

		return nil, err
	}

	return cache, nil
}

// isMine - is mine
func (cache *SimTradingDBCache) isMine(params *tradingpb.SimTradingParams) bool {
	if cache.StrategyName == params.Strategies[0].Name &&
		cache.AssetMarket == params.Assets[0].Market &&
		cache.AssetCode == params.Assets[0].Code &&
		cache.StartTs == params.StartTs &&
		cache.EndTs == params.EndTs {

		return true
	}

	return false
}

// init - initial
func (cache *SimTradingDBCache) init(ctx context.Context, db *SimTradingDB, params *tradingpb.SimTradingParams) error {
	pbcache, err := db.getSimTradingNodes(ctx, params)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDBCache.init:getSimTradingNodes",
			zap.Error(err))

		return err
	}

	if pbcache == nil {
		return nil
	}

	for _, v := range pbcache.Nodes {
		hash, err := cache.hashParams(v.Params)
		if err != nil {
			tradingdb2utils.Warn("SimTradingDBCache.init:hashParams",
				zap.Error(err))

			return err
		}

		cache.mapCache[hash] = &SimTradingDBCacheNode{
			node: v,
		}
	}

	cache.pbcache = pbcache

	return nil
}

// hashParams - hash parameters
func (cache *SimTradingDBCache) hashParams(params *tradingpb.SimTradingParams) (string, error) {
	msg := proto.Clone(params)
	msgParams, isok := msg.(*tradingpb.SimTradingParams)
	if !isok {
		tradingdb2utils.Warn("SimTradingDBCache.hashParams:Clone",
			zap.Error(ErrInvalidSimTradingParams))

		return "", ErrInvalidSimTradingParams
	}

	msgParams.Title = ""

	buf, err := proto.Marshal(msgParams)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDBCache.hashParams:Marshal",
			zap.Error(err))

		return "", err
	}

	h := sha1.New()
	h.Write(buf)
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs), nil
}

// getSimTrading - get PNLData
func (cache *SimTradingDBCache) getSimTrading(ctx context.Context, db *SimTradingDB, params *tradingpb.SimTradingParams) (
	*tradingpb.PNLData, error) {

	hash, err := cache.hashParams(params)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDBCache.getSimTrading:hashParams",
			zap.Error(err))

		return nil, err
	}

	v, isok := cache.mapCache[hash]
	if !isok {
		return nil, nil
	}

	err = db.updSimTradingNodes(ctx, params, cache.pbcache)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDBCache.getSimTrading:updSimTradingNodes",
			zap.Error(err))

		return nil, err
	}

	return db.getPNLData(ctx, v.node.Key)
}
