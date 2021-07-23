package mover

import (
	"github.com/roman-wb/websocket-mover/internal"
)

type Logger interface {
	Errorf(format string, args ...interface{})
	Error(i ...interface{})
	Debug(i ...interface{})
}

type Broker interface {
	Clients() map[internal.Client]struct{}
}

type Worker struct {
	logger Logger
	broker Broker

	ownerId string
	x       int
	y       int
}

func NewWorker(logger Logger, broker Broker) *Worker {
	return &Worker{
		logger: logger,
		broker: broker,
	}
}

func (w *Worker) OnConnected(client internal.Client) {
	w.logger.Debug("OnConnected", client.Conn().RemoteAddr())

	w.sendStateTo(client)
	w.broadcastClientConnected(client)
}

func (w *Worker) OnMessage(client internal.Client, raw []byte) {
	w.logger.Debug("OnMessage", client.Conn().RemoteAddr(), string(raw))

	message, err := NewMessage(raw)
	if err != nil {
		w.logger.Error(err)
		return
	}

	switch message.Type {
	case TypeTryMove:
		if w.tryMove(client, message) {
			w.broadcastStateWithout(client)
		}
	case TypeTryLock:
		if w.tryLock(client) {
			w.broadcastState()
		}
	case TypeTryUnlock:
		if w.tryUnlock(client) {
			w.broadcastStateWithout(client)
		}
	default:
		w.logger.Errorf("unknown message type `%v`", message)
	}
}

func (w *Worker) OnClosed(client internal.Client) {
	w.logger.Debug("OnClosed", client.Conn().RemoteAddr())

	if w.tryUnlock(client) {
		w.broadcastState()
	}
	w.broadcastClientDisconnected(client)
}

func (w *Worker) tryLock(client internal.Client) bool {
	if w.ownerId == "" {
		w.ownerId = client.Id()
		return true
	}
	return false
}

func (w *Worker) tryUnlock(client internal.Client) bool {
	if w.ownerId == client.Id() {
		w.ownerId = ""
		return true
	}
	return false
}

func (w *Worker) tryMove(client internal.Client, message Message) bool {
	if w.ownerId == client.Id() && (w.x != message.X || w.y != message.Y) {
		w.x, w.y = message.X, message.Y
		return true
	}
	return false
}

func (w *Worker) sendStateTo(client internal.Client) {
	message, err := NewPrivateStateMessage(client, w)
	if err != nil {
		w.logger.Error(err)
		return
	}
	client.Send(message)
}

func (w *Worker) broadcastStateWithout(current internal.Client) {
	message, err := NewStateMessage(w)
	if err != nil {
		w.logger.Error(err)
		return
	}
	for client := range w.broker.Clients() {
		if client != current {
			client.Send(message)
		}
	}
}

func (w *Worker) broadcastState() {
	message, err := NewStateMessage(w)
	if err != nil {
		w.logger.Error(err)
		return
	}
	for client := range w.broker.Clients() {
		client.Send(message)
	}
}

func (w *Worker) broadcastClientConnected(currentClient internal.Client) {
	message, err := NewClientConnectedMessage()
	if err != nil {
		w.logger.Error(err)
		return
	}
	for client := range w.broker.Clients() {
		if client != currentClient {
			client.Send(message)
		}
	}
}

func (w *Worker) broadcastClientDisconnected(currentClient internal.Client) {
	message, err := NewClientDisconnectedMessage()
	if err != nil {
		w.logger.Error(err)
		return
	}
	for client := range w.broker.Clients() {
		if client != currentClient {
			client.Send(message)
		}
	}
}
