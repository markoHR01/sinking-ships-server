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
	turn     int
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
	noAttacks := 0

	for !match.IsGameOver() {
		attacker := m.getAttacker()
		defender := m.getDefender()

		message, timeout, err := GetPlayerAttack(attacker)
		if err != nil { return err }

		if timeout || !attackIsOK(message) {
			if noAttacks++; noAttacks > NoAttackLimit {
				match.Quit()
				return nil
			}

			m.SendNoAttack()
			m.nextTurn()
			continue
		}

		m.PlayTurn(message)
		m.nextTurn()
	}

	return nil
}

func (m *Match) PlayTurn(msg Message) {
	x, _ := strconv.Atoi(msg["X"])
	y, _ := strconv.Atoi(msg["Y"])

	hit, sunk, sunkIndex := Attack(m, x, y)

	m.SendAttackResult(x, y, hit, sunk, sunkIndex)
	m.SendEnemyAttack(x, y)
}

func (m *Match) IsGameOver() bool {
	p1Loss := true
	for _, ship := range m.ships1 {
		if !ship.Sunk {
			p1Loss = false
			break
		}
	}

	p2Loss := true
	for _, ship := range m.ships2 {
		if !ship.Sunk {
			p2Loss = false
			break
		}
	}

	return p1Loss || p2Loss
}

func (m *Match) SendMatchFound() {
	msg := Message{"type": "MatchFound"}
	m.player1.SendMessage(msg)
	m.player2.SendMessage(msg)
}

func (m *Match) SendMatchStart() {
	m.turn = rand.Intn(2)

	msg1 := Message{
		"type": "MatchStart",
		"start-first": strconv.FormatBool(m.turn == 0),
	}
	msg2 := Message{
		"type": "MatchStart",
		"start-first": strconv.FormatBool(m.turn == 1),
	}
	m.player1.SendMessage(msg1)
	m.player2.SendMessage(msg2)
}

func (m *Match) SendNoAttack() {
	msg := Message{"type": "NoAttack"}
	m.player1.SendMessage(msg)
	m.player2.SendMessage(msg)
}

func (m *Match) SendAttackResult(
	x, y int,
	hit, sunk bool,
	sunkIndex int,
) {
	msg := Message{
		"type": "AttackResult",
		"X":    strconv.Itoa(x),
		"Y":    strconv.Itoa(y),
		"hit":  strconv.FormatBool(hit),
		"sunk": strconv.FormatBool(sunk),
	}
	if sunk {
		msg["sunk-index"] = strconv.Itoa(sunkIndex)
	}
	m.getAttacker().SendMessage(msg)
}

func (m *Match) SendEnemyAttack(
	x, y int,
) {
	msg := Message{
		"type": "EnemyAttack",
		"X":    strconv.Itoa(x),
		"Y":    strconv.Itoa(y),
	}
	m.getDefender().SendMessage(msg)
}

func (m *Match) getAttacker() *Client {
	if m.turn == 0 {
		return m.player1
	}
	return m.player2
}

func (m *Match) getDefender() *Client {
	if m.turn == 0 {
		return m.player2
	}
	return m.player1
}

func (m *Match) nextTurn() {
	m.turn = 1 - m.turn
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
	go queueWorker(m.player1, m.queue)

	m.player2.MakeQuitChannel()
	go queueWorker(m.player2, m.queue)
}
