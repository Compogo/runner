package runner

import "errors"

// Ошибки, возвращаемые при работе с Runner.
var (
	// TaskAlreadyExistsError возникает при попытке запустить уже существующий процесс.
	TaskAlreadyExistsError = errors.New("task already exists")

	// TaskUndefinedError возникает при попытке остановить несуществующий процесс.
	TaskUndefinedError = errors.New("task is undefined")
)
