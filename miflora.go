package main

import (
	"bytes"
	"fmt"
	"tinygo.org/x/bluetooth"
)

//         bytes   0-1: temperature in 0.1 °C
//        byte      2: unknown
//        bytes   3-6: brightness in Lux (MiFlora only)
//        byte      7: moisture in %
//        byted   8-9: conductivity in µS/cm
//        bytes 10-15: unknown

// Stolen from https://github.com/basnijholt/miflora/blob/master/miflora/miflora_poller.py

const MiFlUUID = "0000fe95-0000-1000-8000-00805f9b34fb"
const MiFlTempByte0 = 0
const MiFlTempByte1 = 1
const MiFlLuxByte0 = 3
const MiFlLuxByte1 = 4
const MiFlLuxByte2 = 5
const MiFlLuxByte3 = 6
const MiFlMoisByte0 = 7
const MiFlCondByte0 = 8
const MiFlCondByte1 = 9

type MiFloraSensor struct {
	name    string
	mac     bluetooth.MAC
	uuid    bluetooth.UUID
	data    []byte
	temp    float32
	lux     float32
	moist   float32
	conduct float32
	batt    float32
}

type MiFloraPacket struct {
	Temperature  float32 `json:"temperature"`
	Brightness   float32 `json:"brightness"`
	Moisture     float32 `json:"moisture"`
	Conductivity float32 `json:"conductivity"`
	Battery      float32 `json:"battery"`
}

// NewBatteryDriver creates a MiFloraSensor
func NewMiFloraSensor(mac bluetooth.MAC) *MiFloraSensor {
	uuid, _ := bluetooth.ParseUUID(MiFlUUID)
	n := &MiFloraSensor{
		name: "UNKNOWN",
		uuid: uuid,
		data: nil,
		mac:  mac,
	}
	return n
}

func (b *MiFloraSensor) getTemperature() float32 {
	decimal := (uint16(b.data[MiFlTempByte0]) * 256) + uint16(b.data[MiFlTempByte1])
	return float32(decimal - 5600) / 1000
}

func (b *MiFloraSensor) getBrightness() float32 {
	decimal := (uint16(b.data[MiFlLuxByte3]) * 2^24) + (uint16(b.data[MiFlLuxByte2]) * 2^16) + (uint16(b.data[MiFlLuxByte1]) * 2^8) + uint16(b.data[MiFlLuxByte0])
	return float32(decimal) / 10
}

func (b *MiFloraSensor) getMoisture() (value float32) {
	return float32(b.data[MiFlMoisByte0] - 100)
}

func (b *MiFloraSensor) getConductivity() float32 {
	decimal := (uint16(b.data[MiFlCondByte0]) * 256) + uint16(b.data[MiFlCondByte1])
	return float32(decimal) / 10
}

func (b *MiFloraSensor) getBattery() (value float32) {
	return float32(0)
}

// Will return UNKNOWN if name is not known.
func (b *MiFloraSensor) Name() string {
	return b.name
}

// Update device with whatever data was recieved by from the BLE sensor.
func (b *MiFloraSensor) UpdateDevice(update *bluetooth.ScanResult) (change bool, failure error) {
	if b.name == "UNKNOWN" && update.LocalName() != "" {
		b.name = update.LocalName()
		change = true
	}

	data, pErr := update.GetServiceData(b.uuid)
	if pErr != nil {
		failure = fmt.Errorf("service data for UUID %s not found in scan packet", MiFlUUID)
		change = false
		return
	}

	if len(data) < 10 {
		failure = fmt.Errorf("service data in UUID %s is too short", MiFlUUID)
		change = false
		return
	}

	if len(b.data) < 10 {
		b.data = data
	} else if bytes.Compare(b.data[MiFlTempByte0:MiFlCondByte1], data[MiFlTempByte0:MiFlCondByte1]) != 0 { // ugh
		b.data = data
	} else {
		return
	}

	temp := b.getTemperature()
	if temp != b.temp {
		b.temp = temp
		change = true
	}
	lux := b.getBrightness()
	if lux != b.lux {
		b.lux = lux
		change = true
	}
	moist := b.getMoisture()
	if moist != b.moist {
		b.moist = moist
		change = true
	}
	conduct := b.getConductivity()
	if conduct != b.conduct {
		b.conduct = conduct
		change = true
	}
	batt := b.getBattery()
	if batt != b.batt {
		b.batt = batt
		change = true
	}

	return
}

func (b *MiFloraSensor) Packet() MiFloraPacket {
	// If GetPacket is called before data is ready, then don't crash.
	if len(b.data) < 13 || b.name == "UNKNOWN" {
		return MiFloraPacket{}
	}

	return MiFloraPacket{
		Temperature:  b.temp,
		Brightness:   b.lux,
		Moisture:     b.moist,
		Conductivity: b.conduct,
		Battery:      b.batt,
	}
}
