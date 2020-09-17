package tradingdb2

import "errors"

var (
	// ErrInvalidMarket - invalid market
	ErrInvalidMarket = errors.New("invalid market")
	// ErrInvalidSymbol - invalid symbol
	ErrInvalidSymbol = errors.New("invalid symbol")
	// ErrInvalidTag - invalid tag
	ErrInvalidTag = errors.New("invalid tag")
)
