package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"tinygo.org/x/bluetooth"
)

func TestAtcSensor(t *testing.T) {
	mac, _ := bluetooth.ParseMAC("AA:BB:CC:11:22:34")

	t.Run("Test Name updated", func(t *testing.T) {
		update := &bluetooth.ScanResult{
			Address: bluetooth.MACAddress{},
			RSSI:    50,
			AdvertisementPayload: &bluetooth.InternalAdvertisementFields{
				AdvertisementFields: bluetooth.AdvertisementFields{
					LocalName:    "NEW_NAME",
					ServiceUUIDs: []bluetooth.UUID{},
					ServiceData:  map[string][]byte{
						"0000181a-0000-1000-8000-00805f9b34fb": {0xaa, 0xbb, 0xcc, 0x11, 0x22, 0x34, 0x01, 0x18, 0x29, 0x59, 0x0b, 0xc2, 0x13},
					},
				},
			},
		}
		sensor := NewATCSensor(mac)

		response, err := sensor.UpdateDevice(update)
		actual := sensor.Name()

		assert.Equal(t, true, response)
		assert.Equal(t, "NEW_NAME", actual)
		assert.Nil(t, err)
	})

	t.Run("Test data output", func(t *testing.T) {
		update := &bluetooth.ScanResult{
			Address: bluetooth.MACAddress{},
			RSSI:    50,
			AdvertisementPayload: &bluetooth.InternalAdvertisementFields{
				AdvertisementFields: bluetooth.AdvertisementFields{
					LocalName:    "NEW_NAME",
					ServiceUUIDs: []bluetooth.UUID{},
					ServiceData:  map[string][]byte{
						"0000181a-0000-1000-8000-00805f9b34fb": {0xaa, 0xbb, 0xcc, 0x11, 0x22, 0x34, 0x01, 0x14, 0x29, 0x59, 0x0b, 0xc2, 0x13},
					},
				},
			},
		}
		expected := AtcPacket{
			Temperature: 27.6,
			Humidity:    41,
			Battery:     89,
		}
		sensor := NewATCSensor(mac)
		sensor.UpdateDevice(update)

		actual := sensor.Packet()

		assert.Equal(t, expected, actual)
	})

	t.Run("Test data output", func(t *testing.T) {
		update := &bluetooth.ScanResult{
			Address: bluetooth.MACAddress{},
			RSSI:    50,
			AdvertisementPayload: &bluetooth.InternalAdvertisementFields{
				AdvertisementFields: bluetooth.AdvertisementFields{
					LocalName:    "NEW_NAME",
					ServiceUUIDs: []bluetooth.UUID{},
					ServiceData:  map[string][]byte{
						"0000181a-0000-1000-8000-00805f9b34fb": {0xaa, 0xbb, 0xcc, 0x11, 0x22, 0x34, 0x01, 0x14, 0x29, 0x59, 0x0b, 0xc2, 0x13},
					},
				},
			},
		}
		expectedPacket := AtcPacket{
			Temperature: 27.6,
			Humidity:    41,
			Battery:     89,
		}
		sensor := NewATCSensor(mac)
		response1, err1 := sensor.UpdateDevice(update)
		response2, err2 := sensor.UpdateDevice(update)

		actual := sensor.Packet()
		assert.Equal(t, true, response1)
		assert.Equal(t, false, response2)
		assert.Nil(t, err1)
		assert.Nil(t, err2)
		assert.Equal(t, expectedPacket, actual)
	})
}
