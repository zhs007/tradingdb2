package tradingdb2task

import "errors"

var (
	// ErrInvalidTask - invalid task
	ErrInvalidTask = errors.New("invalid task")

	// ErrTaskFail - task fail
	ErrTaskFail = errors.New("task fail")
)
