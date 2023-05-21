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

	sensorData [7]Value
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

	// SETUP RELAY
	relay = [4]machine.Pin{
		machine.D5,
		machine.D4,
		machine.D3,
		machine.D2,
	}
	for i := 0; i < 4; i++ {
		relay[i].Configure(machine.PinConfig{Mode: machine.PinOutput})
	}

	// SETUP THE MOTOR
	dirPin = machine.D10
	stepPin = machine.D9
	sleepPin = machine.D8
	dirPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	stepPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	sleepPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	sleepPin.High()

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

	// DISCOVERY MESSAGE
	println("Marshalling Discovery Message, if no action after this, increase stack size with --stack-size 10KB")
	data, err := json.Marshal(ShortDiscoveryMsg)
	if err != nil {
		println("ERROR DISCOVERY", err)
	}
	token := cl.Publish("discovery", 0, false, data)
	token.Wait()
	if token.Error() != nil {
		println(token.Error().Error())
	}

	//go publishing()

	value := distanceSensor.Read()
	var dt time.Time
	var temp int32
	var n int
	eepromData := make([]byte, 48)

	sensorData[DISTANCE] = Value{
		ID: "c",
	}
	sensorData[DISTANCE_RAW] = Value{
		ID: "cr",
	}
	sensorData[MEMORY] = Value{
		ID: "m",
	}
	sensorData[TEMPERATURE] = Value{
		ID: "t",
	}
	sensorData[PRESSURE] = Value{
		ID: "p",
	}
	sensorData[HUMIDITY] = Value{
		ID: "h",
	}
	sensorData[RTC] = Value{
		ID: "rtc",
	}

	for {
		value = distanceSensor.Read()
		println("Distance:", value)
		sensorData[DISTANCE].Value = value
		sensorData[DISTANCE_RAW].Value = value

		dt, err = rtc.ReadTime()
		if err != nil {
			println("Error reading date:", err)
		} else {
			println(dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), dt.Second())
		}
		sensorData[RTC].Value = dt.Unix()

		temp, _ = rtc.ReadTemperature()
		println("Temperature (RTC):", temp)
		sensorData[TEMPERATURE].Value = temp

		temp, _ = temperatureSensor.ReadTemperature()
		println("Temperature (BME280):", temp)
		sensorData[TEMPERATURE].Value = temp
		temp, _ = temperatureSensor.ReadPressure()
		println("Pressure (BME280):", temp)
		sensorData[PRESSURE].Value = temp
		temp, _ = temperatureSensor.ReadHumidity()
		println("Humidity (BME280):", temp)
		sensorData[HUMIDITY].Value = temp

		n, err = eeprom.Read(eepromData)
		println(n, err)
		for i := 0; i < 48; i++ {
			print(eepromData[i])
		}
		println("==========")
		sensorData[MEMORY].Value = eepromData

		data, err = json.Marshal(sensorData)
		if err != nil {
			println("ERROR MARSHALLING SENSOR DATA", err)
		} else {
			token = cl.Publish(DeviceID, 0, false, data)
			token.Wait()
			if token.Error() != nil {
				println(token.Error().Error())
				println("Retrying publish...")
				connectToMQTT()
				token = cl.Publish(DeviceID, 0, false, data)
				token.Wait()
				if token.Error() != nil {
					println(token.Error().Error())
				}
			}
		}

		time.Sleep(time.Second * 60)
	}

	/*i := 0
	d := 0
	for {
		if d == 0 {
			relay[i].High()
		} else {
			relay[i].Low()
		}
		i++
		if i >= 4 {
			i = 0
			if d == 0 {
				d = 1
			} else {
				d = 0
			}
		}
		time.Sleep(1 * time.Second)
	}

	for {
		println("high")

		dirPin.High()
		for s := 0; s < 6000; s++ {
			stepPin.High()
			time.Sleep(1 * time.Millisecond)
			stepPin.Low()
			time.Sleep(1 * time.Millisecond)
		}

		time.Sleep(2 * time.Second)
		println("low")
		dirPin.Low()
		for s := 0; s < 600; s++ {
			stepPin.High()
			time.Sleep(1 * time.Millisecond)
			stepPin.Low()
			time.Sleep(1 * time.Millisecond)
		}
		time.Sleep(2 * time.Second)
	}*/

}
