package broadcast

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type BroadcastClient struct {
	conn *websocket.Conn
	send chan any
}

func CreateBroadcastClient(conn *websocket.Conn) *BroadcastClient {
	return &BroadcastClient{
		conn: conn,
		send: make(chan any, 32),
	}
}

func (c *BroadcastClient) WritePump() {
	defer func() {
		c.conn.Close()
	}()

	for msg := range c.send {
		if err := c.conn.WriteJSON(msg); err != nil {
			// If write fails, bail out – connection is probably dead
			fmt.Println("write error:", err)
			return
		}
	}
}

func (c *BroadcastClient) Send(msg any) {
	select {
	case c.send <- msg:
	default:
		fmt.Println("send buffer full, dropping message")
	}
}
