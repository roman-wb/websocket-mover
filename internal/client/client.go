package client

import (
	"time"

	"github.com/roman-wb/websocket-mover/internal"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 10 * time.Second
	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

type Logger interface {
	Errorf(format string, args ...interface{})
}

type Broker interface {
	Register(client internal.Client)
	Unregister(client internal.Client)
	ClientMessage(client internal.Client, message []byte)
}

type Client struct {
	logger   Logger
	broker   Broker
	conn     *websocket.Conn
	id       string
	outbound chan []byte
}

func NewClient(logger Logger, broker Broker, conn *websocket.Conn) *Client {
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait)) //nolint:errcheck
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait)) //nolint:errcheck
		return nil
	}) //nolint:errcheck
	return &Client{
		logger:   logger,
		broker:   broker,
		conn:     conn,
		id:       uuid.New().String(),
		outbound: make(chan []byte, 100),
	}
}

func (client *Client) Conn() *websocket.Conn {
	return client.conn
}

func (client *Client) Id() string {
	return client.id
}

func (client *Client) Send(message []byte) {
	client.outbound <- message
}

func (client *Client) Close() {
	client.conn.Close()
	close(client.outbound)
}

func (client *Client) ReadPump() {
	defer func() {
		client.broker.Unregister(client)
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				client.logger.Errorf("Unexpected close %v", err)
			}
			break
		}
		client.broker.ClientMessage(client, message)
	}
}

func (client *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.broker.Unregister(client)
	}()

	for {
		client.conn.SetWriteDeadline(time.Now().Add(writeWait)) //nolint:errcheck
		select {
		case message, ok := <-client.outbound:
			if !ok {
				err := client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					client.logger.Errorf("CloseMessage %v", err)
				}
				return
			}
			err := client.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.logger.Errorf("TextMessage %v", err)
			}
		case <-ticker.C:
			err := client.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				client.logger.Errorf("PingMessage %v", err)
				return
			}
		}
	}
}
