package worker

type Pool interface {
	AddTask(t Task)
	AddTasks(t []Task)
	Start()
	StartWithWorkerCreator(f func() (Worker, error)) error
	Stop()
}

type pool struct {
	workers   []Worker
	numWorker int

	//newTasks is used to enqueue new tasks while running. The main go routine will append these tasks to queue.
	//newTasks chan bool
	//queue is the channel, worker consume their tasks from
	queue chan Task
}

// PoolConfig Configuration for a worker pool.
type PoolConfig struct {
	NumWorker int
	QueueSize int
}

// NewPool is used to create a new pool instance
func NewPool(c *PoolConfig) Pool {
	if c == nil {
		c = &PoolConfig{
			NumWorker: 1,
			QueueSize: 10,
		}
	}
	if c.NumWorker <= 0 || c.QueueSize <= 0 {
		panic("NumWorker and QueueSize must be greater than 0")
	}

	return &pool{
		numWorker: c.NumWorker,
		queue:     make(chan Task, c.QueueSize),
	}
}

// AddTask adds a task to the queue. Blocking until the Task is enqueued.
func (p *pool) AddTask(t Task) {
	p.queue <- t
}

// AddTasks add a bunch of tasks to the queue. Block until every Task is enqueued.
func (p *pool) AddTasks(tasks []Task) {
	for _, t := range tasks {
		p.queue <- t
	}
}

func (p *pool) Start() {
	for i := 0; i < p.numWorker; i++ {
		w := &worker{
			quit: make(chan bool),
		}
		w.SetQueue(p.queue)
		p.workers = append(p.workers, w)
		go w.Start()
	}
}
func (p *pool) StartWithWorkerCreator(f func() (Worker, error)) error {
	for i := 0; i < p.numWorker; i++ {
		w, err := f()
		if err != nil {
			return err
		}
		w.SetQueue(p.queue)
		p.workers = append(p.workers, w)
		go w.Start()
	}
	return nil
}

// Stop stops background workers
func (p *pool) Stop() {
	for _, w := range p.workers {
		w.Stop()
	}
	p.workers = make([]Worker, 0)
}
