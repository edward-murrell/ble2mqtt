package main

import (
	"encoding/json"
	"fmt"
	"gobot.io/x/gobot/platforms/mqtt"
	"tinygo.org/x/bluetooth"
)

type appLoop struct {
	config *Config
	sensors map[bluetooth.MAC]AtcSensor
	mqttAdaptor *mqtt.Adaptor
}

func startListening(adapter *bluetooth.Adapter, sensors map[bluetooth.MAC]AtcSensor, config *Config, mqttAdaptor *mqtt.Adaptor) {
	loop := &appLoop{
		config:      config,
		sensors:     sensors,
		mqttAdaptor: mqttAdaptor,
	}

	err := adapter.Scan(loop.handlePacket)
	panicCheck("scan error", err)
}

func (loop *appLoop) handlePacket(adapter *bluetooth.Adapter, blePacket bluetooth.ScanResult) {
	for mac, sensor := range loop.sensors {
		if blePacket.Address.String() != mac.String() {
			continue
		}
		change, failure := sensor.UpdateDevice(&blePacket)
		if failure != nil {
			println(failure.Error())
		}
		if change {
			jsonBytes, err := json.Marshal(sensor.Packet())
			if err != nil {
				fmt.Printf("error marshalling packet: %s", err)
				continue
			}

			topic := fmt.Sprintf(loop.config.MQTT.Path, sensor.Name()) // TODO: Move into sensor?
			fmt.Printf("Publishing to topic %s: %s\n", topic, string(jsonBytes))
			loop.mqttAdaptor.Publish(topic, jsonBytes)
		}
	}
}