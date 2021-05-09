package tradingdb2grpc

import (
	"context"
	"io"
	"net"
	"sort"

	tradingdb2 "github.com/zhs007/tradingdb2"
	tradingdb2task "github.com/zhs007/tradingdb2/task"
	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	tradingdb2ver "github.com/zhs007/tradingdb2/ver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Serv - tradingdb2 Service
type Serv struct {
	tradingpb.UnimplementedTradingDB2Server
	lis           net.Listener
	grpcServ      *grpc.Server
	DB2           *tradingdb2.DB2
	Cfg           *tradingdb2.Config
	DBSimTrading  *tradingdb2.SimTradingDB
	SimTradingDB2 *tradingdb2.SimTradingDB2
	MgrNodes      *Node2Mgr
	tasksMgr      *tradingdb2task.TasksMgr
}

// NewServ -
func NewServ(cfg *tradingdb2.Config) (*Serv, error) {

	db, err := tradingdb2.NewDB2(cfg)
	if err != nil {
		tradingdb2utils.Error("NewServ.NewDB2",
			zap.Error(err))

		return nil, err
	}

	dbSimTrading, err := tradingdb2.NewSimTradingDB(cfg.DBPath, "", cfg.DBEngine)
	if err != nil {
		tradingdb2utils.Error("NewServ.NewSimTradingDB",
			zap.Error(err))

		return nil, err
	}

	simTradingDB2, err := tradingdb2.NewSimTradingDB2(cfg.DBPath, "", cfg.DBEngine)
	if err != nil {
		tradingdb2utils.Error("NewServ.NewSimTradingDB2",
			zap.Error(err))

		return nil, err
	}

	tasksMgr := tradingdb2task.NewTasksMgr()
	// if err != nil {
	// 	tradingdb2utils.Error("NewServ.NewSimTradingDB2",
	// 		zap.Error(err))

	// 	return nil, err
	// }

	lis, err := net.Listen("tcp", cfg.BindAddr)
	if err != nil {
		tradingdb2utils.Error("NewServ.Listen",
			zap.Error(err))

		return nil, err
	}

	grpcServ := grpc.NewServer()

	mgrNodes, err := NewNode2Mgr(cfg)
	if err != nil {
		tradingdb2utils.Error("NewServ.NewNode2Mgr",
			zap.Error(err))

		return nil, err
	}

	serv := &Serv{
		lis:           lis,
		grpcServ:      grpcServ,
		DB2:           db,
		DBSimTrading:  dbSimTrading,
		Cfg:           cfg,
		MgrNodes:      mgrNodes,
		SimTradingDB2: simTradingDB2,
		tasksMgr:      tasksMgr,
	}

	tradingpb.RegisterTradingDB2Server(grpcServ, serv)

	tradingdb2utils.Info("NewServ OK.",
		zap.String("addr", cfg.BindAddr),
		zap.String("ver", tradingdb2ver.Version))

	return serv, nil
}

// Start - start a service
func (serv *Serv) Start(ctx context.Context) error {
	go func() {
		go serv.MgrNodes.Start()
	}()

	return serv.grpcServ.Serve(serv.lis)
}

// Stop - stop service
func (serv *Serv) Stop() {
	serv.lis.Close()

	serv.MgrNodes.Stop()
	serv.DBSimTrading.Stop()
}

// checkBasicRequest - check BasicRequest
func (serv *Serv) checkBasicRequest(req *tradingpb.BasicRequestData) error {
	if req != nil &&
		(req.Token == "" || tradingdb2utils.IndexOfStringSlice(serv.Cfg.Tokens, req.Token, 0) < 0) {

		return tradingdb2.ErrInvalidToken
	}

	return nil
}

// UpdCandles - update candles
func (serv *Serv) UpdCandles(stream tradingpb.TradingDB2_UpdCandlesServer) error {
	candles := &tradingpb.Candles{}
	times := 0

	for {
		req, err := stream.Recv()
		if req != nil && (err == nil || err == io.EOF) {
			if times == 0 {
				err = serv.checkBasicRequest(req.BasicRequest)
				if err != nil {
					tradingdb2utils.Error("Serv.UpdCandles:Token",
						zap.Int("length", len(candles.Candles)),
						zap.Int("times", times),
						zap.String("token", req.BasicRequest.Token),
						zap.Error(err))

					return err
				}

				if req.Candles.Market == "" ||
					req.Candles.Symbol == "" ||
					req.Candles.Tag == "" {

					tradingdb2utils.Error("Serv.UpdCandles:Market|Symbol|Tag",
						zap.Int("length", len(candles.Candles)),
						zap.Int("times", times),
						zap.String("market", req.Candles.Market),
						zap.String("symbol", req.Candles.Symbol),
						zap.String("tag", req.Candles.Tag),
						zap.Error(tradingdb2.ErrInvalidUpdCandlesParams))

					return tradingdb2.ErrInvalidUpdCandlesParams
				}

				candles.Market = req.Candles.Market
				candles.Symbol = req.Candles.Symbol
				candles.Tag = req.Candles.Tag
			}

			times++

			tradingdb2.MergeCandles(candles, req.Candles)
		}

		if err != nil {
			if err == io.EOF {
				err := serv.DB2.UpdCandles(stream.Context(), candles)
				if err != nil {
					tradingdb2utils.Error("Serv.UpdCandles:DB.UpdCandles",
						zap.Int("length", len(candles.Candles)),
						zap.Int("times", times),
						zap.Error(err))

					return err
				}

				return stream.SendAndClose(&tradingpb.ReplyUpdCandles{
					LengthOK: int32(len(candles.Candles)),
				})
			}

			tradingdb2utils.Error("Serv.UpdCandles:Recv",
				zap.Int("length", len(candles.Candles)),
				zap.Int("times", times),
				zap.Error(err))

			return err
		}
	}
}

// GetCandles - get candles
func (serv *Serv) GetCandles(req *tradingpb.RequestGetCandles, stream tradingpb.TradingDB2_GetCandlesServer) error {
	err := serv.checkBasicRequest(req.BasicRequest)
	if err != nil {
		tradingdb2utils.Error("Serv.GetCandles:checkToken",
			zap.String("token", req.BasicRequest.Token),
			zap.Strings("tokens", serv.Cfg.Tokens),
			zap.Error(err))

		return err
	}

	// var tags []string
	// if len(req.Tags) > 0 {
	// 	tags = req.Tags
	// } else if len(req.Tags) == 0 && req.Tag != "" {
	// 	tags = append(tags, req.Tag)
	// }

	candles, err := serv.DB2.GetCandles(stream.Context(), req.Market, req.Symbol, req.TsStart, req.TsEnd, req.Offset)
	if err != nil {
		tradingdb2utils.Error("Serv.GetCandles:DB.GetCandles",
			tradingdb2utils.JSON("params", req),
			zap.Error(err))

		return err
	}

	// sentnums := 0

	return tradingdb2.BatchCandles(candles, serv.Cfg.BatchCandleNums, func(lst *tradingpb.Candles) error {
		cc := &tradingpb.ReplyGetCandles{
			Candles: lst,
		}

		// if sentnums == 0 {
		// 	cc.Candles.Market = candles.Market
		// 	cc.Candles.Symbol = candles.Symbol
		// 	cc.Candles.Tag = candles.Tag
		// }

		return stream.Send(cc)
	})
}

// UpdSymbol - update symbol
func (serv *Serv) UpdSymbol(ctx context.Context, req *tradingpb.RequestUpdSymbol) (*tradingpb.ReplyUpdSymbol, error) {
	err := serv.checkBasicRequest(req.BasicRequest)
	if err != nil {
		tradingdb2utils.Error("Serv.UpdSymbol:checkToken",
			zap.String("token", req.BasicRequest.Token),
			zap.Strings("tokens", serv.Cfg.Tokens),
			zap.Error(err))

		return nil, err
	}

	symbol, err := serv.DB2.GetSymbol(ctx, req.Symbol.Market, req.Symbol.Symbol)
	if err != nil {
		tradingdb2utils.Error("Serv.UpdSymbol:DB.GetSymbol",
			zap.Error(err))

		return nil, err
	}

	if symbol != nil && symbol.Fund != nil && req.Symbol.Fund != nil {
		symbol.Fund = tradingdb2.MergeFund(symbol.Fund, req.Symbol.Fund)
	} else {
		symbol = req.Symbol
	}

	err = serv.DB2.UpdSymbol(ctx, symbol)
	if err != nil {
		tradingdb2utils.Error("Serv.UpdSymbol:DB.UpdSymbol",
			zap.Error(err))

		return nil, err
	}

	res := &tradingpb.ReplyUpdSymbol{
		IsOK: true,
	}

	return res, nil
}

// GetSymbol - get symbol
func (serv *Serv) GetSymbol(ctx context.Context, req *tradingpb.RequestGetSymbol) (*tradingpb.ReplyGetSymbol, error) {
	err := serv.checkBasicRequest(req.BasicRequest)
	if err != nil {
		tradingdb2utils.Error("Serv.GetSymbol:checkToken",
			zap.String("token", req.BasicRequest.Token),
			zap.Strings("tokens", serv.Cfg.Tokens),
			zap.Error(err))

		return nil, err
	}

	si, err := serv.DB2.GetSymbol(ctx, req.Market, req.Symbol)
	if err != nil {
		tradingdb2utils.Error("Serv.GetSymbol:DB.GetSymbol",
			zap.String("market", req.Market),
			zap.String("sytmbol", req.Symbol),
			zap.Error(err))

		return nil, err
	}

	res := &tradingpb.ReplyGetSymbol{
		Symbol: si,
	}

	return res, nil
}

// GetSymbols - get symbols
func (serv *Serv) GetSymbols(req *tradingpb.RequestGetSymbols, stream tradingpb.TradingDB2_GetSymbolsServer) error {
	err := serv.checkBasicRequest(req.BasicRequest)
	if err != nil {
		tradingdb2utils.Error("Serv.GetSymbols:checkToken",
			zap.String("token", req.BasicRequest.Token),
			zap.Strings("tokens", serv.Cfg.Tokens),
			zap.Error(err))

		return err
	}

	var symbols []string
	if len(req.Symbols) > 0 {
		symbols = req.Symbols
	} else {
		arr, err := serv.DB2.GetMarketSymbols(stream.Context(), req.Market)
		if err != nil {
			tradingdb2utils.Error("Serv.GetSymbols:GetMarketSymbols",
				zap.String("Market", req.Market),
				zap.Error(err))

			return err
		}

		symbols = arr
	}

	for _, v := range symbols {
		si, err := serv.DB2.GetSymbol(stream.Context(), req.Market, v)
		if err != nil {
			tradingdb2utils.Error("Serv.GetSymbols:DB.GetSymbol",
				zap.String("market", req.Market),
				zap.String("sytmbol", v),
				zap.Error(err))

			// return err
		}

		if si != nil {
			res := &tradingpb.ReplyGetSymbol{
				Symbol: si,
			}

			err = stream.Send(res)
			if err != nil {
				tradingdb2utils.Error("Serv.GetSymbols:Send",
					zap.String("sytmbol", v),
					zap.Error(err))

				// return err
			}
		}
	}

	return nil
}

// SimTrading - simTrading
func (serv *Serv) checkParams(params *tradingpb.SimTradingParams) error {
	if params.StartTs > 0 && params.EndTs > 0 {
		if params.StartTs >= params.EndTs {
			return ErrInvalidParamsTs
		}
	}

	return nil
}

// SimTrading - simTrading
func (serv *Serv) SimTrading(ctx context.Context, req *tradingpb.RequestSimTrading) (*tradingpb.ReplySimTrading, error) {
	tradingdb2utils.Info("Serv.SimTrading",
		tradingdb2utils.JSON("request", req))

	err := serv.checkBasicRequest(req.BasicRequest)
	if err != nil {
		tradingdb2utils.Error("Serv.SimTrading:checkToken",
			zap.String("token", req.BasicRequest.Token),
			zap.Strings("tokens", serv.Cfg.Tokens),
			zap.Error(err))

		return nil, err
	}

	for _, asset := range req.Params.Assets {
		if !serv.DB2.HasCandles(ctx, asset.Market, asset.Code) {
			tradingdb2utils.Error("Serv.SimTrading:HasCandles",
				tradingdb2utils.JSON("asset", asset),
				zap.Error(ErrNoAsset))

			return nil, ErrNoAsset
		}
	}

	err = serv.checkParams(req.Params)
	if err != nil {
		tradingdb2utils.Error("Serv.SimTrading:checkParams",
			tradingdb2utils.JSON("params", req.Params),
			zap.Error(err))

		return nil, err
	}

	params, err := serv.DB2.FixSimTradingParams(ctx, req.Params)
	if err != nil {
		tradingdb2utils.Error("Serv.SimTrading:FixSimTradingParams",
			zap.Error(err))

		return nil, err
	}

	if !req.IgnoreCache {
		pnl, err := serv.DBSimTrading.GetSimTrading(ctx, params)
		if err != nil {
			tradingdb2utils.Error("Serv.SimTrading:GetSimTrading",
				zap.Error(err))

			return nil, err
		}

		if pnl != nil {
			tradingdb2utils.Debug("Serv.SimTrading:Cached")

			pnl.Title = req.Params.Title

			return &tradingpb.ReplySimTrading{
				Pnl: []*tradingpb.PNLData{
					pnl,
				},
			}, nil
		}
	}

	res := &tradingpb.ReplySimTrading{}

	// for _, v := range params.Baselines {
	// 	bp := GenCalcBaseline(v, params.StartTs, params.EndTs)
	// 	br, err := serv.MgrNodes.CalcPNL2(ctx, bp, nil)
	// 	if err != nil {
	// 		tradingdb2utils.Error("Serv.SimTrading:CalcPNL baseline",
	// 			tradingdb2utils.JSON("parameter", bp),
	// 			zap.Error(err))

	// 		return nil, err
	// 	}

	// 	res.Baseline = append(res.Baseline, br.Pnl...)
	// }

	reply, err := serv.MgrNodes.CalcPNL2(ctx, params, nil)
	if err != nil {
		tradingdb2utils.Error("Serv.SimTrading:CalcPNL2",
			zap.Error(err))

		return nil, err
	}

	res.Pnl = reply.Pnl

	if len(reply.Pnl) > 0 {
		reply.Pnl[0].Title = req.Params.Title

		err = serv.DBSimTrading.UpdSimTrading(ctx, params, reply.Pnl[0])
		if err != nil {
			tradingdb2utils.Error("Serv.SimTrading:UpdSimTrading",
				zap.Error(err))

			return nil, err
		}
	}

	return res, nil
}

// procIgnoreReply - simulation trading
func (serv *Serv) procIgnoreReply(lstIgnore []*tradingpb.ReplySimTrading, minNums int, maxIgnoreNums int, isend bool) (
	[]*tradingpb.ReplySimTrading, []*tradingpb.ReplySimTrading) {

	if len(lstIgnore) <= minNums {
		return lstIgnore, nil
	}

	sort.SliceStable(lstIgnore, func(i, j int) bool {
		if len(lstIgnore[i].Pnl) <= 0 || len(lstIgnore[j].Pnl) <= 0 {
			return true
		}

		if len(lstIgnore[i].Pnl[0].Total.LstCtrl) <= 0 {
			return false
		}

		if len(lstIgnore[j].Pnl[0].Total.LstCtrl) <= 0 {
			return true
		}

		return lstIgnore[i].Pnl[0].Total.TotalReturns > lstIgnore[j].Pnl[0].Total.TotalReturns
	})

	if maxIgnoreNums > 0 {
		if isend {
			var newlst []*tradingpb.ReplySimTrading
			var ignorelst []*tradingpb.ReplySimTrading

			for i, v := range lstIgnore {
				if i < minNums {
					newlst = append(newlst, v)
				} else if i-minNums < maxIgnoreNums {
					ignorelst = append(ignorelst, v)
				} else {
					break
				}
			}

			return newlst, ignorelst
		}

		var newlst []*tradingpb.ReplySimTrading

		for i, v := range lstIgnore {
			if i < minNums {
				newlst = append(newlst, v)
			} else if i-minNums < maxIgnoreNums {
				setIgnoreReplySimTrading(v)

				newlst = append(newlst, v)
			} else {
				break
			}
		}

		return newlst, nil
	}

	var newlst []*tradingpb.ReplySimTrading
	for i := 0; i < minNums; i++ {
		newlst = append(newlst, lstIgnore[i])
	}

	return newlst, nil
}

// SimTrading2 - simulation trading
func (serv *Serv) SimTrading2(stream tradingpb.TradingDB2_SimTrading2Server) error {
	stt := NewSimTradingTasksMgr()
	dbcache := &tradingdb2.SimTradingDBCacheObj{}

	// ignoreTotalReturn := 1.0 // 忽略总回报低于这个值的数据返回，主要用于大批量训练，减少数据处理量。但实际运算会执行，且cache数据是完整的
	minNums := 10 // 被忽略的数据里，还是保留这个条数返回，默认是10
	maxIgnoreNums := 0
	// sortBy := "totalreturn" // 按什么字段来排序，默认是 totalreturn

	var lstIgnore []*tradingpb.ReplySimTrading

	stt.Start()

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			tradingdb2utils.Debug("Serv.SimTrading2:EOF")

			stt.Stop()

			if len(lstIgnore) > minNums {
				lstlast, lstlost := serv.procIgnoreReply(lstIgnore, minNums, maxIgnoreNums, true)

				for _, v := range lstlost {
					setIgnoreReplySimTrading(v)

					err := stream.Send(v)
					if err != nil {
						tradingdb2utils.Error("Serv.SimTrading2:simTrading:SendIgnore lost",
							zap.Error(err))
					}
				}

				lstIgnore = lstlast
			}

			for _, v := range lstIgnore {
				err := stream.Send(v)
				if err != nil {
					tradingdb2utils.Error("Serv.SimTrading2:simTrading:SendIgnore last",
						zap.Error(err))
				}
			}

			if dbcache.Cache != nil {
				err = dbcache.Cache.SaveCache(stream.Context(), serv.DBSimTrading)
				if err != nil {
					tradingdb2utils.Error("Serv.SimTrading2:simTrading:SaveCache",
						zap.Error(err))
				}
			}

			tradingdb2utils.Debug("Serv.SimTrading2:End")

			return nil
		}

		if err != nil {
			tradingdb2utils.Error("Serv.SimTrading2",
				zap.Error(err))

			stt.Stop()

			return err
		}

		if in != nil {
			// if in.IgnoreTotalReturn != 0 {
			// 	ignoreTotalReturn = float64(in.IgnoreTotalReturn)
			// }

			if in.MinNums > 0 {
				minNums = int(in.MinNums)
			}

			if in.MaxIgnoreNums > 0 {
				maxIgnoreNums = int(in.MaxIgnoreNums)
			}

			// if in.SortBy != "" {
			// 	sortBy = in.SortBy
			// }

			// 这个接口不是阻塞的，错误没法直接传递到外面来

			serv.simTrading(stream.Context(), stt, dbcache, in,
				func(req *tradingpb.RequestSimTrading, reply *tradingpb.ReplySimTrading, err error, inCache bool, dbcache *tradingdb2.SimTradingDBCacheObj) {
					if err != nil {
						tradingdb2utils.Error("Serv.SimTrading2:simTrading:OnEnd",
							zap.Error(err))

						// errST = err
					} else if reply != nil {
						isSendNow := true

						if len(reply.Pnl) > 0 {
							reply.Pnl[0].Title = req.Params.Title

							if !inCache {
								// 如果不是incache，才需要更新缓存
								newdbcache, err := serv.DBSimTrading.UpdSimTradingEx(stream.Context(), req.Params, reply.Pnl[0], dbcache.Cache)
								dbcache.Cache = newdbcache

								if err != nil {
									tradingdb2utils.Error("Serv.SimTrading2:UpdSimTrading",
										zap.Error(err))

									return
								}
							}

							// 不发明细
							if req.IgnoreTotalReturn > 0 && reply.Pnl[0].Total != nil && reply.Pnl[0].Total.TotalReturns < req.IgnoreTotalReturn {
								isSendNow = false

								lstIgnore = append(lstIgnore, reply)
								if len(lstIgnore) > minNums*10 {
									// if len(lstIgnore) > minNums*10 && len(lstIgnore) > maxIgnoreNums {
									lstlast, _ := serv.procIgnoreReply(lstIgnore, minNums, maxIgnoreNums, false)

									// for _, v := range lstlost {
									// 	setIgnoreReplySimTrading(v)

									// 	err := stream.Send(v)
									// 	if err != nil {
									// 		tradingdb2utils.Error("Serv.SimTrading2:simTrading:SendIgnore onend lost",
									// 			zap.Error(err))
									// 	}
									// }

									lstIgnore = lstlast
								}
							}
						}

						if isSendNow {
							err := stream.Send(reply)
							if err != nil {
								tradingdb2utils.Error("Serv.SimTrading2:simTrading:Send",
									zap.Error(err))
							}
						}
					}
				})
		}
	}
}

