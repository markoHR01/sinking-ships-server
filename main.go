package main

import (
	"math/rand"
	"time"
	"net"
	"fmt"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	listener, err := net.Listen("tcp", ServerPort)
	if err != nil {
		panic(err)
	}

	queue := NewQueue()
	go queue.Run()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Main :", err)
			continue
		}

		client := NewClient(conn)
		go queueWorker(client, queue)
	}
}
