package gopool

type workRequest struct {
	//向worker发送任务
	jobChan chan<- interface{}

	//从当前任务中读取结果
	retChan <-chan interface{}

	//用来撤销一个正在运行的任务，当被调用时不必再从retChan读取
	interruptFunc func()
}

//workerWrapper包含一个worker实现，并且封装进一个协程
//管理worker及goroutine的生命周期
type workerWrapper struct {
	worker Worker
	interruptChan chan struct{}

	reqChan chan<- workRequest

	closeChan chan struct{}

	closedChan chan struct{}
}

func newWorkerWrapper(reqChan chan<- workRequest, worker Worker) *workerWrapper {
	w := workerWrapper{
		worker:        worker,
		interruptChan: make(chan struct{}),
		reqChan:       reqChan,
		closeChan:     make(chan struct{}),
		closedChan:    make(chan struct{}),
	}

	go w.run()
	return &w
}

func (w *workerWrapper) interrupt() {
	close(w.interruptChan)
	w.worker.Interrupt()
}

func (w *workerWrapper) run() {
	jobChan, retChan := make(chan interface{}), make(chan interface{})
	defer func() {
		w.worker.Terminate()
		close(retChan)
		close(w.closedChan)
	}()

	for {
		w.worker.BlockUntilReady()
		select {
		case w.reqChan <- workRequest{
			jobChan:       jobChan,
			retChan:       retChan,
			interruptFunc: w.interrupt,
		}:
			select {
			case payload := <-jobChan:
				result := w.worker.Process(payload)
				select {
				case retChan <- result:
				case <-w.interruptChan:
					w.interruptChan = make(chan struct{})
				}
			case _,_ = <-w.interruptChan:
				w.interruptChan = make(chan struct{})
			}
		case <-w.closeChan:
			return
		}
	}
}

func (w *workerWrapper) stop() {
	close(w.closeChan)
}

func (w *workerWrapper) join() {
	<-w.closedChan
}