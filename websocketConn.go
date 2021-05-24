package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type websocketConnection struct {
	forward chan *MQTTmessage
	join    chan *websocketClient
	leave   chan *websocketClient
	clients map[*websocketClient]bool
}

func newConn() *websocketConnection {
	return &websocketConnection{
		forward: make(chan *MQTTmessage),
		join:    make(chan *websocketClient),
		leave:   make(chan *websocketClient),
		clients: make(map[*websocketClient]bool),
	}
}

func (wc *websocketConnection) run() {
	for {
		select {
		case client := <-wc.join:
			wc.clients[client] = true
		case client := <-wc.leave:
			delete(wc.clients, client)
			close(client.send)
		case msg := <-wc.forward:
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

func (wc *websocketConnection) wsConn(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error while upgrade to webSockets:", err)
		return
	}

	client := &websocketClient{
		socket:     socket,
		send:       make(chan *MQTTmessage, 1024),
		connection: wc,
	}
	wc.join <- client
	defer func() {
		wc.leave <- client
	}()

	client.write()

}
