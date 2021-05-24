package main

import (
	"github.com/gin-gonic/gin"
)

func execute() {
	websocketconn := newConn()
	go websocketconn.run()

	client := NewMQTTConn(websocketconn)
	go client.StartMQTTConnection()

	r := gin.Default()
	r.LoadHTMLFiles("./assets/index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/ws", func(c *gin.Context) {
		websocketconn.wsConn(c.Writer, c.Request)
	})

	r.Run()
}
