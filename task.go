package runner

import "context"

// Task represents a named executable unit managed by the Runner.
// Each Task has its own lifecycle and can be started and stopped independently.
type Task struct {
	name string

	ctx        context.Context
	cancelFunc context.CancelFunc
	process    Process
}

// NewTask creates a new Task with the given name and process.
// The name is used for logging and stopping the task by name.
func NewTask(name string, process Process) *Task {
	return &Task{name: name, process: process}
}

// String returns the task's name, implementing the fmt.Stringer interface.
func (task *Task) String() string {
	return task.Name()
}

// Name returns the task's identifier.
func (task *Task) Name() string {
	return task.name
}
