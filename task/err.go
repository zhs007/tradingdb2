package tradingdb2task

import "errors"

var (
	// ErrInvalidTask - invalid task
	ErrInvalidTask = errors.New("invalid task")
	// ErrInvalidTaskKey - invalid task key
	ErrInvalidTaskKey = errors.New("invalid task key")

	// ErrTaskFail - task fail
	ErrTaskFail = errors.New("task fail")
)
