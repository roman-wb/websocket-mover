package internal

import "github.com/gorilla/websocket"

type Client interface {
	Conn() *websocket.Conn
	Id() string
	Send(message []byte)
	Close()
	ReadPump()
	WritePump()
}
