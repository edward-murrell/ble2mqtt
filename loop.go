package main

import (
	"encoding/json"
	"fmt"
	"gobot.io/x/gobot/platforms/mqtt"
	"tinygo.org/x/bluetooth"
)

func startListening(adapter *bluetooth.Adapter, sensors map[bluetooth.MAC]AtcSensor, config *Config, mqttAdaptor *mqtt.Adaptor) {
	// Enable BLE interface.
	err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		for mac, sensor := range sensors {
			if device.Address.String() != mac.String() {
				continue
			}
			change, failure := sensor.UpdateDevice(device.AdvertisementPayload)
			if failure != nil {
				println(failure.Error())
			}
			if change {
				jsonBytes, err := json.Marshal(sensor.Packet())
				if err != nil {
					fmt.Printf("error marshalling packet: %s", err)
					continue
				}

				topic := fmt.Sprintf(config.MQTT.Path, sensor.Name()) // TODO: Move into sensor?
				fmt.Printf("Publishing to topic %s: %s\n", topic, string(jsonBytes))
				mqttAdaptor.Publish(topic, jsonBytes)
			}
		}
	})
	panicCheck("scan error", err)
}
