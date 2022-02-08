package main

import (
	"encoding/json"
	"fmt"
	"tinygo.org/x/bluetooth"
)

type sensorStack map[bluetooth.MAC]AtcSensor

type appLoop struct {
	config      *Config
	sensors     sensorStack
	mqttAdaptor mqttClient
}

func startListening(adapter *bluetooth.Adapter, sensors sensorStack, config *Config, mqttAdaptor mqttClient) {
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
			fmt.Printf("Publishing to topic %s: %s...", topic, string(jsonBytes))
			success := loop.mqttAdaptor.Publish(topic, jsonBytes)
			if !success {
				fmt.Printf("failed\n")
			} else {
				fmt.Printf("success\n")
			}

		}
	}
}

type mqttClient interface {
	Publish(topic string, message []byte) bool
}
