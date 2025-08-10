package main

import "sync"

type Queue struct {
	mu   sync.Mutex
	waiting  []*Client
}

func NewQueue() *Queue {
	return &Queue{
		waiting: make([]*Client, 0),
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
