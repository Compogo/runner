package runner

// Middleware defines the interface for task middleware.
// Middlewares can wrap a task's Process function to add cross-cutting concerns
// such as logging, metrics, tracing, panic recovery, or authentication.
//
// Middlewares are applied in a chain, where each middleware receives the task
// and the next Process in the chain, and returns a new Process that may
// perform operations before and/or after calling the next one.
type Middleware interface {
	// Middleware wraps a task's process function.
	// It receives the task being executed and the next Process in the chain,
	// and returns a new Process that will be executed instead.
	//
	// The returned Process can:
	//   - Perform setup before calling next.Process(ctx)
	//   - Call next.Process(ctx) to continue the chain
	//   - Perform cleanup after next.Process(ctx) returns
	//   - Short-circuit the chain by returning early without calling next
	Middleware(task *Task, next Process) Process
}

// MiddlewareFunc is a function adapter that allows ordinary functions to be
// used as Middleware implementations.
//
// Example:
//
//	var loggingMiddleware = MiddlewareFunc(func(task *Task, next Process) Process {
//	    return ProcessFunc(func(ctx context.Context) error {
//	        log.Info("starting task", "name", task.Name())
//	        err := next.Process(ctx)
//	        log.Info("finished task", "name", task.Name(), "error", err)
//	        return err
//	    })
//	})
type MiddlewareFunc func(task *Task, next Process) Process

// Middleware implements the Middleware interface by calling the underlying function.
func (m MiddlewareFunc) Middleware(task *Task, next Process) Process {
	return m(task, next)
}
