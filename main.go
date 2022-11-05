package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/ds3231"
	"tinygo.org/x/drivers/vl6180x"
)

func main() {
	machine.I2C0.Configure(machine.I2CConfig{})

	rtc := ds3231.New(machine.I2C0)
	rtc.Configure()

	distanceSensor := vl6180x.New(machine.I2C0)
	connected := distanceSensor.Connected()
	if !connected {
		println("VL6180X device not found")
		return
	}

	distanceSensor.Configure(true)

	valid := rtc.IsTimeValid()
	if !valid {
		date := time.Date(2022, 11, 05, 16, 11, 07, 0, time.UTC)
		rtc.SetTime(date)
	}

	running := rtc.IsRunning()
	if !running {
		err := rtc.SetRunning(true)
		if err != nil {
			println("Error configuring RTC")
		}
	}

	value := distanceSensor.Read()
	for {
		dt, err := rtc.ReadTime()
		if err != nil {
			println("Error reading date:", err)
		} else {
			println(dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), dt.Second())
		}
		temp, _ := rtc.ReadTemperature()
		println("Temperature:", temp)
		value = distanceSensor.Read()
		println("Distancia:", value)

		time.Sleep(time.Second * 1)
	}
}
