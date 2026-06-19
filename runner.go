package runner

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/Compogo/compogo"
	typesErrors "github.com/Compogo/types/errors"
	"github.com/Compogo/types/linker"
	"github.com/Compogo/types/set"
)

// Runner управляет жизненным циклом процессов (Task).
// Предоставляет возможности:
//   - Запуска процессов с поддержкой middleware
//   - Остановки процессов по имени или объекту
//   - Graceful shutdown через io.Closer
//   - Проверки наличия процессов
//
// Runner потокобезопасен и может использоваться из нескольких горутин.
type Runner interface {
	io.Closer
	// RunProcess запускает один процесс.
	RunProcess(Process) error

	// RunProcesses запускает несколько процессов.
	RunProcesses(...Process) error

	// StopProcess останавливает процесс по объекту.
	StopProcess(Process) error

	// StopProcessByName останавливает процесс по имени.
	StopProcessByName(string) error

	// HasProcess проверяет наличие процесса по объекту.
	HasProcess(Process) bool

	// HasProcessByName проверяет наличие процесса по имени.
	HasProcessByName(string) bool

	// Use добавляет middleware для всех процессов.
	Use(middlewares ...Middleware)
}

// runner — внутренняя реализация Runner.
type runner struct {
	wg      sync.WaitGroup
	rwMutex sync.RWMutex

	processes     set.Set[Process]
	linkProcesses *linker.Linker[string, Process]
	middlewares   []Middleware

	closer compogo.Closer
	logger compogo.Logger
}

// newRunner создаёт новый Runner.
func newRunner(closer compogo.Closer, logger compogo.Logger) *runner {
	return &runner{
		processes:     set.NewSet[Process](),
		linkProcesses: linker.NewLinker[string, Process](linker.KeyStringNormalizer[Process]()),
		closer:        closer,
		logger:        logger.GetLogger("runner"),
	}
}

func (runner *runner) Use(middlewares ...Middleware) {
	runner.rwMutex.Lock()
	defer runner.rwMutex.Unlock()
	runner.middlewares = append(runner.middlewares, middlewares...)
}

func (runner *runner) getMiddlewares() []Middleware {
	runner.rwMutex.RLock()
	defer runner.rwMutex.RUnlock()
	return runner.middlewares
}

func (runner *runner) Close() (err error) {
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

func (runner *runner) RunProcesses(process ...Process) (err error) {
	for _, proc := range process {
		if err = runner.RunProcess(proc); err != nil {
			return err
		}
	}

	return nil
}

func (runner *runner) RunProcess(process Process) (err error) {
	if runner.HasProcess(process) {
		return fmt.Errorf("[runner] process '%s': %w", process.Name(), TaskAlreadyExistsError)
	}

	middlewares := runner.getMiddlewares()
	processFunc := process.Process
	for i := len(middlewares) - 1; i >= 0; i-- {
		processFunc = middlewares[i].Middleware(process, processFunc)
	}

	runner.addProcess(process)

	runner.wg.Go(func(processFunc ProcessFunc) func() {
		return func() {
			defer runner.removeProcess(process)
			if err := processFunc(runner.closer.GetContext()); err != nil {
				runner.logger.Errorf("process '%s' executed failed: %s", process.Name(), err.Error())
			}
		}
	}(processFunc))

	return nil
}

func (runner *runner) StopProcessByName(name string) error {
	runner.rwMutex.Lock()
	process, err := runner.linkProcesses.Get(name)
	runner.rwMutex.Unlock()

	if errors.Is(err, typesErrors.DoesNotExistError) {
		return fmt.Errorf("[runner] task '%s': %w", name, TaskUndefinedError)
	}

	if err != nil {
		return err
	}

	return runner.StopProcess(process)
}

func (runner *runner) StopProcess(process Process) (err error) {
	if err := process.Close(); err != nil {
		return err
	}

	runner.removeProcess(process)

	return nil
}

func (runner *runner) addProcess(process Process) {
	runner.rwMutex.Lock()
	defer runner.rwMutex.Unlock()

	runner.processes.Add(process)
	runner.linkProcesses.Add(process.Name(), process)
}

func (runner *runner) removeProcess(process Process) {
	runner.rwMutex.Lock()
	defer runner.rwMutex.Unlock()

	runner.processes.Remove(process)
	runner.linkProcesses.Remove(process.Name())
}

func (runner *runner) HasProcess(process Process) bool {
	runner.rwMutex.RLock()
	defer runner.rwMutex.RUnlock()

	return runner.processes.Contains(process)
}

func (runner *runner) HasProcessByName(name string) bool {
	runner.rwMutex.RLock()
	defer runner.rwMutex.RUnlock()

	return runner.linkProcesses.Has(name)
}
