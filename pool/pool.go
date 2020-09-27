package pool

import "fmt"

//定义任务类型
type Task struct {
	f func() error //无参函数
}

func NewTask(f func() error) *Task {
	return &Task{f:f}
}

func (t *Task) Execute() {
	t.f()
}

//协程池
type Pool struct {
	//接收task入口
	EntryChannel chan *Task

	//最大worker数量
	worker_num int

	//内部任务就绪队列
	JobsChannel chan *Task
}

func NewPool(cap int) *Pool {
	return &Pool{
		EntryChannel: make(chan *Task),
		worker_num:   cap,
		JobsChannel:  make(chan *Task),
	}
}

func (p *Pool) worker(workID int) {
	for task := range p.JobsChannel {
		task.Execute()
		fmt.Println("worker ID ", workID, " 执行完毕任务")
	}
}

func (p *Pool) Run() {
	//1.首先根据协程池worker数量限定，开启固定数量worker
	//每个worker用一个Goroutine承载
	for i:=0; i<p.worker_num; i++ {
		go p.worker(i)
	}

	//2.从EntryChannel入口取外界传递过来的任务
	//并将任务放入JobsChannel中
	for task := range p.EntryChannel {
		p.JobsChannel <- task
	}

	//3.执行完毕关闭Channel
	close(p.JobsChannel)
	close(p.EntryChannel)
}