package main

import (
	"time"
	"errors"
	"strconv"
	"fmt"
	"sort"
)

func getShipPositions(m *Match) ([]Message, error) {
	timelimit := time.Now().Add(WaitPlayerShips)
	var messages [2]Message

	players := []*Client{m.player1, m.player2}

	for time.Now().Before(timelimit) {
		for i, p := range players {
			if messages[i] == nil {
				if msg, err := p.ReadMessage(); err != nil {
					return nil, errors.New("Client : Error")
				} else if msg != nil {
					if msg["type"] == "ShipPositions" {
						messages[i] = msg
					}
				}
			}
		}

		if messages[0] != nil && messages[1] != nil {
			return messages[:], nil
		}

		time.Sleep(GSPWaitNextRead)
	}

	return nil, errors.New("Timeout")
}

func shipIsOK(ship [][2]int) bool {
	if len(ship) == 0 {
		return false
	}

	xInLine := true
	yInLine := true
	for _, part := range ship {
		if part[0] != ship[0][0] {
			xInLine = false
		}
		if part[1] != ship[0][1] {
			yInLine = false
		}
	}
	if !(xInLine || yInLine) {
		return false
	}

	if xInLine {
		ys := make([]int, len(ship))
		for i, p := range ship {
			ys[i] = p[1]
		}
		sort.Ints(ys)
		for i := 1; i < len(ys); i++ {
			if ys[i] != ys[i-1] + 1 {
				return false
			}
		}
	} else {
		xs := make([]int, len(ship))
		for i, p := range ship {
			xs[i] = p[0]
		}
		sort.Ints(xs)
		for i := 1; i < len(xs); i++ {
			if xs[i] != xs[i-1] + 1 {
				return false
			}
		}
	}

	return true
}

func getShips(messages []Message) ([][]Ship, error) {
	var ships [][]Ship

	for _, m := range messages {
		var shipXY [][2]int

		idx := 0
		for i := 0; i < ShipPartsTotal; i++ {
			key := strconv.Itoa(idx)
			val, ok := m[key]
			if !ok {
				return nil, errors.New("Message : Corrupted")
			}

			var x, y int
			n, err := fmt.Sscanf(val, "X%dY%d", &x, &y)
			if err != nil || n != 2 {
				return nil, errors.New("Message : Corrupted")
			}

			shipXY = append(shipXY, [2]int{x, y})
			idx++
		}

		seen := make(map[[2]int]bool)
		for _, xy := range shipXY {
			x, y := xy[0], xy[1]

			if x < 0 || x >= 10 || y < 0 || y >= 10 {
				return nil, errors.New("Ship : Out-of-Bounds")
			}

			if seen[xy] {
				return nil, errors.New("Ship : Overlap")
			}
			seen[xy] = true
		}

		start := 0
		for _, size := range ShipSizes {
			end := start + size
			ship := shipXY[start:end]

			if !shipIsOK(ship) {
				return nil, errors.New("Ship : NotOK")
			}

			start = end
		}

		start = 0
		var playerShips []Ship
		for _, size := range ShipSizes {
			end := start + size
			s := shipXY[start:end]

			ship := Ship{}
			for _, xy := range s {
				sp := ShipPart{X: xy[0], Y: xy[1]}
				ship.Parts = append(ship.Parts, sp)
			}
			playerShips = append(playerShips, ship)

			start = end
		}

		ships = append(ships, playerShips)
	}

	return ships, nil
}

func GetPlayerShips(m *Match) ([][]Ship, error) {
	messages, err := getShipPositions(m)
	if err != nil { return nil, err }

	ships, err := getShips(messages)
	if err != nil { return nil, err }

	return ships, nil
}
