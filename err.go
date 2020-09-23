package tradingdb2

import "errors"

var (
	// ErrInvalidMarket - invalid market
	ErrInvalidMarket = errors.New("invalid market")
	// ErrInvalidSymbol - invalid symbol
	ErrInvalidSymbol = errors.New("invalid symbol")
	// ErrInvalidTag - invalid tag
	ErrInvalidTag = errors.New("invalid tag")

	// ErrInvalidUpdCandlesParams - invalid UpdCandles Params
	ErrInvalidUpdCandlesParams = errors.New("invalid UpdCandles Params")

	// ErrDuplicateChildNode - duplicate child node
	ErrDuplicateChildNode = errors.New("duplicate child node")

	// ErrInvalidToken - invalid token
	ErrInvalidToken = errors.New("invalid token")
)
