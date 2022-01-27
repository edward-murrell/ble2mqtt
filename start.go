package main

import (
	"encoding/json"
	"fmt"
	"gobot.io/x/gobot/platforms/mqtt"
	"log"
	"os"
	"strings"
	"tinygo.org/x/bluetooth"
)


func main() {
	var adapter = bluetooth.DefaultAdapter // TODO, allow other than hci0
	panicCheck("enable BLE stack", adapter.Enable())

	mqttAdaptor := getMqttConnection(os.Args[1])

	sensors := getSensors(os.Args)

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

				topic := fmt.Sprintf("sensor/%s/state", sensor.Name())
				fmt.Printf("Publishing to topic %s: %s\n", topic, string(jsonBytes))
				mqttAdaptor.Publish(topic, jsonBytes)
			}
		}
	})
	panicCheck("scan error", err)
}

func getSensors(args []string) map[bluetooth.MAC]AtcSensor {
	if len(args) < 3 {
		log.Fatal("No MAC address(es) specified")
	}

	sensors := make(map[bluetooth.MAC]AtcSensor, len(args)-2)

	for idx, arg := range args[2:] {
		arg = strings.Trim(arg, " ")
		mac, parseE := bluetooth.ParseMAC(arg)
		if parseE != nil {
			log.Fatalf("fatal error on arg %d, %s %s", idx+2, parseE.Error(), arg)
		}
		sensors[mac] = *NewATCSensor(mac)
	}

	return sensors
}

func getMqttConnection(address string) *mqtt.Adaptor {
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
