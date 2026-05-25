package bootstrap

import (
	"ble2mqtt/internal/adapt"
	"ble2mqtt/internal/bt"
	"ble2mqtt/internal/config"
	"ble2mqtt/internal/runner"
	"errors"
	"fmt"
	"gobot.io/x/gobot/platforms/mqtt"
	"log/slog"
	"strings"
	"tinygo.org/x/bluetooth"
)

func MakeTasks(cfg *config.Config, logger *slog.Logger) (*runner.RunMap, error) {
	// Move all this into a bootstrap
	adapter, contErr := MakeController(cfg)
	if contErr != nil {
		return nil, contErr
	}

	mqttAdaptor, mqttError := getMqttConnection(cfg)
	if mqttError != nil {
		return nil, mqttError
	}

	sensors, senErr := getSensors(cfg)
	if senErr != nil {
		return nil, senErr
	}

	loop := bt.NewAppLoop(cfg, logger, adapter, sensors, mqttAdaptor)

	return &runner.RunMap{
		"loop": loop.Run,
	}, nil
}

func MakeController(config *config.Config) (*bluetooth.Adapter, error) {
	var controller = bluetooth.DefaultAdapter // TODO, allow other than hci0
	conErr := controller.Enable()
	return controller, conErr
}

func getSensors(config *config.Config) (*bt.SensorStack, error) {
	if len(config.Sensors) == 0 {
		return nil, errors.New("no configured sensors found in configuration file")
	}

	sensors := make(bt.SensorStack, len(config.Sensors))

	for idx, sensorCfg := range config.Sensors {
		rawMac := strings.Trim(sensorCfg.MAC, " ")
		mac, parseE := bluetooth.ParseMAC(rawMac)
		if parseE != nil {
			return nil, fmt.Errorf("fatal error on sensorCfg %d, %s %s", idx+2, parseE.Error(), sensorCfg)
		}
		// Parsing ensures that MAC formats are identical.
		sensors[mac.String()] = adapt.NewATCSensor(mac)
	}

	return &sensors, nil
}

func getMqttConnection(config *config.Config) (*mqtt.Adaptor, error) {
	address := fmt.Sprintf("tcp://%s:%d", config.MQTT.Host, config.MQTT.Port)
	mqttAdaptor := mqtt.NewAdaptor(address, "ble2mqtt")
	mqttAdaptor.SetAutoReconnect(true)
	mqttError := mqttAdaptor.Connect()
	return mqttAdaptor, mqttError
}

func panicCheck(action string, err error) {
	if err != nil {
		panic("Failed while " + action + ": " + err.Error())
	}
}
