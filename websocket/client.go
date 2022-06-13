package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/eric-orenge/message-server/models"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		wsData := models.WsData{}
		err = json.Unmarshal(p, &wsData)
		if err != nil {
			fmt.Println("Error unmarshal")
		}
		message := models.Message{Type: messageType, Body: wsData, ClientID: c.ID}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
