package gopool

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrPoolNotRuning = errors.New("the pool is not running")
	ErrJobNotFunc = errors.New("generic worker not given a func()")
	ErrWorkerClosed = errors.New("worker was closed")
	ErrJobTimeOut = errors.New("job request time out")
)

type Worker interface {
	//同步执行任务并返回结果
	Process(interface{}) interface{}
	//在worker可以处理下一个任务前，阻塞协程
	BlockUntilReady()
	//任务撤销时调用
	Interrupt()
	//worker从协程池删除
	Terminate()
}

type closureWorker struct {
	processor func(interface{}) interface{}
}

func (w *closureWorker) Process(payload interface{}) interface{} {
	return w.processor(payload)
}

func (w *closureWorker) BlockUntilReady() {}
func (w *closureWorker) Interrupt() {}
func (w *closureWorker) Terminate() {}

type callbackWorker struct {}

func (w *callbackWorker) Process(payload interface{}) interface{} {
	f, ok := payload.(func())
	if !ok {
		return ErrJobNotFunc
	}
	f()
	return nil
}
func (w *callbackWorker) BlockUntilReady() {}
func (w *callbackWorker) Interrupt() {}
func (w *callbackWorker) Terminate() {}

type Pool struct {
	queuedJobs int64
	ctor func() Worker
	workers []*workerWrapper
	reqChan chan workRequest
	workerMut sync.Mutex
}

func New(n int, ctor func() Worker) *Pool {
	p := &Pool{
		ctor:       ctor,
		reqChan:    make(chan workRequest),
	}
	p.SetSize(n)
	return p
}

func (p *Pool) Process(payload interface{}) interface{} {
	atomic.AddInt64(&p.queuedJobs, 1)
	request, open := <-p.reqChan
	if !open {
		panic(ErrPoolNotRuning)
	}
	request.jobChan <- payload
	payload, open = <- request.retChan
	if !open {
		panic(ErrWorkerClosed)
	}
	atomic.AddInt64(&p.queuedJobs, -1)
	return payload
}

func (p *Pool) SetSize(n int) {
	p.workerMut.Lock()
	defer p.workerMut.Unlock()

	lWorkers := len(p.workers)
	if lWorkers == n {
		return
	}
	for i:= lWorkers; i<n; i++ {
		p.workers = append(p.workers, newWorkerWrapper(p.reqChan, p.ctor()))
	}

	//异步停止大于N的workers
	for i:=n; i<lWorkers;i++ {
		p.workers[i].stop()
	}

	//同步等待大于N的workers停止
	for i:=n; i<lWorkers;i++ {
		p.workers[i].join()
	}

	p.workers = p.workers[:n]
}