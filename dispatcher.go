package main

import (
	"container/heap"
	"fmt"
)

type Pool []*Worker

type Balancer struct {
	Pool *Pool
	Done chan *Worker
}

func Dispatch(jobRequests <-chan *Job, done chan *Worker) {
	var p Pool
	heap.Init(&p)

	b := &Balancer {
		Pool: &p,
		Done: done,
	}

	b.Balance(jobRequests)
}

func (b *Balancer) Balance(jobRequests <-chan *Job) {
	for {
		select {
		case job := <-jobRequests:
			b.dispatch(job)
			fmt.Println(b.Pool)
		case worker := <-b.Done:
			b.complete(worker)
		}
	}
}

func (b *Balancer) dispatch(job *Job) {
	w := heap.Pop(b.Pool).(*Worker)
	w.jobChan <- job
	w.pending += 1
	heap.Push(b.Pool, w)
}

func (b *Balancer) complete(worker *Worker) {
	worker.pending -= 1
	heap.Fix(b.Pool, worker.index)
}

func (p Pool) Len () int { return len(p) }

func (p Pool) Less(i, j int) bool { return p[i].pending < p[j].pending }

func (p Pool) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].index = i
	p[j].index = j
}

func (p *Pool) Push(w interface{}) {
	worker := w.(*Worker)
	worker.index = p.Len()
	*p = append(*p,  worker)
}

func (p *Pool) Pop() interface{} {
	old := *(p)
	n := len(old)
	item := old[n-1]
	item.index = -1
	*(p) = old[0 : n-1]
	return item
}
