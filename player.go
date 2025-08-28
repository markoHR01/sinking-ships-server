package main

import (
	"time"
	"errors"
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

func GetPlayerShips(m *Match) ([]*Ship, error) {
	messages, err := getShipPositions(m)
	if err != nil { return nil, err }

	ships, err := getShips()
	if err != nil { return nil, err }

	err := shipsAreOK(ships)
	if err != nil { return nil, err }

	return ships, nil
}
