package main

import (
	"fmt"
	"net"
	"time"
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
