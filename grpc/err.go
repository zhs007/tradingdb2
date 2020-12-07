package tradingdb2grpc

import "errors"

var (
	// ErrNodeNotFree - Node is not free
	ErrNodeNotFree = errors.New("Node is not free")

	// ErrNoNode - No node
	ErrNoNode = errors.New("No node")

	// ErrNoRuningTaskInNode2Mgr - No runing task in Node2Mgr
	ErrNoRuningTaskInNode2Mgr = errors.New("No runing task in Node2Mgr")

	// ErrDuplicateRuningTaskInNode2Mgr - Duplicate runing task in Node2Mgr
	ErrDuplicateRuningTaskInNode2Mgr = errors.New("Duplicate runing task in Node2Mgr")
)
