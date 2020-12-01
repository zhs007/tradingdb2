package tradingdb2

import (
	"context"
	"crypto/sha1"
	"fmt"
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
	cache, err := db.GetSimTradingNodes(ctx, params)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.UpdSimTrading:GetSimTradingNodes",
			zap.Error(err))

		return err
	}

	if cache == nil {
		cache = &tradingpb.SimTradingCache{}
	}

	for _, v := range cache.Nodes {
		if db.isSameSimTradingParams(params, v.Params) {
			v.LastTs = time.Now().Unix()
			return db.UpdSimTradingNodes(ctx, params, cache)
		}
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
	}

	cache.Nodes = append(cache.Nodes, node)

	err = db.UpdSimTradingNodes(ctx, params, cache)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.UpdSimTrading:UpdSimTradingNodes",
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

// GetSimTrading - get candles
func (db *SimTradingDB) GetSimTrading(ctx context.Context, params *tradingpb.SimTradingParams) (
	*tradingpb.PNLData, error) {

	cache, err := db.GetSimTradingNodes(ctx, params)
	if err != nil {
		tradingdb2utils.Warn("SimTradingDB.GetSimTrading:GetSimTradingNodes",
			zap.Error(err))

		return nil, err
	}

	if cache == nil {
		return nil, nil
	}

	for _, v := range cache.Nodes {
		if db.isSameSimTradingParams(params, v.Params) {
			v.LastTs = time.Now().Unix()
			db.UpdSimTradingNodes(ctx, params, cache)

			return db.getPNLData(ctx, v.Key)
		}
	}

	return nil, nil
}

// GetSimTradingNodes - get simtrading nodes
func (db *SimTradingDB) GetSimTradingNodes(ctx context.Context, params *tradingpb.SimTradingParams) (
	*tradingpb.SimTradingCache, error) {

	key := makeSimTradingNodesDBKey(params.Strategies[0].Name, params.Assets[0].Market, params.Assets[0].Code, params.StartTs, params.EndTs)

	buf, err := db.AnkaDB.Get(ctx, simtradingDBName, key)
	if err != nil {
		if err == ankadb.ErrNotFoundKey {
			return nil, nil
		}

		tradingdb2utils.Warn("SimTradingDB.GetSimTradingNodes:Get",
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

// UpdSimTradingNodes - update simtrading nodes
func (db *SimTradingDB) UpdSimTradingNodes(ctx context.Context, params *tradingpb.SimTradingParams, cache *tradingpb.SimTradingCache) error {
	key := makeSimTradingNodesDBKey(params.Strategies[0].Name, params.Assets[0].Market, params.Assets[0].Code, params.StartTs, params.EndTs)

	buf, err := proto.Marshal(cache)
	if err != nil {
		return err
	}

	err = db.AnkaDB.Set(ctx, simtradingDBName, key, buf)
	if err != nil {
		return err
	}

	return nil
}

// isSameSimTradingParams - is same SimTradingParams
func (db *SimTradingDB) isSameSimTradingParams(v0 *tradingpb.SimTradingParams, v1 *tradingpb.SimTradingParams) bool {
	if len(v0.Strategies) == len(v1.Strategies) {
		for si, sv0 := range v0.Strategies {
			sv1 := v1.Strategies[si]
			if !IsSameStrategy(sv0, sv1) {
				return false
			}
		}
	}

	return true
}

// // GetPNLData - get simtrading PNLData
// func (db *SimTradingDB) GetPNLData(ctx context.Context, params *tradingpb.SimTradingParams) (
// 	*tradingpb.PNLData, error) {

// 	cache, err := db.GetSimTradingNodes(ctx, params)
// 	if err != nil {
// 		tradingdb2utils.Warn("SimTradingDB.GetPNLData:GetSimTradingNodes",
// 			zap.Error(err))

// 		return nil, err
// 	}

// 	if cache == nil {
// 		return nil, nil
// 	}

// 	for _, v := range cache.Nodes {
// 		if db.isSameSimTradingParams(params, v.Params) {
// 			v.LastTs = time.Now().Unix()
// 			db.UpdSimTradingNodes(ctx, params, cache)

// 			return db.getPNLData(ctx, v.Key)
// 		}
// 	}

// 	return nil, nil
// }

// getPNLData - get candles
func (db *SimTradingDB) getPNLData(ctx context.Context, key string) (
	*tradingpb.PNLData, error) {

	buf, err := db.AnkaDB.Get(ctx, simtradingDBName, key)
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

	err = db.AnkaDB.Set(ctx, simtradingDBName, key, buf)
	if err != nil {
		return err
	}

	return nil
}
