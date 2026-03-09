package runner

import "errors"

var (
	// TaskAlreadyExistsError is returned when attempting to run a task that is
	// already registered and running.
	TaskAlreadyExistsError = errors.New("task already exists")

	// TaskUndefinedError is returned when attempting to stop a task that is not
	// registered or has already completed.
	TaskUndefinedError = errors.New("task is undefined")
)
