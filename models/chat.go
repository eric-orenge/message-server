package models

type Message struct {
	ID       string `json:"id"`
	Type     int    `json:"type"`
	Body     WsData `json:"body"`
	ClientID string `json:"clientID"`
}

type WsData struct {
	Command string      `json:"command"`
	Data    interface{} `json:"data"`
}
