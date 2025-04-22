package worker

type pool struct {
	workers   []*worker
	numWorker int

	//newTasks is used to enqueue new tasks while running. The main go routine will append these tasks to queue.
	//newTasks chan bool
	//queue is the channel, worker consume their tasks from
	queue chan *task
	//quit is the control channel to stop the main go routine
	quit chan bool
	//lock sync.Mutex
}

// PoolConfig Configuration for a worker pool.
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
	if c.NumWorker <= 0 || c.QueueSize <= 0 {
		panic("NumWorker and QueueSize must be greater than 0")
	}

	return &pool{
		numWorker: c.NumWorker,
		queue:     make(chan *task, c.QueueSize),
		quit:      make(chan bool),
	}
}

// AddTask adds a task to the queue. Blocking until the Task is enqueued.
func (p *pool) AddTask(t *task) {
	p.queue <- t

}

// AddTasks add a bunch of tasks to the queue. Block until every Task is enqueued.
func (p *pool) AddTasks(tasks []*task) {
	for _, t := range tasks {
		p.queue <- t
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
}

// Stop stops background workers
func (p *pool) Stop() {
	for _, w := range p.workers {
		w.Stop()
	}
	p.quit <- true
}
