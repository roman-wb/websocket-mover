package broker

import "github.com/roman-wb/websocket-mover/internal"

type Listner interface {
	OnConnected(conn internal.Client)
	OnMessage(conn internal.Client, message []byte)
	OnClosed(conn internal.Client)
}

type Broker struct {
	clients map[internal.Client]struct{}

	register   chan internal.Client
	unregister chan internal.Client

	listner Listner
}

func NewBroker() *Broker {
	return &Broker{
		clients: make(map[internal.Client]struct{}),

		register:   make(chan internal.Client),
		unregister: make(chan internal.Client),
	}
}

func (b *Broker) Run(listner Listner) {
	b.listner = listner

	for {
		select {
		case client := <-b.register:
			b.clients[client] = struct{}{}
			go client.ReadPump()
			go client.WritePump()
			b.listner.OnConnected(client)
		case client := <-b.unregister:
			if _, ok := b.clients[client]; ok {
				delete(b.clients, client)
				client.Close()
				b.listner.OnClosed(client)
			}
		}
	}
}

func (b *Broker) Register(client internal.Client) {
	b.register <- client
}

func (b *Broker) Unregister(client internal.Client) {
	b.unregister <- client
}

func (b *Broker) Clients() map[internal.Client]struct{} {
	return b.clients
}

func (b *Broker) ClientMessage(client internal.Client, message []byte) {
	b.listner.OnMessage(client, message)
}
