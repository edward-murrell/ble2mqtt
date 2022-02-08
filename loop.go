package main

import (
	"encoding/json"
	"fmt"
	"tinygo.org/x/bluetooth"
)

type sensorStack map[string]AtcSensor

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
	sensor, ok := loop.sensors[blePacket.Address.String()]
	if !ok {
		return
	}
	change, failure := sensor.UpdateDevice(&blePacket)
	if failure != nil {
		println(failure.Error())
	}
	if change {
		jsonBytes, err := json.Marshal(sensor.Packet())
		if err != nil {
			fmt.Printf("error marshalling packet: %s", err)
			return
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

type mqttClient interface {
	Publish(topic string, message []byte) bool
}
