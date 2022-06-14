package websocket

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/eric-orenge/message-server/models"
	"github.com/eric-orenge/message-server/redis"
	"github.com/eric-orenge/message-server/utils"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	PendingAck map[string]models.Message
	Broadcast  chan models.Message
	RedisCtrl  *redis.RedisCtrl
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
		if msg, ok := pool.PendingAck[msgID]; ok {
			//add to redis cache
			msgObject, _ := json.Marshal(msg)
			err := pool.RedisCtrl.SetMessage(msgID, msgObject)
			if err != nil {
				log.Println(err)
				return
			}
			// add to client queue
			clientID := msg.Body.Data.(map[string]interface{})["to"].(string)
			err = pool.RedisCtrl.PushMessage(clientID, msgID)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("Message ", msgID, " ack still pending")

		} else {
			//ack has been removed
			log.Println("Message ", msgID, " ack cleared")
		}
	}
}
func (pool *Pool) Start() {
	// connect to redis
	pool.RedisCtrl.Address = os.Getenv("REDIS_URL")
	pool.RedisCtrl.Password = os.Getenv("REDIS_PASSWORD")
	db, _ := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 64)
	pool.RedisCtrl.DB = db

	err := pool.RedisCtrl.Connect()
	if err != nil {
		panic(err)
	}

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
					// send ack to client who sent the message
				}
			} else if message.Body.Command == "newID" {
				log.Println("Client requesting for new ID")
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
			} else if message.Body.Command == "setID" {
				log.Println("Set client ID")
				for client := range pool.Clients {
					if client.ID == message.ClientID {
						//client ID set
						client.ID = message.Body.Data.(string)
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
