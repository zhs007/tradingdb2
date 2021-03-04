package tradingdb2grpc

import "errors"

var (
	// ErrNodeNotFree - Node is not free
	ErrNodeNotFree = errors.New("Node is not free")

	// ErrNoNode - No node
	ErrNoNode = errors.New("No node")

	// ErrNoAsset - No asset
	ErrNoAsset = errors.New("No asset")

	// ErrNoRuningTaskInNode2Mgr - No runing task in Node2Mgr
	ErrNoRuningTaskInNode2Mgr = errors.New("No runing task in Node2Mgr")

	// ErrDuplicateRuningTaskInNode2Mgr - Duplicate runing task in Node2Mgr
	ErrDuplicateRuningTaskInNode2Mgr = errors.New("Duplicate runing task in Node2Mgr")

	// ErrDuplicateTaskIndex - Duplicate task index
	ErrDuplicateTaskIndex = errors.New("Duplicate task index")
	// ErrNoTaskIndex - No task index
	ErrNoTaskIndex = errors.New("No task index")

	// ErrInvalidParamsTs - Invalid timestamp in params
	ErrInvalidParamsTs = errors.New("Invalid timestamp in params")
)
