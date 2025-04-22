package worker

type worker struct {
	queue chan *Task
	quit  chan bool
}

func (w *worker) Start() {
	for {
		select {
		case <-w.quit:
			return
		case t := <-w.queue:
			t.execute()
		}
	}
}

// Stop is used to stop a running worker
func (w *worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
