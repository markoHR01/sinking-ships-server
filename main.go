package main

import (
	"fmt"
	"net"
	"strings"
)

const serverPort = ":22335"

func readMessage(data []byte) map[string]string {
	message := make(map[string]string)
	input := string(data)

	msgLines := strings.Split(input, "\n")
	for _, msgLine := range msgLines {
		keyValue := strings.SplitN(msgLine, "=", 2)
		if len(keyValue) == 2 {
			key := keyValue[0]
			value := keyValue[1]
			message[key] = value
		}
	}

	return message
}

func processConnection(conn net.Conn) {
	defer conn.Close()

	var buffer [1024] byte

	for {
		size, err := conn.Read(buffer[:])
		if err != nil {
			fmt.Println("conn.Read(buffer[:]) failed:", err)
			return
		}

		message := readMessage(buffer[:size])

		switch message["type"] {

		case "JoinQueue":
			conn.Write([]byte("type=QueueJoined\nSTOP\n"))
		case "LeaveQueue":
			conn.Write([]byte("type=QueueLeft\nSTOP\n"))
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

		go processConnection(conn)
	}
}
