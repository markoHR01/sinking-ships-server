package main

import (
	"fmt"
	"net"
	"time"
	"math/rand"
)

const serverPort = ":22335"

func queueWorker(client *Client, queue *Queue) {
	for {
		select {
		case <- client.quit:
			return
		default:
			message, err := client.ReadMessage()
			if err != nil {
				fmt.Println("client.ReadMessage() failed:", err)
				return
			}

			switch message["type"] {
			case "JoinQueue":
				queue.Join(client)
				client.SendMessage(Message{"type": "QueueJoined"})
			case "LeaveQueue":
				queue.Leave(client)
				client.SendMessage(Message{"type": "QueueLeft"})
			}
		}
	}
}

func matchWorker(client1, client2 *Client) {

	time.Sleep(10 * time.Millisecond)

	close(client1.quit)
	close(client2.quit)

	client1.SendMessage(Message{"type": "MatchFound"})
	client2.SendMessage(Message{"type": "MatchFound"})

	time.Sleep(25 * time.Second)

	var m1, m2 Message
	var err1, err2 error
	for i := 0; i < 6; i++ {
		if m1 == Message{} {
			m1, err1 = client1.ReadMessage()
			if err1 != nil {
				// can't continue
			}
		}
		if m2 == Message{} {
			m2, err2 = client2.ReadMessage()
			if err2 != nil {
				// can't continue
			}
		}

		if m1 != Message{} && m2 != Message{} {
			break
		} else {
			time.Sleep(500 * time.Millisecond)
		}
	}

	if m1 == Message{} || !f(m1) {
		// can't continue
	}

	if m2 == Message{} || !f(m2) {
		// can't continue
	}

	// Extract data from the message

	startFirst := rand.Intn(2) == 0

	client1.SendMessage(Message{"type": "MatchStart",
                                "start-first": startFirst})
	client2.SendMessage(Message{"type": "MatchStart",
                                "start-first": !startFirst})
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	listener, err := net.Listen("tcp", serverPort)
	if err != nil {
		panic(err)
	}

	queue := NewQueue()
	go queue.Run(matchWorker)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept() failed:", err)
			continue
		}

		client := NewClient(conn)
		go queueWorker(client, queue)
	}
}
