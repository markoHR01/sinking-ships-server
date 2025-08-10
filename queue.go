package main

import (
	"sync"
	"time"
)

type Queue struct {
	mu   sync.Mutex
	waiting  []*Client
}

func NewQueue() *Queue {
	return &Queue{
		waiting: make([]*Client, 0),
	}
}

func (q *Queue) Run(matchWorker func(*Client, *Client)) {
	for {
		time.Sleep(1 * time.Second)

		q.mu.Lock()
		for len(q.waiting) >= 2 {
			c1 := q.waiting[0]
			c2 := q.waiting[1]
			q.waiting = q.waiting[2:]

			go matchWorker(c1, c2)
		}
		q.mu.Unlock()
	}
}

func (q *Queue) Join(client *Client) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.waiting = append(q.waiting, client)
}

func (q *Queue) Leave(client *Client) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for i, c := range q.waiting {
		if c == client {
			q.waiting = append(q.waiting[:i], q.waiting[i+1:]...)
			break
		}
	}
}
