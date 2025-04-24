package domain

import "github.com/gorilla/websocket"

type WSClient struct {
	Conn *websocket.Conn
	Send chan []byte
}
