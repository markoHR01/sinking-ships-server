package main

import (
	"fmt"
	"net"
)

const serverPort = ":22335"

func processClientMessages(client *Client) {
	defer client.Close()

	for {
		message, err := client.ReadMessage()
		if err != nil {
			fmt.Println("client.ReadMessage() failed:", err)
			return
		}

		switch message["type"] {
		case "JoinQueue":
			client.SendMessage(Message{"type": "QueueJoined"})
		case "LeaveQueue":
			client.SendMessage(Message{"type": "QueueLeft"})
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", serverPort)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept() failed:", err)
			continue
		}

		client := NewClient(conn)
		go processClientMessages(client)
	}
}
