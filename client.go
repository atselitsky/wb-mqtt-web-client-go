package main

import (
	"fmt"
	"time"

	"github.com/atselitsky/wb-mqtt-web-client-go/helpers"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	conn   *websocketConnection
	topics map[string]byte
}

func (client *Client) MqttMesHandler(c mqtt.Client, m mqtt.Message) {
	msg := &MQTTmessage{
		Path:    m.Topic(),
		Message: string(m.Payload()),
		When:    time.Now(),
	}
	client.conn.forward <- msg
	fmt.Printf("Expected message from %s: %s \n", m.Topic(), m.Payload())
}

func createMQTTPahoClient() mqtt.Client {
	client := helpers.ReadClientParams()
	return client.CreateClient()
}

func NewMQTTConn(c *websocketConnection) *Client {
	return &Client{
		conn: c,
		topics: map[string]byte{
			"/test/topic":    1,
			"/test/topic/on": 1,
			"/devices/wb-m1w2_107/controls/Internal Temperature": 1,
			"/devices/hwmon/controls/Board Temperature":          1,
		},
	}
}

func (c *Client) StartMQTTConnection() {
	client := createMQTTPahoClient()
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// token := client.Subscribe("/test/topic", 1, c.MqttMesHandler)
	token := client.SubscribeMultiple(c.topics, c.MqttMesHandler)
	token.Wait()

}
