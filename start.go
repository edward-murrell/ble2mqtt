package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot"
)


func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	envAdapt := NewEnvironmentSensor(bleAdaptor)


	work := func() {
		gobot.Every(5*time.Second, func() {
			fmt.Printf("Temperature level: %.1f C\n", envAdapt.GetTemperature())
			fmt.Printf("Humidity level: %.1f%%\n", envAdapt.GetHumidity())
		})
	}

	robot := gobot.NewRobot("bleBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{envAdapt},
		work,
	)

	robot.Start()
}

