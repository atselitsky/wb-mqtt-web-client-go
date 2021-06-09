package websocketConn

import (
	"log"
	"net/http"

	"github.com/atselitsky/wb-mqtt-web-client-go/pkg/configs"
	"github.com/gorilla/websocket"
)

type WebsocketConnection struct {
	Forward chan *configs.MQTTmessage
	join    chan *websocketClient
	leave   chan *websocketClient
	clients map[*websocketClient]bool
}

func NewConn() *WebsocketConnection {
	return &WebsocketConnection{
		Forward: make(chan *configs.MQTTmessage),
		join:    make(chan *websocketClient),
		leave:   make(chan *websocketClient),
		clients: make(map[*websocketClient]bool),
	}
}

func (wc *WebsocketConnection) Run() {
	for {
		select {
		case client := <-wc.join:
			wc.clients[client] = true
		case client := <-wc.leave:
			delete(wc.clients, client)
			close(client.send)
		case msg := <-wc.Forward:
			for client := range wc.clients {
				client.send <- msg
			}
		}
	}
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (wc *WebsocketConnection) WsConn(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error while upgrade to webSockets:", err)
		return
	}

	client := &websocketClient{
		socket:     socket,
		send:       make(chan *configs.MQTTmessage, 1024),
		connection: wc,
	}
	wc.join <- client
	defer func() {
		wc.leave <- client
	}()

	client.write()

}
