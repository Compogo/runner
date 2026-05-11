package runner

import (
	"context"
	"io"
)

type Process interface {
	io.Closer
	Process(ctx context.Context) error
	Name() string
}

type ProcessFunc func(ctx context.Context) error
