package runner

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Compogo/compogo/closer"
	"github.com/Compogo/compogo/logger"
	"github.com/Compogo/types/linker"
	"github.com/Compogo/types/set"
)

type Default struct {
	wg      sync.WaitGroup
	rwMutex sync.RWMutex

	processes     set.Set[Process]
	linkProcesses *linker.Linker[string, Process]
	middlewares   []Middleware

	closer closer.Closer
	logger logger.Logger
}

func NewDefault(closer closer.Closer, logger logger.Logger) *Default {
	return &Default{
		processes:     set.NewSet[Process](),
		linkProcesses: linker.NewLinker[string, Process](linker.KeyStringNormalizer[Process]()),
		closer:        closer,
		logger:        logger.GetLogger("runner"),
	}
}

func (runner *Default) Use(middlewares ...Middleware) {
	runner.rwMutex.Lock()
	defer runner.rwMutex.Unlock()
	runner.middlewares = append(runner.middlewares, middlewares...)
}

func (runner *Default) getMiddlewares() []Middleware {
	runner.rwMutex.RLock()
	defer runner.rwMutex.RUnlock()
	return runner.middlewares
}

func (runner *Default) Close() (err error) {
	runner.rwMutex.Lock()
	processes := runner.processes.Clone()
	runner.rwMutex.Unlock()

	for process := range processes {
		if err = runner.StopProcess(process); err != nil {
			return err
		}
	}

	runner.wg.Wait()

	return nil
}

func (runner *Default) RunProcesses(process ...Process) (err error) {
	for _, proc := range process {
		if err = runner.RunProcess(proc); err != nil {
			return err
		}
	}

	return nil
}

func (runner *Default) RunProcess(process Process) (err error) {
	if runner.HasProcess(process) {
		return fmt.Errorf("[runner] process '%s': %w", process.Name(), TaskAlreadyExistsError)
	}

	middlewares := runner.getMiddlewares()
	processFunc := process.Process
	for i := len(middlewares) - 1; i >= 0; i-- {
		processFunc = middlewares[i].Middleware(process, processFunc)
	}

	runner.addProcess(process)

	runner.wg.Go(func() {
		defer func() {
			if err := runner.StopProcess(process); err != nil && !errors.Is(err, TaskUndefinedError) {
				runner.logger.Errorf("process '%s' stop failed: %w", process.Name(), err)
			}
		}()

		if err := processFunc(runner.closer.GetContext()); err != nil {
			runner.logger.Errorf("process '%s' executed failed: %s", process.Name(), err.Error())
		}
	})

	return nil
}

func (runner *Default) StopProcessByName(name string) error {
	if !runner.HasProcessByName(name) {
		return fmt.Errorf("[runner] task '%s': %w", name, TaskUndefinedError)
	}

	runner.rwMutex.Lock()
	process, _ := runner.linkProcesses.Get(name)
	runner.rwMutex.Unlock()

	return runner.StopProcess(process)
}

func (runner *Default) StopProcess(process Process) (err error) {
	if err := process.Close(); err != nil {
		return err
	}

	runner.removeProcess(process)

	return nil
}

func (runner *Default) addProcess(process Process) {
	runner.rwMutex.Lock()
	defer runner.rwMutex.Unlock()

	runner.processes.Add(process)
	runner.linkProcesses.Add(process.Name(), process)
}

func (runner *Default) removeProcess(process Process) {
	runner.rwMutex.Lock()
	defer runner.rwMutex.Unlock()

	runner.processes.Remove(process)
	runner.linkProcesses.Remove(process.Name())
}

func (runner *Default) HasProcess(process Process) bool {
	runner.rwMutex.RLock()
	defer runner.rwMutex.RUnlock()

	return runner.processes.Contains(process)
}

func (runner *Default) HasProcessByName(name string) bool {
	runner.rwMutex.RLock()
	defer runner.rwMutex.RUnlock()

	return runner.linkProcesses.Has(name)
}
