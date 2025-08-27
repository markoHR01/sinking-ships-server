package main

import (
	"time"
	"math/rand"
	"strconv"
)

type Match struct {
	player1  *Client
	player2  *Client
	ships1   []Ship
	ships2   []Ship
}

func NewMatch(c1, c2 *Client) *Match {
	match := &Match {
		player1:  c1,
		player2:  c2,
	}

	return match
}

func CreateMatch(c1, c2 *Client) {
	time.Sleep(WaitQueueExit)
	c1.QuitWorkers()
	c2.QuitWorkers()

	go matchWorker(NewMatch(c1, c2))
}

func (m *Match) Setup() error {
	// ...
}

func (m *Match) SendMatchFound() {
	msg := Message{"type": "MatchFound"}
	m.player1.SendMessage(msg)
	m.player2.SendMessage(msg)
}

func (m *Match) SendMatchStart() {
	startFirst := rand.Intn(2) == 0
	msg1 := Message{
		"type": "MatchStart",
		"start-first": strconv.FormatBool(startFirst),
	}
	msg2 := Message{
		"type": "MatchStart",
		"start-first": strconv.FormatBool(!startFirst),
	}
	m.player1.SendMessage(msg1)
	m.player2.SendMessage(msg2)
}

func (m *Match) Quit() {
	m.SendMatchQuit()
	m.EndMatch()
}

func (m *Match) SendMatchQuit() {
	msg := Message{"type": "MatchQuit"}
	m.player1.SendMessage(msg)
	m.player2.SendMessage(msg)
}

func (m *Match) EndMatch() {
	m.player1.MakeQuitChannel()
	go queueWorker(m.player1)

	m.player2.MakeQuitChannel()
	go queueWorker(m.player2)
}
