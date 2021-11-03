package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/mqtt"
)


func main() {
	mqttUrl := os.Args[1]
	mqttAdaptor := mqtt.NewAdaptor(mqttUrl, "ble2mqtt")
	mqttError := mqttAdaptor.Connect()

	fmt.Printf("Starting connection to MQTT server: %s... ", mqttUrl)
	if mqttError != nil {
		fmt.Printf("failed - %s\n", mqttError)
		os.Exit(1)
	} else {
		fmt.Printf("success.\n")
	}

	bleAdaptor := ble.NewClientAdaptor(os.Args[2])
	envAdapt := NewEnvironmentSensor(bleAdaptor)
	bleAdaptor.Connect()
	envAdapt.SetName(fmt.Sprintf("%s_%s", envAdapt.Connection().Name(), strings.ReplaceAll(bleAdaptor.Address(), "_", "")))

	work := func() {
		gobot.Every(5*time.Second, func() {
			fmt.Printf("Querying %s\n", envAdapt.Name())
			packet := envAdapt.GetPacket()

			fmt.Printf("Temperature level: %.1f C\n", packet.Temperature)
			fmt.Printf("Humidity level: %.1f%%\n", packet.Humidity)

			jsonBytes, err := json.Marshal(packet)
			if err != nil {
				fmt.Printf("error marshalling packet: %s", err)
				return
			}
			topic := fmt.Sprintf("sensor/%s/state", envAdapt.Name())
			fmt.Printf("Publishing to topic %s: %s\n", topic, string(jsonBytes))

			mqttAdaptor.Publish("asd", jsonBytes)
		})
	}

	robot := gobot.NewRobot("bleBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{envAdapt},
		work,
	)

	robot.Start()
}
