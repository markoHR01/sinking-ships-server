package main

import "fmt"

func queueWorker(client *Client, queue *Queue) {
	for {
		select {
		case <- client.quitWorker:
			return
		default:
			message, err := client.ReadMessage()
			if err != nil {
				fmt.Println("Queue Worker :", err)

				queue.Leave(client)
				return
			}

			switch message["type"] {
			case "JoinQueue":
				queue.Join(client)
				client.SendMessage(Message{"type": "QueueJoined"})
			case "LeaveQueue":
				queue.Leave(client)
				client.SendMessage(Message{"type": "QueueLeft"})
			default:
				if message != nil {
					fmt.Println("Queue Worker : Unknown Type :", message["type"])
				}
			}
		}
	}
}
