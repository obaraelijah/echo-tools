package worker

type taskInterface interface {
	Error() error
}

//Task This defines how a task used for a pool should be look like.
// The Method Error() returns an error if the function F returns an error
type Task struct {
	err error
	F   func() error
}

func (t *Task) Error() error {
	return t.err
}

func (t *Task) execute() {
	t.err = t.F()
}
