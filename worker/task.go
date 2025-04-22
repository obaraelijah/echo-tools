package worker

// Task This defines how a task used for a pool should be look like.
// The Method Error() returns an error if the function F returns an error
type task struct {
	ret chan error
	F   func() error
}

func NewTask(f func() error) *task {
	return &task{
		F:   f,
		ret: make(chan error, 1), // Buffered channel to fetch the result non-blocking
	}
}

func (t *task) WaitForResult() error {
	return <-t.ret
}

func (t *task) execute() {
	t.ret <- t.F()
}
