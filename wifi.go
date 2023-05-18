package main

import (
	"fmt"
	"machine"
	"time"
	"tinygo.org/x/drivers/net/mqtt"
	"tinygo.org/x/drivers/wifinina"
)

// IP address of the MQTT broker to use. Replace with your own info.
const server = "tcp://test.mosquitto.org:1883"

//const server = "ssl://test.mosquitto.org:8883"

// change these to connect to a different UART or pins for the ESP8266/ESP32
var (
	// these are the default pins for the Arduino Nano33 IoT.
	spi = machine.NINA_SPI

	// this is the ESP chip that has the WIFININA firmware flashed on it
	adaptor *wifinina.Device

	cl      mqtt.Client
	topicTx = "rabbitfeeder/tx"
	topicRx = "rabbitfeeder/rx"
)

func subHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("[%s]  ", msg.Topic())
	fmt.Printf("%s\r\n", msg.Payload())
}

func publishing() {
	for i := 0; ; i++ {
		println("Publishing MQTT message...")
		data := []byte(fmt.Sprintf(`{"e":[{"n":"hello %d","v":101}]}`, i))
		token := cl.Publish(topicRx, 0, false, data)
		token.Wait()
		if token.Error() != nil {
			println(token.Error().Error())
		}

		time.Sleep(30000 * time.Millisecond)
	}
}

// connect to access point
func connectToAP() {
	time.Sleep(2 * time.Second)
	println("Connecting to " + WifiSSID)
	err := adaptor.ConnectToAccessPoint(WifiSSID, WifiPassword, 10*time.Second)
	if err != nil { // error connecting to AP
		for {
			println(err)
			time.Sleep(1 * time.Second)
		}
	}

	println("Connected.")

	time.Sleep(2 * time.Second)
	ip, _, _, err := adaptor.GetIP()
	for ; err != nil; ip, _, _, err = adaptor.GetIP() {
		println(err.Error())
		time.Sleep(1 * time.Second)
	}
	println(ip.String())
}

func failMessage(msg string) {
	for {
		println(msg)
		time.Sleep(1 * time.Second)
	}
}