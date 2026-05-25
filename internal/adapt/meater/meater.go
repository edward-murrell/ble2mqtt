package meater

import (
	"bytes"
	"fmt"
	"tinygo.org/x/bluetooth"
)

const ServiceUUID = "a75cc7fc-c956-488f-ac2a-2dbc08b63a04"
const BatteryUUID = "2adb4877-68d8-4884-bd3c-d83853bf27b8"
const Firmware = "00002a26-0000-1000-8000-00805f9b34fb"
const TempByte0 = 6
const TempByte1 = 7
const HumiByte0 = 8
const BattByte0 = 9

type MeaterSensor struct {
	name    string
	mac     bluetooth.MAC
	svcUuid bluetooth.UUID
	data    []byte
	temp    float32
	humi    float32
	batt    float32
}

var (
	serviceUUID  bluetooth.UUID
	batteryUUID  bluetooth.UUID
	firmwareUUID bluetooth.UUID
)

func init() {
	serviceUUID, _ = bluetooth.ParseUUID(ServiceUUID)
	batteryUUID, _ := bluetooth.ParseUUID(BatteryUUID)
	firmwareUUID, _ := bluetooth.ParseUUID(Firmware)
}

type MeaterPacket struct {
	Tip     float32 `json:"tip"`
	Ambient float32 `json:"ambient"`
	Battery float32 `json:"battery"`
}

func NewMeaterSensor(mac bluetooth.MAC) *MeaterSensor {
	uuid, _ := bluetooth.ParseUUID(ServiceUUID)
	n := &MeaterSensor{
		name:    "UNKNOWN",
		svcUuid: uuid, // Why would you do that?
		data:    nil,
		mac:     mac,
	}
	return n
}

func (m *MeaterSensor) getTemperature() float32 {
	decimal := (uint16(m.data[TempByte0]) * 256) + uint16(m.data[TempByte1])
	return float32(decimal) / 10
}

func (m *MeaterSensor) getHumidity() (value float32) {
	return float32(m.data[HumiByte0])
}

func (m *MeaterSensor) getBattery() (value float32) {
	return float32(m.data[BattByte0])
}

// Will return UNKNOWN if name is not known.
func (m *MeaterSensor) Name() string {
	return m.name
}

// Update device with whatever data was recieved by from the BLE sensor.
func (m *MeaterSensor) UpdateDevice(update *bluetooth.ScanResult) (change bool, failure error) {
	if m.name == "UNKNOWN" && update.LocalName() != "" {
		m.name = update.LocalName()
		change = true
	}

	data := extractServiceData(update, serviceUUID)
	if data == nil {
		failure = fmt.Errorf("service data for UUID %s not found in scan packet", ServiceUUID)
		change = false
		return
	}

	if len(data) < 13 {
		failure = fmt.Errorf("service data in UUID %s is too short", ServiceUUID)
		change = false
		return
	}

	if len(m.data) < 13 {
		m.data = data
	} else if bytes.Compare(m.data[TempByte0:BattByte0], data[TempByte0:BattByte0]) != 0 {
		m.data = data
	} else {
		return
	}

	temp := m.getTemperature()
	if temp != m.temp {
		m.temp = temp
		change = true
	}
	humi := m.getHumidity()
	if humi != m.humi {
		m.humi = humi
		change = true
	}
	batt := m.getBattery()
	if batt != m.batt {
		m.batt = batt
		change = true
	}

	return
}

func extractServiceData(update *bluetooth.ScanResult, serviceUUID bluetooth.UUID) []byte {
	for _, elem := range update.ServiceData() {
		if elem.UUID == serviceUUID {
			return elem.Data
		}
		if elem.UUID == BatteryUUID {

		}
	}
	return nil
}

func (m *MeaterSensor) GetState() any {
	// If GetPacket is called before data is ready, then don't crash.
	if len(m.data) < 13 || m.name == "UNKNOWN" {
		return MeaterPacket{}
	}

	return MeaterPacket{
		Tip:     m.temp,
		Ambient: m.humi,
		Battery: m.batt,
	}
}
