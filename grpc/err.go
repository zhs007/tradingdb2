package tradingdb2grpc

import "errors"

var (
	// ErrNodeNotFree - Node is not free
	ErrNodeNotFree = errors.New("node is not free")

	// ErrNoNode - No node
	ErrNoNode = errors.New("no node")

	// ErrNoAsset - No asset
	ErrNoAsset = errors.New("no asset")

	// ErrNoRuningTaskInNode2Mgr - No runing task in Node2Mgr
	ErrNoRuningTaskInNode2Mgr = errors.New("no runing task in Node2Mgr")

	// ErrDuplicateRuningTaskInNode2Mgr - Duplicate runing task in Node2Mgr
	ErrDuplicateRuningTaskInNode2Mgr = errors.New("duplicate runing task in Node2Mgr")

	// ErrDuplicateTaskIndex - Duplicate task index
	ErrDuplicateTaskIndex = errors.New("duplicate task index")
	// ErrNoTaskIndex - No task index
	ErrNoTaskIndex = errors.New("no task index")

	// ErrInvalidParamsTs - Invalid timestamp in params
	ErrInvalidParamsTs = errors.New("invalid timestamp in params")
)
