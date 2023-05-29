package main

const DeviceID = "rabbitf3"

type Discovery struct {
	Home              string `json:"~"`
	Name              string `json:"name,omitempty"`
	UniqueID          string `json:"unique_id,omitempty"`
	ObjectID          string `json:"object_id,omitempty"`
	UnitOfMeasurement string `json:"unit_of_measurement,omitempty"`
	ValueTemplate     string `json:"value_template,omitempty"`
	CommandTopic      string `json:"cmd_t,omitempty"`
	StatusTopic       string `json:"stat_t,omitempty"`
	Device            Device `json:"device,omitempty"`
	Icon              string `json:"icon,omitempty"`
}

type Device struct {
	Identifiers  []string `json:"identifiers,omitempty"`
	Name         string   `json:"name,omitempty"`
	Model        string   `json:"model,omitempty"`
	Manufacturer string   `json:"manufacturer,omitempty"`
}

type SensorState struct {
	Temperature int32  `json:"temperature,omitempty"`
	Humidity    int32  `json:"humidity,omitempty"`
	Pressure    int32  `json:"pressure,omitempty"`
	Distance    uint16 `json:"distance,omitempty"`
	EEPROM      []byte `json:"eeprom,omitempty"`
	Date        string `json:"date,omitempty"`
}

type RelayState struct {
	Relay1 string `json:"relay1,omitempty"`
	Relay2 string `json:"relay2,omitempty"`
	Relay3 string `json:"relay3,omitempty"`
	Relay4 string `json:"relay4,omitempty"`
}

var device = Device{
	Identifiers:  []string{"rabbitf3"},
	Name:         "Rabbit Feeder Supreme",
	Model:        "Rabbit Feeder Supreme F3",
	Manufacturer: "@conejo@social.tinygo.org",
}

var Relay1Discovery = Discovery{
	Home:          "homeassistant/switch/relay1",
	Name:          "Relay 1 (USB)",
	UniqueID:      DeviceID + "_relay1",
	ObjectID:      DeviceID + "_relay1",
	ValueTemplate: "{{ value_json.relay1 }}",
	CommandTopic:  "~/set",
	StatusTopic:   "homeassistant/switch/relays/state",
	Device:        device,
	Icon:          "mdi:usb-port",
}

var Relay2Discovery = Discovery{
	Home:          "homeassistant/switch/relay2",
	Name:          "Relay 2 (USB)",
	UniqueID:      DeviceID + "_relay2",
	ObjectID:      DeviceID + "_relay2",
	ValueTemplate: "{{ value_json.relay2 }}",
	CommandTopic:  "~/set",
	StatusTopic:   "homeassistant/switch/relays/state",
	Device:        device,
	Icon:          "mdi:usb-port",
}

var Relay3Discovery = Discovery{
	Home:          "homeassistant/switch/relay3",
	Name:          "Relay 3 (12V)",
	UniqueID:      DeviceID + "_relay3",
	ObjectID:      DeviceID + "_relay3",
	ValueTemplate: "{{ value_json.relay3 }}",
	CommandTopic:  "~/set",
	StatusTopic:   "homeassistant/switch/relays/state",
	Device:        device,
	Icon:          "mdi:audio-input-stereo-minijack",
}

var Relay4Discovery = Discovery{
	Home:          "homeassistant/switch/relay4",
	Name:          "Relay 4 (12V)",
	UniqueID:      DeviceID + "_relay4",
	ObjectID:      DeviceID + "_relay4",
	ValueTemplate: "{{ value_json.relay4 }}",
	CommandTopic:  "~/set",
	StatusTopic:   "homeassistant/switch/relays/state",
	Device:        device,
	Icon:          "mdi:audio-input-stereo-minijack",
}

var TemperatureDiscovery = Discovery{
	Home:              "homeassistant/sensor/temperature",
	Name:              "Temperature",
	UniqueID:          DeviceID + "_temp",
	ObjectID:          DeviceID + "_temp",
	UnitOfMeasurement: "ºC",
	ValueTemplate:     "{{ value_json.temperature / 1000 }}",
	StatusTopic:       "homeassistant/switch/sensors/state",
	Device:            device,
	Icon:              "mdi:thermometer",
}

var HumidityDiscovery = Discovery{
	Home:              "homeassistant/sensor/humidity",
	Name:              "Humidity",
	UniqueID:          DeviceID + "_humidity",
	ObjectID:          DeviceID + "_humidity",
	UnitOfMeasurement: "%",
	ValueTemplate:     "{{ value_json.humidity / 100 }}",
	StatusTopic:       "homeassistant/switch/sensors/state",
	Device:            device,
	Icon:              "mdi:water-percent",
}

var PressureDiscovery = Discovery{
	Home:              "homeassistant/sensor/pressure",
	Name:              "Pressure",
	UniqueID:          DeviceID + "_pressure",
	ObjectID:          DeviceID + "_pressure",
	UnitOfMeasurement: "Pa",
	ValueTemplate:     "{{ value_json.pressure / 1000 }}",
	StatusTopic:       "homeassistant/switch/sensors/state",
	Device:            device,
	Icon:              "mdi:air-filter",
}

var DistanceDiscovery = Discovery{
	Home:              "homeassistant/sensor/distance",
	Name:              "Distance",
	UniqueID:          DeviceID + "_dist",
	ObjectID:          DeviceID + "_dist",
	UnitOfMeasurement: "mm",
	ValueTemplate:     "{{ value_json.distance }}",
	StatusTopic:       "homeassistant/switch/sensors/state",
	Device:            device,
	Icon:              "mdi:gauge-full",
}

var RTCDiscovery = Discovery{
	Home:          "homeassistant/statestream/rtc",
	Name:          "RTC",
	UniqueID:      DeviceID + "_rtc",
	ObjectID:      DeviceID + "_rtc",
	ValueTemplate: "{{ value_json.timestamp }}",
	StatusTopic:   "homeassistant/switch/sensors/state",
	Device:        device,
	Icon:          "mdi:clock-digital",
}

var EEPROMDiscovery = Discovery{
	Home:          "homeassistant/text/eeprom",
	Name:          "EEPROM",
	UniqueID:      DeviceID + "_eeprom",
	ObjectID:      DeviceID + "_eeprom",
	ValueTemplate: "{{ value_json.eeprom }}",
	StatusTopic:   "homeassistant/switch/sensors/state",
	Device:        device,
	Icon:          "mdi:text-box",
}

var MotorDiscovery = Discovery{
	Home:         "homeassistant/binary_sensor/motor",
	Name:         "Motor",
	UniqueID:     DeviceID + "_motor",
	ObjectID:     DeviceID + "_motor",
	CommandTopic: "~/set",
	StatusTopic:  "~/state",
	Device:       device,
	Icon:         "mdi:engine",
}
