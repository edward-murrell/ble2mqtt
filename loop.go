package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"tinygo.org/x/bluetooth"
)

type sensorStack map[string]AtcSensor

type appLoop struct {
	config      *Config
	logger      log.FieldLogger
	sensors     sensorStack
	mqttAdaptor mqttClient
}

func startListening(logger *log.Logger, adapter *bluetooth.Adapter, sensors sensorStack, config *Config, mqttAdaptor mqttClient) {
	loop := &appLoop{
		config:      config,
		logger:      logger,
		sensors:     sensors,
		mqttAdaptor: mqttAdaptor,
	}

	err := adapter.Scan(loop.handlePacket)
	panicCheck("scan error", err)
}

func (loop *appLoop) handlePacket(adapter *bluetooth.Adapter, blePacket bluetooth.ScanResult) {
	mac := blePacket.Address.String()
	sensor, ok := loop.sensors[mac]
	if !ok {
		return
	}
	change, failure := sensor.UpdateDevice(&blePacket)
	if failure != nil {
		loop.logger.Errorf("error updating sensor %s: %s", mac, failure)
	}
	if change {
		jsonBytes, err := json.Marshal(sensor.Packet())
		if err != nil {
			loop.logger.Errorf("error marshalling packet: %s", err)
			return
		}

		topic := fmt.Sprintf(loop.config.MQTT.Path, sensor.Name()) // TODO: Move into sensor?
		success := loop.mqttAdaptor.Publish(topic, jsonBytes)
		if !success {
			loop.logger.Errorf("Failed to publish to topic %s, data %s.", topic, string(jsonBytes))
		} else {
			loop.logger.Infof("Published to topic %s, data %s", topic, string(jsonBytes))
		}
	}
}

type mqttClient interface {
	Publish(topic string, message []byte) bool
}
