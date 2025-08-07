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

func (q *Queue) Join(client *Client) {
	q.mu.Lock()
	q.waiting = append(q.waiting, client)

	if len(q.waiting) >= 2 {
		client1 := q.waiting[0]
		client2 := q.waiting[1]
		q.waiting = q.waiting[2:]
		q.mu.Unlock()

		go startMatch(client1, client2)
	} else {
		q.mu.Unlock()
	}
}

func startMatch(c1, c2 *Client) {
	time.Sleep(10 * time.Millisecond)

	c1.SendMessage(Message{"type": "MatchFound"})
	c2.SendMessage(Message{"type": "MatchFound"})
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
