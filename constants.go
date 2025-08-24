package main

import "time"

const (
	ServerPort       = ":22335"
	QueueInterval    = 1 * time.Second
	SendMessageWait  = 100 * time.Millisecond
)
