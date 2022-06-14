package api

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

var (
	PendingAck []string
	Messages   []models.Message
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WsEndpoint(pool *websocket.Pool, c echo.Context) error {
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
