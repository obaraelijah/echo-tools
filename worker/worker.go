package worker

type Worker interface {
	SetQueue(chan Task)
	Start()
	Stop()
}

type worker struct {
	queue chan Task
	quit  chan bool
}

func (w *worker) SetQueue(c chan Task) {
	w.queue = c
}

func (w *worker) Start() {
	for {
		select {
		case <-w.quit:
			return
		case t := <-w.queue:
			t.Execute()
		}
	}
}

// Stop is used to stop a running worker
func (w *worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
