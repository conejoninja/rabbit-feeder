package main

import (
	"encoding/json"
	"machine"
	"strings"
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

	cl    mqtt.Client
	token mqtt.Token

	connectedWifi bool
	connectedMQTT bool

	relayID     string
	relayStatus string
)

func relayHandler(topics []string, payload string) {
	if len(topics) < 4 || topics[3] != "set" {
		return
	}

	switch topics[2] {
	case "relay1":
		if payload == "ON" {
			relay[0].High()
		} else if payload == "OFF" {
			relay[0].Low()
		}
		break
	case "relay2":
		if payload == "ON" {
			relay[1].High()
		} else if payload == "OFF" {
			relay[1].Low()
		}
		break
	case "relay3":
		if payload == "ON" {
			relay[2].High()
		} else if payload == "OFF" {
			relay[2].Low()
		}
		break
	case "relay4":
		if payload == "ON" {
			relay[3].High()
		} else if payload == "OFF" {
			relay[3].Low()
		}
		break
	default:
		break
	}
	sendRelayStatus()
}

func subHandler(client mqtt.Client, msg mqtt.Message) {
	println("[", msg.Topic(), "] ", string(msg.Payload()))
	topics := strings.Split(msg.Topic(), "/")
	if topics[0] != "homeassistant" {
		return
	}
	if len(topics) > 1 && topics[1] == "switch" {
		relayHandler(topics, string(msg.Payload()))
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
			println("Waiting 15s before trying to reconnect")
			connectedWifi = false
			time.Sleep(15 * time.Second)
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
	println("Marshalling Discovery Messages, if no action after this, increase stack size with --stack-size 10KB")
	data, err = json.Marshal(Relay1Discovery)
	if err != nil {
		println("[DISCOVERY]", err)
	}
	println("[DISCOVERY]", string(data))
	token := cl.Publish(Relay1Discovery.Home+"/config", 0, false, data)
	token.Wait()
	if token.Error() != nil {
		println("[DISCOVERY]", token.Error().Error())
	}

	data, err = json.Marshal(Relay2Discovery)
	if err != nil {
		println("[DISCOVERY]", err)
	}
	println("[DISCOVERY]", string(data))
	token = cl.Publish(Relay2Discovery.Home+"/config", 0, false, data)
	token.Wait()
	if token.Error() != nil {
		println("[DISCOVERY]", token.Error().Error())
	}

	data, err = json.Marshal(Relay3Discovery)
	if err != nil {
		println("[DISCOVERY]", err)
	}
	println("[DISCOVERY]", string(data))
	token = cl.Publish(Relay3Discovery.Home+"/config", 0, false, data)
	token.Wait()
	if token.Error() != nil {
		println("[DISCOVERY]", token.Error().Error())
	}

	data, err = json.Marshal(Relay4Discovery)
	if err != nil {
		println("[DISCOVERY]", err)
	}
	println("[DISCOVERY]", string(data))
	token = cl.Publish(Relay4Discovery.Home+"/config", 0, false, data)
	token.Wait()
	if token.Error() != nil {
		println("[DISCOVERY]", token.Error().Error())
	}

	data, err = json.Marshal(TemperatureDiscovery)
	if err != nil {
		println("[DISCOVERY]", err)
	}
	println("[DISCOVERY]", string(data))
	token = cl.Publish(TemperatureDiscovery.Home+"/config", 0, false, data)
	token.Wait()
	if token.Error() != nil {
		println("[DISCOVERY]", token.Error().Error())
	}

	data, err = json.Marshal(PressureDiscovery)
	if err != nil {
		println("[DISCOVERY]", err)
	}
	println("[DISCOVERY]", string(data))
	token = cl.Publish(PressureDiscovery.Home+"/config", 0, false, data)
	token.Wait()
	if token.Error() != nil {
		println("[DISCOVERY]", token.Error().Error())
	}

	data, err = json.Marshal(HumidityDiscovery)
	if err != nil {
		println("[DISCOVERY]", err)
	}
	println("[DISCOVERY]", string(data))
	token = cl.Publish(HumidityDiscovery.Home+"/config", 0, false, data)
	token.Wait()
	if token.Error() != nil {
		println("[DISCOVERY]", token.Error().Error())
	}

	data, err = json.Marshal(DistanceDiscovery)
	if err != nil {
		println("[DISCOVERY]", err)
	}
	println("[DISCOVERY]", string(data))
	token = cl.Publish(DistanceDiscovery.Home+"/config", 0, false, data)
	token.Wait()
	if token.Error() != nil {
		println("[DISCOVERY]", token.Error().Error())
	}

	data, err = json.Marshal(EEPROMDiscovery)
	if err != nil {
		println("[DISCOVERY]", err)
	}
	println("[DISCOVERY]", string(data))
	token = cl.Publish(EEPROMDiscovery.Home+"/config", 0, false, data)
	token.Wait()
	if token.Error() != nil {
		println("[DISCOVERY]", token.Error().Error())
	}

	data, err = json.Marshal(RTCDiscovery)
	if err != nil {
		println("[DISCOVERY]", err)
	}
	println("[DISCOVERY]", string(data))
	token = cl.Publish(RTCDiscovery.Home+"/config", 0, false, data)
	token.Wait()
	if token.Error() != nil {
		println("[DISCOVERY]", token.Error().Error())
	}

	data, err = json.Marshal(MotorDiscovery)
	if err != nil {
		println("[DISCOVERY]", err)
	}
	println("[DISCOVERY]", string(data))
	token = cl.Publish(MotorDiscovery.Home+"/config", 0, false, data)
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
