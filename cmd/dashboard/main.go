package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var subscriptions map[string]bool
var token mqtt.Token
var c mqtt.Client

func main() {

	subscriptions = make(map[string]bool)

	opts := mqtt.NewClientOptions().AddBroker(MQTTProtocol + "://" + MQTTServer + ":" + MQTTPort)
	opts.SetClientID(MQTTClientID)
	opts.SetUsername(MQTTUser)
	opts.SetPassword(MQTTPassword)
	opts.SetDefaultPublishHandler(defaultHandler)

	c = mqtt.NewClient(opts)
	if token = c.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token)
		log.Println(token.Error())
		panic(token.Error())
	}
	defer c.Disconnect(250)

	// Discover new devices when they connect to the network
	if token = c.Subscribe("discovery", 0, discoveryHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	/*
		if token = c.Subscribe("events", 0, eventsHandler); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}*/

	for {
	}
}

/*
var discoveryHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("DISCOVERY: %s\n", msg.Payload())
	var device common.Device
	err := json.Unmarshal(msg.Payload(), &device)
	if err == nil {
		//db.AddDevice([]byte(device.ID), device)
		if v, ok := subscriptions[device.ID]; !ok || !v {
			subscriptions[device.ID] = true
			if token = c.Subscribe(device.ID, 0, nil); token.Wait() && token.Error() != nil {
				subscriptions[device.ID] = false
				log.Println(token.Error())
				os.Exit(1)
			}

		}
	} else {
		log.Println(err)
	}
}
*/

var defaultHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("[%s]: %s\n", msg.Topic(), msg.Payload())
}

var discoveryHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("[%s]: %s\n", msg.Topic(), msg.Payload())
	var discovery Device
	err := json.Unmarshal(msg.Payload(), &discovery)
	if err != nil {
		fmt.Println("Error Unmarshalling DISCOVERY", err)
		return
	}

	if v, ok := subscriptions[discovery.ID]; !ok || !v {
		subscriptions[discovery.ID] = true
		if token = c.Subscribe(discovery.ID, 0, nil); token.Wait() && token.Error() != nil {
			subscriptions[discovery.ID] = false
			log.Println(token.Error())
			os.Exit(1)
		}

	}
}
