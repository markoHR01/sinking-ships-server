package main

import (
	"net"
	"sync"
	"sync/atomic"
	"fmt"
	"time"
	"errors"
)

type Client struct {
	conn        net.Conn
	in          chan Message
	out         chan Message
	alive       atomic.Bool
	kill        sync.Once
	quitWorker  chan struct{}
}

func NewClient(conn net.Conn) *Client {
	client := &Client {
		conn:        conn,
		in:          make(chan Message, 3),
		out:         make(chan Message, 3),
		quitWorker:  make(chan struct{}),
	}
	client.alive.Store(true)

	go client.receive()
	go client.send()

	return client
}

func (c *Client) QuitWorkers() {
	close(c.quitWorker)
}

func (c *Client) MakeQuitChannel() {
	c.quitWorker = make(chan struct{})
}

func (c *Client) Close() {
	c.kill.Do(func() {
		c.alive.Store(false)
		c.conn.Close()
	})
}

func (c *Client) ReadMessage() (Message, error) {
	if !c.alive.Load() {
		err := errors.New("Read Message : Error : Connection Closed")
		return nil, err
	}

	select {
	case m := <- c.in:
		return m, nil
	default:
		return nil, nil
	}
}

func (c *Client) SendMessage(m Message) (bool, error) {
	if !c.alive.Load() {
		err := errors.New("Send Message : Error : Connection Closed")
		return false, err
	}

	select {
	case c.out <- m:
		return true, nil
	default:
		fmt.Println("Output Buffer Overload : Message Dropped")

		return false, nil
	}
}

func (c *Client) receive() {
	var buffer [1024] byte

	for {
		size, err := c.conn.Read(buffer[:])
		if err != nil {
			fmt.Println("Client : Read Error : ", err)

			c.Close()
			return
		}

		message := Deserialize(buffer[:size])

		select {
		case c.in <- message:
		default:
			fmt.Println("Input Buffer Overload : Message Dropped")
		}
	}
}

func (c *Client) send() {

	for {
		select {
		case message := <- c.out:
			payload := []byte(Serialize(message))

			if _, err := c.conn.Write(payload); err != nil {
				fmt.Println("Client : Send Error : ", err)

				c.Close()
				return
			}
		default:
			if !c.alive.Load() {
				return
			}

			time.Sleep(SendMessageWait)
		}
	}
}
