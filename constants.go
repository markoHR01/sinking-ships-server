package main

import "time"

const (
	ServerPort       = ":22335"
	QueueInterval    = 1 * time.Second
	SendMessageWait  = 100 * time.Millisecond
	WaitQueueExit    = 10 * time.Millisecond
	WaitPlayerSetup  = 25 * time.Second
	WaitPlayerShips  = 3 * time.Second
	GSPWaitNextRead  = 450 * time.Millisecond
	ShipPartsTotal   = 17
)

var ShipSizes = [5]int{2, 3, 3, 4, 5}
