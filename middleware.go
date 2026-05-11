package runner

type Middleware interface {
	Middleware(Process, ProcessFunc) ProcessFunc
}

type MiddlewareFunc func(Process, ProcessFunc) ProcessFunc

func (m MiddlewareFunc) Middleware(p Process, next ProcessFunc) ProcessFunc {
	return m(p, next)
}
