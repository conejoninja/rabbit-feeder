package main

import (
	"encoding/json"
	"machine"
	"time"

	"tinygo.org/x/drivers/at24cx"
	"tinygo.org/x/drivers/wifinina"

	"tinygo.org/x/drivers/bme280"
	"tinygo.org/x/drivers/ds3231"
	"tinygo.org/x/drivers/vl6180x"
)

var (
	dirPin   machine.Pin
	stepPin  machine.Pin
	sleepPin machine.Pin
	relay    [4]machine.Pin
)

var (
	distanceSensor    vl6180x.Device
	rtc               ds3231.Device
	temperatureSensor bme280.Device
	eeprom            at24cx.Device

	distanceSensorEnabled    bool
	temperatureSensorEnabled bool
	rtcEnabled               bool
	eepromEnabled            bool

	sensorState SensorState
	relayState  RelayState
	data        []byte
	err         error

	distance   uint16
	dt         time.Time
	temp       int32
	n          int
	eepromData []byte
)

const (
	Alarm1     = 0
	LastAlarm1 = 4
	NextAlarm1 = 12
	Quantity1  = 20
	Alarm2     = 24
	LastAlarm2 = 28
	NextAlarm3 = 36
	Quantity4  = 44
	NextRecord = 48
)

const (
	DISTANCE = iota
	DISTANCE_RAW
	MEMORY
	TEMPERATURE
	PRESSURE
	HUMIDITY
	RTC
)

func main() {

	time.Sleep(5 * time.Second)
	println("online")

	// SETUP RELAY
	relay = [4]machine.Pin{
		machine.D5,
		machine.D4,
		machine.D3,
		machine.D2,
	}
	for i := 0; i < 4; i++ {
		relay[i].Configure(machine.PinConfig{Mode: machine.PinOutput})
		relay[i].Low()
	}

	// SETUP THE MOTOR
	dirPin = machine.D10
	stepPin = machine.D9
	sleepPin = machine.D8
	dirPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	stepPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	sleepPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	sleepPin.High()

	machine.I2C0.Configure(machine.I2CConfig{})

	// SETUP DISTANCE SENSOR
	distanceSensor = vl6180x.New(machine.I2C0)
	distanceSensorEnabled = distanceSensor.Connected()
	if !distanceSensorEnabled {
		println("VL6180X device not found")
	} else {
		distanceSensor.Configure(true)
	}

	// SETUP BME280
	temperatureSensor = bme280.New(machine.I2C0)
	temperatureSensor.Configure()

	// SETUP RTC
	rtc = ds3231.New(machine.I2C0)
	rtc.Configure()

	valid := rtc.IsTimeValid()
	if !valid {
		println("DATE IS NOT VALID")
		date := time.Date(2023, 05, 14, 15, 49, 07, 0, time.UTC)
		rtc.SetTime(date)
	}

	rtcEnabled = rtc.IsRunning()
	if !rtcEnabled {
		err := rtc.SetRunning(true)
		if err != nil {
			println("Error configuring RTC")
		} else {
			rtcEnabled = true
		}
	}

	// SETUP EEPROM
	eeprom := at24cx.New(machine.I2C0)
	eeprom.Configure(at24cx.Config{})
	eepromEnabled = true // assume it's working
	eepromData = make([]byte, 48)

	// Configure SPI for 8Mhz, Mode 0, MSB First
	spi.Configure(machine.SPIConfig{
		Frequency: 8 * 1e6,
		SDO:       machine.NINA_SDO,
		SDI:       machine.NINA_SDI,
		SCK:       machine.NINA_SCK,
	})

	// Init esp8266/esp32
	adaptor = wifinina.New(spi,
		machine.NINA_CS,
		machine.NINA_ACK,
		machine.NINA_GPIO0,
		machine.NINA_RESETN)
	adaptor.Configure()

	connectToAP()
	connectToMQTT()
	publishDiscovery()
	// Let discovery message to be processed and other devices subscribe to it
	time.Sleep(2 * time.Second)

	for {

		sendSensorStatus()
		sendRelayStatus()

		time.Sleep(time.Second * 60)
	}
}

func sendSensorStatus() {
	distance = distanceSensor.Read()
	println("Distance:", distance)
	sensorState.Distance = distance

	dt, err = rtc.ReadTime()
	if err != nil {
		println("Error reading date:", err)
	} else {
		println(dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), dt.Second())
	}
	sensorState.Date = dt.Format(time.RFC3339)

	/*temp, _ = rtc.ReadTemperature()
	println("Temperature (RTC):", temp)
	sensorState.Temperature = temp*/

	temp, _ = temperatureSensor.ReadTemperature()
	println("Temperature (BME280):", temp)
	sensorState.Temperature = temp
	temp, _ = temperatureSensor.ReadPressure()
	println("Pressure (BME280):", temp)
	sensorState.Pressure = temp
	temp, _ = temperatureSensor.ReadHumidity()
	println("Humidity (BME280):", temp)
	sensorState.Humidity = temp

	n, err = eeprom.Read(eepromData)
	println(n, err)
	for i := 0; i < 48; i++ {
		print(eepromData[i])
	}
	println("==========")
	sensorState.EEPROM = eepromData

	data, err = json.Marshal(sensorState)
	if err != nil {
		println("ERROR MARSHALLING SENSOR DATA", err)
	} else {
		publishData(sensorStateTopic, &data)
	}
}

func sendRelayStatus() {
	relayState.Relay1 = "OFF"
	relayState.Relay2 = "OFF"
	relayState.Relay3 = "OFF"
	relayState.Relay4 = "OFF"
	if relay[0].Get() {
		relayState.Relay1 = "ON"
	}
	if relay[1].Get() {
		relayState.Relay2 = "ON"
	}
	if relay[2].Get() {
		relayState.Relay3 = "ON"
	}
	if relay[3].Get() {
		relayState.Relay4 = "ON"
	}
	data, err = json.Marshal(relayState)
	if err != nil {
		println("ERROR MARSHALLING RELAY DATA", err)
	} else {
		publishData(relayStateTopic, &data)
	}
}
