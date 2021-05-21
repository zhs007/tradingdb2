package tradingdb2task

import "errors"

var (
	// ErrInvalidTask - invalid task
	ErrInvalidTask = errors.New("invalid task")
	// ErrInvalidTaskKey - invalid task key
	ErrInvalidTaskKey = errors.New("invalid task key")

	// ErrDuplicateKey - duplicate key
	ErrDuplicateKey = errors.New("duplicate key")

	// ErrNoKey - non-key
	ErrNoKey = errors.New("non-key")

	// ErrTaskFail - task fail
	ErrTaskFail = errors.New("task fail")
)
