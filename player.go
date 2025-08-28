package main

func GetPlayerShips(m *Match) ([]*Ship, error) {
	messages, err := getShipPositions(m)
	if err != nil { return nil, err }

	ships, err := getShips()
	if err != nil { return nil, err }

	err := shipsAreOK(ships)
	if err != nil { return nil, err }

	return ships, nil
}
