package main

import "time"

func matchWorker(match *Match) {
	match.SendMatchFound()

	time.Sleep(WaitPlayerSetup)

	// Receive
	// Parse
	// Validate

	match.SendMatchStart()
}
