package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const serverPort = ":22335"

func processConnection(conn net.Conn) {
	defer conn.Close()
	msgBuffer := bufio.NewReader(conn)

	for {
		msg, err := msgBuffer.ReadString('\n')
		if err != nil {
			fmt.Println("msgBuffer.ReadString('\n') failed:", err)
			return
		}

		msg = strings.TrimSpace(msg)
		if strings.HasPrefix(msg, "type=") {
			msgType := strings.TrimPrefix(msg, "type=")

			switch msgType {
			case "JoinQueue":
				conn.Write([]byte("type=QueueJoined\n"))
			case "LeaveQueue":
				conn.Write([]byte("type=QueueLeft\n"))
			default:
				// Unknown - do nothing
			}
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
