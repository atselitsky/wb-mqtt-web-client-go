package main

import "github.com/gorilla/websocket"

type websocketClient struct {
	socket     *websocket.Conn
	send       chan *MQTTmessage
	connection *websocketConnection
}

func (c *websocketClient) write() {
	defer c.socket.Close()
	for msg := range c.send {
		// err := conn.WriteMessage(websocket.TextMessage, msg)
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}

}

// func (c *websocketClient) read() {
// 	defer c.socket.Close()
// 	for {
// 		do nothin
// 	}
// }
