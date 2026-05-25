package bt

import (
	"ble2mqtt/internal/adapt"
	"ble2mqtt/internal/config"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"log/slog"
	"tinygo.org/x/bluetooth"
)

type SensorStack map[string]*adapt.AtcSensor

type AppLoop struct {
	adaptor     *bluetooth.Adapter
	config      *config.Config
	logger      *slog.Logger
	sensors     *SensorStack
	mqttAdaptor mqttClient
}

func NewAppLoop(config *config.Config, logger *slog.Logger, adaptor *bluetooth.Adapter, sensors *SensorStack, mqttAdaptor mqttClient) *AppLoop {
	return &AppLoop{config: config, logger: logger, adaptor: adaptor, sensors: sensors, mqttAdaptor: mqttAdaptor}
}

func (loop *AppLoop) Run(ctx context.Context) error {
	var err error
	go func() {
		err = loop.adaptor.Scan(loop.handlePacket)
	}()
	select {
	case <-ctx.Done():
	}
	return err
}

func (loop *AppLoop) handlePacket(adapter *bluetooth.Adapter, blePacket bluetooth.ScanResult) {
	mac := blePacket.Address.String()
	sensors := *loop.sensors
	sensor, ok := sensors[mac]
	if !ok {
		return
	}
	change, failure := sensor.UpdateDevice(&blePacket)
	if failure != nil {
		loop.logger.Error("error updating sensor", slog.String("mac", mac), "failure", failure.Error())
	}
	if change {
		log.Debugf("change detected for %s", mac)
		jsonBytes, err := json.Marshal(sensor.Packet())
		if err != nil {
			loop.logger.Error("error marshalling packet", "error", err.Error())
			return
		}

		topic := fmt.Sprintf(loop.config.MQTT.Path, sensor.Name()) // TODO: Move into sensor?
		success := loop.mqttAdaptor.Publish(topic, jsonBytes)
		if !success {
			loop.logger.Error("failed to publish to topic", slog.String("topic", topic), slog.String("data", string(jsonBytes)))
		} else {
			loop.logger.Info("published to topic", slog.String("topic", topic), slog.String("data", string(jsonBytes)))
		}
	}
}

type mqttClient interface {
	Publish(topic string, message []byte) bool
}
