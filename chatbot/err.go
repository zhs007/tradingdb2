package trdb2chatbot

import "errors"

var (
	// ErrInvalidServ - invalid serv
	ErrInvalidServ = errors.New("invalid serv")

	// ErrCmdInvalidParams - invalid params
	ErrCmdInvalidParams = errors.New("invalid params")
)
