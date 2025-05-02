package worker

import "context"

type Task interface {
	WaitForResult() error
	Execute()
	ExecuteWithContext(ctx context.Context)
}

// task This defines how a task used for a pool should be look like.
// The Method Error() returns an error if the function F returns an error
type task struct {
	ret      chan error
	F        func() error
	FWithCTX func(ctx context.Context) error
}

func NewTask(f func() error) Task {
	return &task{
		F:   f,
		ret: make(chan error, 1), // Buffered channel to fetch the result non-blocking
	}
}

func NewTaskWithContext(f func(ctx context.Context) error) Task {
	return &task{
		FWithCTX: f,
		ret:      make(chan error, 1),
	}
}

func (t *task) WaitForResult() error {
	return <-t.ret
}

func (t *task) Execute() {
	t.ret <- t.F()
}

func (t *task) ExecuteWithContext(ctx context.Context) {
	t.ret <- t.FWithCTX(ctx)
}
