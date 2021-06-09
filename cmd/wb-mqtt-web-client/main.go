package main

import (
	"net/http"

	"github.com/atselitsky/wb-mqtt-web-client-go/pkg/devices"
	"github.com/atselitsky/wb-mqtt-web-client-go/pkg/mqttConn"
	"github.com/atselitsky/wb-mqtt-web-client-go/pkg/rules"
	"github.com/atselitsky/wb-mqtt-web-client-go/pkg/websocketConn"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	// TO allow CORS
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}
func main() {
	rules.GetRulesFiles()
	websocketconn := websocketConn.NewConn()
	go websocketconn.Run()

	client := mqttConn.NewMQTTConn(websocketconn)
	go client.StartMQTTConnection()

	r := gin.Default()
	// r.Use(CORS())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Access-Control-Allow-Headers",
			"Content-Type", "Content-Length", "Accept-Encoding",
			"X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.StaticFS("/static", http.Dir("../../web/build/static"))
	r.LoadHTMLFiles("../../web/build/index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/ws", func(c *gin.Context) {
		websocketconn.WsConn(c.Writer, c.Request)
	})

	r.GET("/config/rs485", func(c *gin.Context) {
		devices.GetRS485Config(c)
	})

	r.POST("/config/rs485", devices.PostRS485Config)

	r.Run()
}

// err := os.Remove("users.json")
// if err != nil {
// 	fmt.Println(err)
// }
// jsonFile, err := os.Create("users.json")
// if err != nil {
// 	fmt.Println(err)
// }
// defer jsonFile.Close()
// // decoder := json.NewDecoder(c.Request.Body)
// // decoder.Decode(&jsonFile)
// scanner := bufio.NewScanner(c.Request.Body)
// for scanner.Scan() {
// 	//fmt.Println(scanner.Text())
// 	//jsonFile.WriteString(scanner.Text())
// 	fmt.Fprintln(jsonFile, scanner.Text())
// }
// //json.NewDecoder(c.Request.Body).Decode(encoder)
// c.JSON(200, gin.H{
// 	"status": "posted",
// })
// // devices.PostRS485Config(c.Writer, c.Request)
