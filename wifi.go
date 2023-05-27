package main

import (
	"encoding/json"
	"machine"
	"strconv"
	"time"

	"tinygo.org/x/drivers/net/mqtt"
	"tinygo.org/x/drivers/wifinina"
)

// IP address of the MQTT broker to use. Replace with your own info.
const server = MQTTProtocol + "://" + MQTTServer + ":" + MQTTPort

//const server = "ssl://test.mosquitto.org:8883"

// change these to connect to a different UART or pins for the ESP8266/ESP32
var (
	// these are the default pins for the Arduino Nano33 IoT.
	spi = machine.NINA_SPI

	// this is the ESP chip that has the WIFININA firmware flashed on it
	adaptor *wifinina.Device

	cl             mqtt.Client
	token          mqtt.Token
	topicPublish   = DeviceID
	topicSubscribe = DeviceID + "-call"

	connectedWifi bool
	connectedMQTT bool

	relayID     string
	relayStatus string
)

func subHandler(client mqtt.Client, msg mqtt.Message) {
	println("[", msg.Topic(), "] ", string(msg.Payload()))
	var fns []Method
	err := json.Unmarshal(msg.Payload(), &fns)
	if err != nil {
		println("ERROR UnMarshalling rabbitf3-call payload", err)
		return
	}
	for _, f := range fns {
		println(f.Name)
		switch f.Name {
		case "info":
			data, err = json.Marshal(DiscoveryMsg)
			if err != nil {
				println("[INFO]", err)
			}
			publishData("discovery", &data)
			break
		case "gm":
			break
		case "sm":
			break
		case "grtc":
			break
		case "relay":
			relayID = ""
			relayStatus = ""
			for _, p := range f.Params {
				if p.ID == "r" {
					relayID = p.Value.(string)
				} else if p.ID == "s" {
					relayStatus = p.Value.(string)
				}
			}
			if relayID != "" && relayStatus != "" {
				i, _ := strconv.Atoi(relayID)
				if relayStatus == "1" || relayStatus == "on" {
					relay[i].High()
				} else if relayStatus == "0" || relayStatus == "off" {
					relay[i].Low()
				}
			}
			break
		case "food":
			break
		default:
		}
	}
}

func connectToAP() {
	time.Sleep(2 * time.Second)
	for {
		println("Connecting to " + WifiSSID)
		err := adaptor.ConnectToAccessPoint(WifiSSID, WifiPassword, 10*time.Second)
		if err == nil { // error connecting to AP
			println("Connected.")

			time.Sleep(2 * time.Second)
			ip, _, _, err := adaptor.GetIP()
			for ; err != nil; ip, _, _, err = adaptor.GetIP() {
				println("[GET IP]", err.Error())
				time.Sleep(1 * time.Second)
			}
			println("[IP]", ip.String())
			connectedWifi = true
			break
		} else {
			println("[CONNECT TO AP]", err)
			println("Waiting 30s before trying to reconnect")
			connectedWifi = false
			time.Sleep(30 * time.Second)
		}
	}
}

func connectToMQTT() {
	println("Connecting to MQTT")

	opts := mqtt.NewClientOptions().AddBroker(server)
	opts.SetClientID(MQTTClientID)
	opts.SetUsername(MQTTUser)
	opts.SetPassword(MQTTPassword)

	cl = mqtt.NewClient(opts)

	if token := cl.Connect(); token.Wait() && token.Error() != nil {
		println("[MQTT CONNECT]", token.Error().Error())
		connectedMQTT = false
	}
	println("Connected to MQTT")

	// subscribe
	token := cl.Subscribe("#", 0, subHandler)
	token.Wait()
	if token.Error() != nil {
		println("[MQTT SUBSCRIBE]", token.Error().Error())
		connectedMQTT = false
	}
	connectedMQTT = true
}

func publishDiscovery() {
	// DISCOVERY MESSAGE
	println("Marshalling Discovery Message, if no action after this, increase stack size with --stack-size 10KB")
	data, err = json.Marshal(ShortDiscoveryMsg)
	if err != nil {
		println("[DISCOVERY]", err)
	}
	println("[DISCOVERY]", string(data))
	token := cl.Publish("discovery", 0, false, data)
	token.Wait()
	if token.Error() != nil {
		println("[DISCOVERY]", token.Error().Error())
	}
}

func publishData(topic string, data *[]byte) {
	println("[PUBLISH DATA]", "#"+topic, "MSG TO SEND", string(*data))
	token = cl.Publish(topic, 0, false, *data)
	token.Wait()
	if token.Error() != nil {
		println("[PUBLISH DATA]", token.Error().Error())
		println("[PUBLISH DATA]", "Retrying publish...")
		connectToMQTT()
		token = cl.Publish(topic, 0, false, *data)
		token.Wait()
		if token.Error() != nil {
			println("[PUBLISH DATA]", token.Error().Error())
		}
	}
}