// simTrading - simTrading
func (serv *Serv) simTrading(ctx context.Context, mgrTasks *SimTradingTasksMgr, dbcache *tradingdb2.SimTradingDBCacheObj,
	req *tradingpb.RequestSimTrading, onEnd FuncOnSimTradingTaskEnd) {
	// tradingdb2utils.Info("Serv.simTrading",
	// 	tradingdb2utils.JSON("request", req))

	err := serv.checkBasicRequest(req.BasicRequest)
	if err != nil {
		tradingdb2utils.Error("Serv.simTrading:checkToken",
			zap.String("token", req.BasicRequest.Token),
			zap.Strings("tokens", serv.Cfg.Tokens),
			zap.Error(err))

		onEnd(req, nil, err, false, dbcache)

		return
	}

	for _, asset := range req.Params.Assets {
		if !serv.DB2.HasCandles(ctx, asset.Market, asset.Code) {
			tradingdb2utils.Error("Serv.simTrading:HasCandles",
				tradingdb2utils.JSON("asset", asset),
				zap.Error(ErrNoAsset))

			onEnd(req, nil, ErrNoAsset, false, dbcache)

			return
		}
	}

	err = serv.checkParams(req.Params)
	if err != nil {
		tradingdb2utils.Error("Serv.simTrading:checkParams",
			tradingdb2utils.JSON("params", req.Params),
			zap.Error(err))

		onEnd(req, nil, err, false, dbcache)

		return
	}

	params, err := serv.DB2.FixSimTradingParams(ctx, req.Params)
	if err != nil {
		tradingdb2utils.Error("Serv.simTrading:FixSimTradingParams",
			zap.Error(err))

		onEnd(req, nil, err, false, dbcache)

		return
	}

	if !req.IgnoreCache {
		newdbcache, pnl, err := serv.DBSimTrading.GetSimTradingEx(ctx, params, dbcache.Cache)
		dbcache.Cache = newdbcache

		if err != nil {
			tradingdb2utils.Error("Serv.simTrading:GetSimTrading",
				zap.Error(err))

			onEnd(req, nil, err, false, dbcache)

			return
		}

		if pnl != nil {
			// tradingdb2utils.Debug("Serv.simTrading:Cached")

			onEnd(req, &tradingpb.ReplySimTrading{
				Pnl: []*tradingpb.PNLData{
					pnl,
				},
			}, err, true, dbcache)

			return
		}
	}

	err = mgrTasks.AddTask(serv.MgrNodes, dbcache, req, onEnd)
	if err != nil {
		tradingdb2utils.Error("Serv.simTrading:AddTask",
			zap.Error(err))
	}
}

