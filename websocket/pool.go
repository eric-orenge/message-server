package websocket

import (
	"context"
	"log"
	"time"

	"github.com/eric-orenge/message-server/models"
	"github.com/eric-orenge/message-server/utils"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	PendingAck map[string]models.Message
	Broadcast  chan models.Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		PendingAck: make(map[string]models.Message),
		Broadcast:  make(chan models.Message),
	}
}

func (pool *Pool) handleAckReceipts(msgID string) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second) // Cancel in 5 seconds
	defer cancel()

	select {
	case <-ctx.Done(): // When time is out
		log.Println("Periodic sync done: ", msgID)
		// message has not been received after 30 seconds
		// cache message in redis
		if _, ok := pool.PendingAck[msgID]; ok {
			log.Println("Message ", msgID, " ack still pending")
			//cache

		} else {
			//ack has been removed
			log.Println("Message ", msgID, " ack cleared")
			ctx.Done()
		}
	}
}
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			wsData := models.WsData{Command: "newUser"}
			for client, _ := range pool.Clients {
				log.Println(client)
				client.Conn.WriteJSON(models.Message{Type: 1, Body: wsData})
			}
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			log.Println("Size of Connection Pool: ", len(pool.Clients))
			wsData := models.WsData{Command: "userDisconnected"}
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(models.Message{Type: 1, Body: wsData})
			}
		case message := <-pool.Broadcast:
			log.Println("Sending message to all clients in Pool")
			if message.Body.Command == "ack" { //ack message was received.
				msgID := message.Body.Data.(string)
				// remove from pending acks since it has been received
				if msg, ok := pool.PendingAck[msgID]; ok {
					delete(pool.PendingAck, msg.ID)
				}
			} else if message.Body.Command == "newID" {
				for client := range pool.Clients {
					if client.ID == message.ClientID {
						message.Body.Data = client.ID
						message.Body.Command = "clientID"
						if err := client.Conn.WriteJSON(message); err != nil {
							log.Println(err)
							return
						}
					}
				}
			} else { // do not broadcast ack receipts
				log.Println("Size of Connection Pool: ", len(pool.Clients))
				for client := range pool.Clients {
					message.ID = utils.GetID()
					if err := client.Conn.WriteJSON(message); err != nil {
						log.Println(err)
						return
					}
					// pending ack
					pool.PendingAck[message.ID] = message
					//wait for ack for message sent
					go pool.handleAckReceipts(message.ID)
					// if not received, add message to the queue

				}
			}

		}
	}
}
