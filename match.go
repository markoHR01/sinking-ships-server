package main

import (
	"time"
	"math/rand"
	"strconv"
)

type Match struct {
	queue    *Queue
	player1  *Client
	player2  *Client
	ships1   []Ship
	ships2   []Ship
}

func NewMatch(q *Queue, c1, c2 *Client) *Match {
	match := &Match {
		queue:    q,
		player1:  c1,
		player2:  c2,
	}

	return match
}

func CreateMatch(q *Queue, c1, c2 *Client) {
	time.Sleep(WaitQueueExit)
	c1.QuitWorkers()
	c2.QuitWorkers()

	go matchWorker(NewMatch(q, c1, c2))
}

func (m *Match) Setup() error {
	ships, err := GetPlayerShips(m)
	if err != nil { return err }

	m.ships1 = ships[0]
	m.ships2 = ships[1]
	return nil
}

func (m *Match) Play error {

	for !match.IsGameOver() {
		player := match.getAttacker()

		message, timeout, err := GetPlayerAttack(player)
		if err != nil { return err }

		if timeout {
			// Send No Attack, Switch Turns
			// No Attack > 3 ? Match Quit
			// continue
		}

		// Send Attack Result, Enemy Attack
		// Switch Turns
		// continue
	}

	return nil
}

func (m *Match) IsGameOver() bool {} // Check if player won/lost

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

func (m *Match) getAttacker() *Client {} // Get player (turn == true)

func (m *Match) getDefender() *Client {} // Get player (turn == false)

func (m *Match) nextTurn() {} // Switch players

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
	go queueWorker(m.player1, m.queue)

	m.player2.MakeQuitChannel()
	go queueWorker(m.player2, m.queue)
}
