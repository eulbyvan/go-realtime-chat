package model

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Client struct {
	Id     string
	Socket *websocket.Conn
	Send   chan []byte
}

func (c *Client) Read(manager *ClientManager) {
	defer func() {
		manager.Unregister <- c
		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			manager.Unregister <- c
			c.Socket.Close()
			break
		}

		jsonMessage, _ := json.Marshal(&Message{Sender: c.Id, Content: string(message)})
		manager.Broadcast <- jsonMessage
	}
}

func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for message := range c.Send {
		if err := c.Socket.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}

	// Send a close message if the loop exits
	c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
}
