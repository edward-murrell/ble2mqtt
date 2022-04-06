package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"tinygo.org/x/bluetooth"
)

func TestMiFloraSensor(t *testing.T) {
	mac, _ := bluetooth.ParseMAC("AA:BB:CC:11:22:34")

	t.Run("Test Name updated", func(t *testing.T) {
		update := createFakeMiFloraResult(bluetooth.MACAddress{}, "NEW_NAME", []byte{0x98, 0x00, 0xc7, 0x0b, 0x91, 0x6c, 0x8d, 0x7c, 0xc4, 0x0d, 0x08, 0x10, 0x01, 0x2f})
		sensor := NewMiFloraSensor(mac)

		response, err := sensor.UpdateDevice(update)
		actual := sensor.Name()

		assert.Equal(t, true, response)
		assert.Equal(t, "NEW_NAME", actual)
		assert.Nil(t, err)
	})

	t.Run("Test data output from a sensor that has not seen any data before", func(t *testing.T) {
		update := createFakeMiFloraResult(bluetooth.MACAddress{}, "NEW_NAME", []byte{0x98, 0x00, 0xc7, 0x0b, 0x91, 0x6c, 0x8d, 0x7c, 0xc4, 0x0d, 0x08, 0x10, 0x01, 0x2f})
		expected := MiFloraPacket{
			Temperature: 28.96,
			Brightness: 34.8,
			Moisture: 108, // 8
			Conductivity: 3622, // should be ~268
			Battery: 0,
		}

		sensor := NewMiFloraSensor(mac)
		sensor.UpdateDevice(update)

		actual := sensor.Packet()

		assert.Equal(t, expected, actual)
	})
	t.Run("Different data packet", func(t *testing.T) {
		update := createFakeMiFloraResult(bluetooth.MACAddress{}, "NEW_NAME1", []byte{0x71, 0x20, 0x98, 0x00, 0xd4, 0xa2, 0x9b, 0x6d, 0x8d, 0x7c, 0xc4, 0x0d, 0x08, 0x10, 0x01, 0x0b})
		expected := MiFloraPacket{
			Temperature: 23.3,
			Brightness: 650,
			Moisture: 9,
			Conductivity: 0,
			Battery: 100,
		}

		sensor := NewMiFloraSensor(mac)
		sensor.UpdateDevice(update)

		actual := sensor.Packet()

		assert.Equal(t, expected, actual)
	})
}

func createFakeMiFloraResult(mac bluetooth.MACAddress, name string, payload []byte) *bluetooth.ScanResult {
	return &bluetooth.ScanResult{
		Address: mac,
		RSSI:    50,
		AdvertisementPayload: &bluetooth.InternalAdvertisementFields{
			AdvertisementFields: bluetooth.AdvertisementFields{
				LocalName:    name,
				ServiceUUIDs: []bluetooth.UUID{},
				ServiceData: map[string][]byte{
					"0000fe95-0000-1000-8000-00805f9b34fb": payload,
				},
			},
		},
	}
}
