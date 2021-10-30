package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/mqtt"
	"gobot.io/x/gobot"
)



func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	envAdapt := NewEnvironmentSensor(bleAdaptor)
	mqttAdaptor := mqtt.NewAdaptor("tcp://192.168.183.123:1883", "ble2mqtt")
	mqttAdaptor.Connect()
	mqttDriver := mqtt.NewDriver(mqttAdaptor, "home-assistant/sensor/ABC/state")

	work := func() {
		gobot.Every(5*time.Second, func() {
			packet := envAdapt.GetPacket()

			fmt.Printf("Temperature level: %.1f C\n", packet.Temperature)
			fmt.Printf("Humidity level: %.1f%%\n", packet.Humidity)

//			var jsonBytes []byte
			jsonBytes, err := json.Marshal(packet)
			if err != nil {
				fmt.Printf("Error: %s", err)
			}
			fmt.Printf("Send %s\n", string(jsonBytes))



			mqttDriver.Publish(jsonBytes)
		})
	}

	robot := gobot.NewRobot("bleBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{envAdapt},
		work,
	)

	robot.Start()
}

