package worker

type pool struct {
	workers   []*worker
	numWorker int

	//newTasks is used to enqueue new tasks while running. The main go routine will append these tasks to queue.
	newTasks chan *Task
	//queue is the channel, worker consume their tasks from
	queue chan *Task
	//quit is the control channel to stop the main go routine
	quit chan bool
}

// PoolConfig Configuration for a worker pool. If you want AddTask / AddTasks to be blocking set QueueSize to 0.
type PoolConfig struct {
	NumWorker int
	QueueSize int
}

func NewPool(c *PoolConfig) *pool {
	if c == nil {
		c = &PoolConfig{
			NumWorker: 1,
			QueueSize: 10,
		}
	}
	if c.NumWorker <= 0 {
		panic("NumWorker must be greater than 0")
	}

	return &pool{
		numWorker: c.NumWorker,
		newTasks:  make(chan *Task, c.QueueSize),
		queue:     make(chan *Task),
		quit:      make(chan bool),
	}
}

// AddTask adds a task to the queue
func (p *pool) AddTask(task *Task) {
	p.newTasks <- task
}

// AddTasks add a bunch of tasks to the queue
func (p *pool) AddTasks(tasks []*Task) {
	for _, t := range tasks {
		p.newTasks <- t
	}
}

func (p *pool) Start() {
	for i := 0; i < p.numWorker; i++ {
		w := &worker{
			queue: p.queue,
			quit:  make(chan bool),
		}
		p.workers = append(p.workers, w)
		go w.Start()
	}

	go func() {
		for {
			select {
			case t := <-p.newTasks:
				p.queue <- t
			case <-p.quit:
				return
			}
		}
	}()
}

//Stop stops background workers
func (p *pool) Stop() {
	for _, w := range p.workers {
		w.Stop()
	}
	p.quit <- true
}
