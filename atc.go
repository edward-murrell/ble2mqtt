package main

import (
	"bytes"
	"fmt"
	"tinygo.org/x/bluetooth"
)

const UUID = "0000181a-0000-1000-8000-00805f9b34fb"

type AtcSensor struct {
	name string
	mac bluetooth.MAC
	uuid bluetooth.UUID
	data []byte
}

type AtcPacket struct {
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
	Battery     float32 `json:"battery"`
}

// NewBatteryDriver creates a AtcSensor
func NewATCSensor(mac bluetooth.MAC) *AtcSensor {

	uuid, _ := bluetooth.ParseUUID(UUID)
	n := &AtcSensor{
		name: "UNKNOWN",
		uuid: uuid,
		data: nil,
		mac: mac,
	}
	return n
}

func (b *AtcSensor) getTemperature() (float32) {
	decimal := (uint16(b.data[6]) * 256) + uint16(b.data[7])
	return float32(decimal) / 10
}

func (b *AtcSensor) getHumidity() (value float32) {
	return float32(b.data[8])
}

func (b *AtcSensor) getBattery() (value float32) {
	return float32(b.data[9])
}

// Will return UNKNOWN if name is not known.
func (b *AtcSensor) Name() string  {
	return b.name
}

// Update device with whatever data was recieved by from the BLE sensor.
func (b *AtcSensor) UpdateDevice(update bluetooth.AdvertisementPayload) (change bool, failure error) {
	if b.name == "UNKNOWN" && update.LocalName() != "" {
		b.name = update.LocalName()
		change = true
	}

	data, pErr := update.GetServiceData(b.uuid)
	if pErr != nil {
		failure = fmt.Errorf("recieved empty data for %s", UUID)
		change = false
		return
	}

	if len(data) < 13 {
		failure = fmt.Errorf("recieved service on UUID %s for ATC that it doesn't care about", UUID)
		change = false
		return
	}

	if bytes.Compare(b.data, data) != 0 {
		b.data = data
		change = true
	}

	return
}

func (b *AtcSensor) Packet() AtcPacket {
	// If GetPacket is called before data is ready, then don't crash.
	if len(b.data) < 13 || b.name == "UNKNOWN" {
		return AtcPacket{}
	}

	return AtcPacket{
		Temperature: b.getTemperature(),
		Humidity:    b.getHumidity(),
		Battery:     b.getBattery(),
	}
}
