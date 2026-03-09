package runner

import "context"

// Process defines the interface for any executable unit that can be run
// as a task. The Process method receives a context that is canceled when
// the task should stop (e.g., during application shutdown).
//
// Example:
//
//	type Server struct {}
//
//	func (s *Server) Process(ctx context.Context) error {
//	    return s.httpServer.ListenAndServe()
//	}
type Process interface {
	// Process executes the task logic. The provided context is canceled
	// when the task needs to stop (graceful shutdown). The method should
	// block until the task completes or the context is done.
	Process(ctx context.Context) error
}

// ProcessFunc is a function adapter that allows ordinary functions to be
// used as Process implementations.
//
// Example:
//
//	task := runner.NewTask("worker", runner.ProcessFunc(func(ctx context.Context) error {
//	    for {
//	        select {
//	        case <-ctx.Done():
//	            return nil
//	        default:
//	            doWork()
//	        }
//	    }
//	}))
type ProcessFunc func(ctx context.Context) error

// Process implements the Process interface by calling the underlying function.
func (p ProcessFunc) Process(ctx context.Context) error {
	return p(ctx)
}
