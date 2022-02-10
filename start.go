package main

import (
	"errors"
	"fmt"
	"gobot.io/x/gobot/platforms/mqtt"
	"strings"
	"tinygo.org/x/bluetooth"
)

func main() {
	config, confErr := getConfig()
	panicCheck("loading configuration", confErr)

	logger := getLogger(config)

	adapter, contErr := getController(config)
	panicCheck("enabling BLE controller", contErr)

	mqttAdaptor, mqttError := getMqttConnection(config)
	panicCheck("creating MQTT connection", mqttError)

	sensors, senErr := getSensors(config)
	panicCheck("loading sensors", senErr)

	startListening(logger, adapter, sensors, config, mqttAdaptor)
}



func getController(config *Config) (*bluetooth.Adapter, error) {
	var controller = bluetooth.DefaultAdapter // TODO, allow other than hci0
	conErr := controller.Enable()
	return controller, conErr
}

func getSensors(config *Config) (*sensorStack, error) {
	if len(config.Sensors) == 0 {
		return nil, errors.New("no configured sensors found in configuration file")
	}

	sensors := make(sensorStack, len(config.Sensors))

	for idx, sensorCfg := range config.Sensors {
		rawMac := strings.Trim(sensorCfg.MAC, " ")
		mac, parseE := bluetooth.ParseMAC(rawMac)
		if parseE != nil {
			return nil, fmt.Errorf("fatal error on sensorCfg %d, %s %s", idx+2, parseE.Error(), sensorCfg)
		}
		// Parsing ensures that MAC formats are identical.
		sensors[mac.String()] = NewATCSensor(mac)
	}

	return &sensors, nil
}

func getMqttConnection(config *Config) (*mqtt.Adaptor, error) {
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
