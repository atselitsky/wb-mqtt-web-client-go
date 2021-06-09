package configs

import (
	"fmt"
	"log"
	"net/url"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/yaml.v3"
)

type ClientParams interface {
	CreateClient() mqtt.Client
	//OnConnectHandler()
	//OnConnectLost()
}

type clientParams struct {
	// clientLocation     string
	ClientID string `yaml:"ClientID"`
	// broker             string
	// port               int
	Username              string `yaml:"Username"`
	Password              string `yaml:"Password"`
	OnConnect             mqtt.OnConnectHandler
	DefaultPublishHandler mqtt.MessageHandler
}

func (o *clientParams) messagePubHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Unexpected message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

func (o *clientParams) onConnectHandler(client mqtt.Client) {
	fmt.Println("Connected")
}

func (o *clientParams) CreateClient() mqtt.Client {
	servers, err := url.Parse("tcp://localhost:1883")
	if err != nil {
		log.Fatal(err)
	}
	options := mqtt.ClientOptions{
		Servers:               []*url.URL{servers},
		ClientID:              o.ClientID,
		Username:              o.Username,
		Password:              o.Password,
		OnConnect:             o.onConnectHandler,
		DefaultPublishHandler: o.messagePubHandler,
	}

	client := mqtt.NewClient(&options)
	return client
}

func openConfigFile(name string) *os.File {
	f, err := os.Open(name)
	if err != nil {
		log.Fatalf("In parseConfig error: %s\n", err)
	}
	return f
}

//  * ReadClientParam reads mqtt client parameters from yaml file, decode it and return a  Client parameters
//  Yaml file path: /config/ClientOptions.yaml
// TODO Change exact path to specified(and add default if it's not)

func ReadClientParams() ClientParams {
	var o *clientParams

	f := openConfigFile("../../pkg/configs/ClientOptions.yaml")
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err := decoder.Decode(&o)
	if err != nil {
		log.Fatal(err)
	}

	return o
}
