package main

type ShipPart struct {
	X   int
	Y   int
	Hit bool
}

type Ship struct {
	Parts []ShipPart
	Sunk  bool
}
