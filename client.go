package main

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"errors"
	"time"
)

type Client struct {
	conn    net.Conn
	in      chan Message
	out     chan Message
	alive   atomic.Bool
	kill    sync.Once
	quit    chan struct{}
}

func NewClient(conn net.Conn) *Client {
	client := &Client {
		conn:  conn,
		in:    make(chan Message, 3),
		out:   make(chan Message, 3),
		quit:  make(chan struct{}),
	}
	client.alive.Store(true)

	go client.receive()
	go client.send()

	return client
}

func (c *Client) ReadMessage() (Message, error) {
	if !c.alive.Load() {
		err := errors.New("Connection is closed")
		return Message{}, err
	}

	select {
	case message := <-c.in:
		return message, nil
	default:
		return Message{}, nil
	}
}

func (c *Client) SendMessage(message Message) (bool, error) {
	if !c.alive.Load() {
		err := errors.New("Connection is closed")
		return false, err
	}

	select {
	case c.out <- message:
		return true, nil
	default:
		return false, nil
	}
}

func (c *Client) receive() {

	var buffer [1024] byte

	for {
		size, err := c.conn.Read(buffer[:])
		if err != nil {
			fmt.Println("c.conn.Read(buffer[:]) failed", err)
			c.Close()
			return
		}

		message := Deserialize(buffer[:size])

		select {
		case c.in <- message:
		default:
		}
	}
}

func (c *Client) send() {

	for {
		var message Message
		select {
		case message = <-c.out:
		default:
			if !c.alive.Load() {
				return
			}
			time.Sleep(100 * time.Millisecond)
			continue
		}

		payload := []byte(Serialize(message))

		_, err := c.conn.Write(payload)
		if err != nil {
			fmt.Println("c.conn.Write(payload) failed:", err)
			c.Close()
			return
		}
	}
}

func (c *Client) Close() {
	c.kill.Do(func() {
		c.alive.Store(false)
		c.conn.Close()
	})
}
