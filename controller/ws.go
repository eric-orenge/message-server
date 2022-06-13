package controller

import (
	"log"
	"net/http"

	"github.com/eric-orenge/message-server/models"
	"github.com/eric-orenge/message-server/utils"
	"github.com/eric-orenge/message-server/websocket"
	ws "github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

type WsData struct {
	Command string      `json:"command"`
	Data    interface{} `json:"data"`
}

// type WsResponse struct {
// 	Success int         `json:"success"`
// 	Message interface{} `json:"message"`
// }

type WebsocketCtrl struct {
	PendingAck []string
	Messages   []models.Message
}

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (wCtrl *WebsocketCtrl) WsEndpoint(pool *websocket.Pool, c echo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// roomid := c.Param("roomid")
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	client := &websocket.Client{
		ID:   utils.GetID(),
		Conn: ws,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
	return nil
}

// func (wCtrl *WebsocketCtrl) SendMessage(message *models.Message, conn *websocket.Conn) {
// 	wsData := WsData{
// 		Command: "receiveMessage",
// 		Data:    message,
// 	}
// 	msgData, _ := json.Marshal(wsData)
// 	conn.WriteMessage(1, msgData)
// }

// func (wCtrl *WebsocketCtrl) Reader(conn *websocket.Conn) {
// 	for {
// 		// read in a message
// 		_, p, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		wsData := WsData{}
// 		err = json.Unmarshal([]byte(p), &wsData)
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}

// 		switch wsData.Command {
// 		case "newID":
// 			clientID := wCtrl.getID()
// 			wsData := WsData{
// 				Command: "clientID",
// 				Data:    clientID,
// 			}
// 			data, _ := json.Marshal(wsData)
// 			err = conn.WriteMessage(1, data) // send client ID
// 			if err != nil {
// 				log.Println(err)
// 			}
// 		case "sendMessage":
// 			receivedData := wsData.Data.(map[string]interface{})
// 			message := models.Message{
// 				ID:   wCtrl.getID(),
// 				To:   receivedData["to"].(string),
// 				From: receivedData["from"].(string),
// 				Text: receivedData["text"].(string),
// 			}
// 			wCtrl.SendMessage(&message, conn)
// 		}
// 	}
// }
