package main

import "time"

const DeviceID = "rabbitf3"

type Device struct {
	ID      string   `json:"id"`
	Name    string   `json:"name,omitempty"`
	Version string   `json:"version,omitempty"`
	Out     []Value  `json:"out,omitempty"`
	Methods []Method `json:"methods,omitempty"`
}

type Value struct {
	ID    string      `json:"id"`
	Type  string      `json:"type,omitempty"`
	Name  string      `json:"name,omitempty"`
	Unit  string      `json:"unit,omitempty"`
	Time  *time.Time  `json:"time,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

type Method struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Params      []Value `json:"params,omitempty"`
}

type Param struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Event struct {
	ID       string     `json:"id"`
	Message  string     `json:"message,omitempty"`
	Priority uint8      `json:"priority,omitempty"`
	Time     *time.Time `json:"time,omitempty"`
	Extra    []Param    `json:"extra,omitempty"`
}

var ShortDiscoveryMsg = Device{
	ID:   DeviceID,
	Name: "Rabbit Feeder Supreme",
}

var DiscoveryMsg = Device{
	ID:      DeviceID,
	Name:    "Rabbit Feeder Supreme",
	Version: "v1.0.0",
	Out: []Value{
		Value{
			ID:   "c",
			Name: "Capacity",
			Unit: "%",
		},
		Value{
			ID:   "cr",
			Name: "Capacity Raw",
			Unit: "mm",
		},
		Value{
			ID:   "m",
			Name: "EEPROM",
			Unit: "byte",
		},
		Value{
			ID:   "t",
			Name: "Temperature",
			Unit: "mC",
		},
		Value{
			ID:   "p",
			Name: "Pressure",
			Unit: "mPa",
		},
		Value{
			ID:   "h",
			Name: "Humidity",
			Unit: "percent",
		},
		Value{
			ID:   "rtc",
			Name: "Datetime",
			Unit: "Time",
		},
	},
	Methods: []Method{
		Method{
			Name:        "gm",
			Description: "Get the EEPROM",
		},
		Method{
			Name:        "sm",
			Description: "Set the EEPROM",
			Params: []Value{
				Value{
					ID:   "p",
					Name: "Byte position",
				},
				Value{
					ID:   "v",
					Name: "Value",
				},
			},
		},
		Method{
			Name:        "grtc",
			Description: "Get the RTC datetime",
		},
	},
}