// simTrading3 - simTrading3
func (serv *Serv) simTrading3(ctx context.Context, req *tradingpb.RequestSimTrading, onEnd FuncOnSimTrading3TaskEnd) {
	// tradingdb2utils.Info("Serv.simTrading",
	// 	tradingdb2utils.JSON("request", req))

	err := serv.checkBasicRequest(req.BasicRequest)
	if err != nil {
		tradingdb2utils.Error("Serv.simTrading3:checkToken",
			zap.String("token", req.BasicRequest.Token),
			zap.Strings("tokens", serv.Cfg.Tokens),
			zap.Error(err))

		onEnd(req, nil, err, false)

		return
	}

	for _, asset := range req.Params.Assets {
		if !serv.DB2.HasCandles(ctx, asset.Market, asset.Code) {
			tradingdb2utils.Error("Serv.simTrading3:HasCandles",
				tradingdb2utils.JSON("asset", asset),
				zap.Error(ErrNoAsset))

			onEnd(req, nil, ErrNoAsset, false)

			return
		}
	}

	err = serv.checkParams(req.Params)
	if err != nil {
		tradingdb2utils.Error("Serv.simTrading:checkParams",
			tradingdb2utils.JSON("params", req.Params),
			zap.Error(err))

		onEnd(req, nil, err, false)

		return
	}

	params, err := serv.DB2.FixSimTradingParams(ctx, req.Params)
	if err != nil {
		tradingdb2utils.Error("Serv.simTrading:FixSimTradingParams",
			zap.Error(err))

		onEnd(req, nil, err, false)

		return
	}

	if !req.IgnoreCache {
		pnl, err := serv.SimTradingDB2.GetSimTrading(ctx, params)

		if err != nil {
			tradingdb2utils.Error("Serv.simTrading3:GetSimTrading",
				zap.Error(err))

			onEnd(req, nil, err, false)

			return
		}

		if pnl != nil {
			// tradingdb2utils.Debug("Serv.simTrading:Cached")

			onEnd(req, &tradingpb.ReplySimTrading{
				Pnl: []*tradingpb.PNLData{
					pnl,
				},
			}, err, true)

			return
		}

		//！！ 这里开始是兼容SimTradingDB的
		isok, err := serv.DBSimTrading.Upgrade2SimTradingDB2(ctx, params, serv.SimTradingDB2)
		if err != nil {
			tradingdb2utils.Error("Serv.simTrading3:Upgrade2SimTradingDB2",
				zap.Error(err))

			onEnd(req, nil, err, false)

			return
		}

		if isok {
			pnl, err = serv.SimTradingDB2.GetSimTrading(ctx, params)
			if err != nil {
				tradingdb2utils.Error("Serv.simTrading3:GetSimTrading&Upgrade2SimTradingDB2",
					zap.Error(err))

				onEnd(req, nil, err, false)

				return
			}

			if pnl != nil {
				// tradingdb2utils.Debug("Serv.simTrading:Cached")

				onEnd(req, &tradingpb.ReplySimTrading{
					Pnl: []*tradingpb.PNLData{
						pnl,
					},
				}, err, true)

				return
			}
		}
	}

	err = serv.tasksMgr.AddTask(req.Params, func(task *tradingdb2task.Task) error {
		onEnd(req, &tradingpb.ReplySimTrading{
			Pnl: []*tradingpb.PNLData{
				task.PNL,
			},
		}, nil, false)

		return nil
	})
	if err != nil {
		tradingdb2utils.Error("Serv.simTrading3:AddTask",
			zap.Error(err))
	}
}

