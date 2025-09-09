package main

import (
	"time"
	"fmt"
)

func matchWorker(match *Match) {
	match.SendMatchFound()

	time.Sleep(WaitPlayerSetup)

	if err := match.Setup(); err != nil {
		fmt.Println("Setup :", err)
		match.Quit()
		return
	}

	match.SendMatchStart()

	if err := match.Play(); err != nil {
		fmt.Println("Match :", err)
		match.Quit()
		return
	}
}
