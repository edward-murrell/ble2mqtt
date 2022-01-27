package main

import (
	"encoding/json"
	"fmt"
	"gobot.io/x/gobot/platforms/mqtt"
	"log"
	"strings"
	"tinygo.org/x/bluetooth"
)


func main() {
	config, conErr := getConfig()
	panicCheck("loading config", conErr)

	var adapter = bluetooth.DefaultAdapter // TODO, allow other than hci0
	panicCheck("enable BLE stack", adapter.Enable())

	mqttAdaptor := getMqttConnection(config)
	sensors := getSensors(config)

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

func getSensors(config *Config) map[bluetooth.MAC]AtcSensor {
	if len(config.Sensors) == 0 {
		log.Fatalf("no configured sensors found in configuration file")
	}

	sensors := make(map[bluetooth.MAC]AtcSensor, len(config.Sensors))

	for idx, sensorCfg := range config.Sensors {
		rawMac := strings.Trim(sensorCfg.MAC, " ")
		mac, parseE := bluetooth.ParseMAC(rawMac)
		if parseE != nil {
			log.Fatalf("fatal error on sensorCfg %d, %s %s", idx+2, parseE.Error(), sensorCfg)
		}
		sensors[mac] = *NewATCSensor(mac)
	}

	return sensors
}

func getMqttConnection(config *Config) *mqtt.Adaptor {
	address := fmt.Sprintf("tcp://%s:%d", config.MQTT.Host, config.MQTT.Port)
	mqttAdaptor := mqtt.NewAdaptor(address, "ble2mqtt")
	mqttError := mqttAdaptor.Connect()
	panicCheck("MQTT Connect", mqttError)
	return mqttAdaptor
}

func panicCheck(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