// SimTrading3 - simulation trading
func (serv *Serv) SimTrading3(stream tradingpb.TradingDB2_SimTrading3Server) error {
	minNums := 10 // 被忽略的数据里，还是保留这个条数返回，默认是10
	maxIgnoreNums := 0
	var lstIgnore []*tradingpb.ReplySimTrading

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			tradingdb2utils.Debug("Serv.SimTrading3:EOF")

			if len(lstIgnore) > minNums {
				lstlast, lstlost := serv.procIgnoreReply(lstIgnore, minNums, maxIgnoreNums, true)

				for _, v := range lstlost {
					setIgnoreReplySimTrading(v)

					err := stream.Send(v)
					if err != nil {
						tradingdb2utils.Error("Serv.SimTrading3:simTrading:SendIgnore lost",
							zap.Error(err))
					}
				}

				lstIgnore = lstlast
			}

			for _, v := range lstIgnore {
				err := stream.Send(v)
				if err != nil {
					tradingdb2utils.Error("Serv.SimTrading3:simTrading:SendIgnore last",
						zap.Error(err))
				}
			}

			tradingdb2utils.Debug("Serv.SimTrading3:End")

			return nil
		}

		if err != nil {
			tradingdb2utils.Error("Serv.SimTrading3",
				zap.Error(err))

			return err
		}

		if in != nil {
			if in.MinNums > 0 {
				minNums = int(in.MinNums)
			}

			if in.MaxIgnoreNums > 0 {
				maxIgnoreNums = int(in.MaxIgnoreNums)
			}

			// 这个接口不是阻塞的，错误没法直接传递到外面来

			serv.simTrading3(stream.Context(), in,
				func(req *tradingpb.RequestSimTrading, reply *tradingpb.ReplySimTrading, err error, inCache bool) {
					if err != nil {
						tradingdb2utils.Error("Serv.SimTrading3:simTrading:OnEnd",
							zap.Error(err))

						// errST = err
					} else if reply != nil {
						isSendNow := true

						if len(reply.Pnl) > 0 {
							reply.Pnl[0].Title = req.Params.Title

							if !inCache {
								// 如果不是incache，才需要更新缓存
								err := serv.SimTradingDB2.UpdSimTrading(stream.Context(), req.Params, reply.Pnl[0])
								if err != nil {
									tradingdb2utils.Error("Serv.SimTrading3:UpdSimTrading",
										zap.Error(err))

									return
								}
							}

							// 不发明细
							if req.IgnoreTotalReturn > 0 && reply.Pnl[0].Total != nil && reply.Pnl[0].Total.TotalReturns < req.IgnoreTotalReturn {
								isSendNow = false

								lstIgnore = append(lstIgnore, reply)
								if len(lstIgnore) > minNums*10 {
									// if len(lstIgnore) > minNums*10 && len(lstIgnore) > maxIgnoreNums {
									lstlast, _ := serv.procIgnoreReply(lstIgnore, minNums, maxIgnoreNums, false)

									// for _, v := range lstlost {
									// 	setIgnoreReplySimTrading(v)

									// 	err := stream.Send(v)
									// 	if err != nil {
									// 		tradingdb2utils.Error("Serv.SimTrading2:simTrading:SendIgnore onend lost",
									// 			zap.Error(err))
									// 	}
									// }

									lstIgnore = lstlast
								}
							}
						}

						if isSendNow {
							err := stream.Send(reply)
							if err != nil {
								tradingdb2utils.Error("Serv.UpdSimTrading:simTrading:Send",
									zap.Error(err))
							}
						}
					}
				})
		}
	}
}

// ReqTradingTask3 - request trading task
func (serv *Serv) ReqTradingTask3(stream tradingpb.TradingDB2_ReqTradingTask3Server) error {
	return nil
}
