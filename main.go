package main

import (
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	HIVE_BROKER = "tcp://broker.hivemq.com:1883"
)

var messageHandler = func(client mqtt.Client, message mqtt.Message) {
	log.Printf("Topic: %s, Message: %s", message.Topic(), message.Payload())
}

func main() {

	// connect with options
	opts := mqtt.NewClientOptions().AddBroker(HIVE_BROKER).SetClientID("sample")

	// start client
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// subscribe to a topic
	if token := c.Subscribe("example/topic", 0, messageHandler); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		os.Exit(1)
	}

	// publish and wait from the server receipt after sending a message
	for i := 0; i < 5; i++ {
		text := log.Sprintf("Hello. this is message %d", i)
		token := c.Publish("example/topic", 0, false, text)
		token.Wait()
	}

	time.Sleep(3 * time.Second)

	// unsubscribe
	if token := c.Unsubscribe("exmaple/topic"); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)
}
