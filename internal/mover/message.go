package mover

import (
	"encoding/json"
	"fmt"

	"github.com/roman-wb/websocket-mover/internal"
)

const (
	TypeNotify = "notify"
	TypeState  = "state"

	TypeTryMove   = "tryMove"
	TypeTryLock   = "tryLock"
	TypeTryUnlock = "tryUnlock"
)

const (
	NotifyClientConnected    = "New client connected"
	NotifyClientDisconnected = "Client disconnected"
)

type Message struct {
	Type    string `json:"type"`
	Notify  string `json:"notify,omitempty"`
	Id      string `json:"id,omitempty"`
	OwnerId string `json:"ownerId,omitempty"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
}

func NewMessage(raw []byte) (Message, error) {
	message := Message{}
	return message, json.Unmarshal(raw, &message)
}

func NewPrivateStateMessage(client internal.Client, w *Worker) ([]byte, error) {
	message, err := json.Marshal(Message{
		Type:    TypeState,
		Id:      client.Id(),
		OwnerId: w.ownerId,
		X:       w.x,
		Y:       w.y,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal NewPrivateStateMessage %v", err)
	}
	return message, nil
}

func NewStateMessage(w *Worker) ([]byte, error) {
	message, err := json.Marshal(Message{
		Type:    TypeState,
		OwnerId: w.ownerId,
		X:       w.x,
		Y:       w.y,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal NewStateMessage %v", err)
	}
	return message, nil
}

func NewClientConnectedMessage() ([]byte, error) {
	message, err := json.Marshal(Message{
		Type:   TypeNotify,
		Notify: NotifyClientConnected,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal NewClientConnectedMessage %v", err)
	}
	return message, nil
}

func NewClientDisconnectedMessage() ([]byte, error) {
	message, err := json.Marshal(Message{
		Type:   TypeNotify,
		Notify: NotifyClientDisconnected,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal NewClientDisconnectedMessage %v", err)
	}
	return message, nil
}
