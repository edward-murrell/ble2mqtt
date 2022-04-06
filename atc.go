package main

import (
	"bytes"
	"fmt"
	"tinygo.org/x/bluetooth"
)

const AtcUUID = "0000181a-0000-1000-8000-00805f9b34fb"
const AtcTempByte0 = 6
const AtcTempByte1 = 7
const AtcHumiByte0 = 8
const AtcBattByte0 = 9

type AtcSensor struct {
	name string
	mac  bluetooth.MAC
	uuid bluetooth.UUID
	data []byte
	temp float32
	humi float32
	batt float32
}

type AtcPacket struct {
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
	Battery     float32 `json:"battery"`
}

// NewBatteryDriver creates a AtcSensor
func NewATCSensor(mac bluetooth.MAC) *AtcSensor {

	uuid, _ := bluetooth.ParseUUID(AtcUUID)
	n := &AtcSensor{
		name: "UNKNOWN",
		uuid: uuid,
		data: nil,
		mac:  mac,
	}
	return n
}

func (b *AtcSensor) getTemperature() float32 {
	decimal := (uint16(b.data[AtcTempByte0]) * 256) + uint16(b.data[AtcTempByte1])
	return float32(decimal) / 10
}

func (b *AtcSensor) getHumidity() (value float32) {
	return float32(b.data[AtcHumiByte0])
}

func (b *AtcSensor) getBattery() (value float32) {
	return float32(b.data[AtcBattByte0])
}

// Will return UNKNOWN if name is not known.
func (b *AtcSensor) Name() string {
	return b.name
}

// Update device with whatever data was recieved by from the BLE sensor.
func (b *AtcSensor) UpdateDevice(update *bluetooth.ScanResult) (change bool, failure error) {
	if b.name == "UNKNOWN" && update.LocalName() != "" {
		b.name = update.LocalName()
		change = true
	}

	data, pErr := update.GetServiceData(b.uuid)
	if pErr != nil {
		failure = fmt.Errorf("service data for AtcUUID %s not found in scan packet", AtcUUID)
		change = false
		return
	}

	if len(data) < 13 {
		failure = fmt.Errorf("service data in AtcUUID %s is too short", AtcUUID)
		change = false
		return
	}

	if len(b.data) < 13 {
		b.data = data
	} else if bytes.Compare(b.data[AtcTempByte0:AtcBattByte0], data[AtcTempByte0:AtcBattByte0]) != 0 {
		b.data = data
	} else {
		return
	}

	temp := b.getTemperature()
	if temp != b.temp {
		b.temp = temp
		change = true
	}
	humi := b.getHumidity()
	if humi != b.humi {
		b.humi = humi
		change = true
	}
	batt := b.getBattery()
	if batt != b.batt {
		b.batt = batt
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
		Temperature: b.temp,
		Humidity:    b.humi,
		Battery:     b.batt,
	}
}
