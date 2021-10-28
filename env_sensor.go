package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
)

type EnvSensor struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

// NewBatteryDriver creates a EnvSensor
func NewEnvironmentSensor(a ble.BLEConnector) *EnvSensor {
	n := &EnvSensor{
		name:       gobot.DefaultName("Environment sensor"),
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	return n
}

// Connection returns the Driver's Connection to the associated Adaptor
func (b *EnvSensor) Connection() gobot.Connection { return b.connection }

// Name returns the Driver name
func (b *EnvSensor) Name() string { return b.name }

// SetName sets the Driver name
func (b *EnvSensor) SetName(n string) { b.name = n }

// adaptor returns BLE adaptor
func (b *EnvSensor) adaptor() ble.BLEConnector {
	return b.Connection().(ble.BLEConnector)
}

func (b *EnvSensor) Start() (err error) {
	return
}

func (b *EnvSensor) Halt() (err error) { return }

func (b *EnvSensor) GetTemperature() (value float32) {
	c, err := b.adaptor().ReadCharacteristic("00002a1f-0000-1000-8000-00805f9b34fb")
	if err != nil {
		fmt.Printf("error %s", err)
		return 0
	}
	value = float32(c[0]) / 10
	return
}
