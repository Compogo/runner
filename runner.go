package runner

import (
	"io"
)

type Runner interface {
	io.Closer
	RunProcess(Process) error
	RunProcesses(...Process) error
	StopProcess(Process) error
	StopProcessByName(string) error
	HasProcess(Process) bool
	HasProcessByName(string) bool
	Use(middlewares ...Middleware)
}
