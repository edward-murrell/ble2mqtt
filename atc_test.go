package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"tinygo.org/x/bluetooth"
)

func TestAtcSensor(t *testing.T) {
	mac, _ := bluetooth.ParseMAC("AA:BB:CC:11:22:34")

	t.Run("Test Name updated", func(t *testing.T) {
		update := createFakeAtcResult(bluetooth.MACAddress{}, "NEW_NAME", []byte{0xaa, 0xbb, 0xcc, 0x11, 0x22, 0x34, 0x01, 0x18, 0x29, 0x59, 0x0b, 0xc2, 0x13})
		sensor := NewATCSensor(mac)

		response, err := sensor.UpdateDevice(update)
		actual := sensor.Name()

		assert.Equal(t, true, response)
		assert.Equal(t, "NEW_NAME", actual)
		assert.Nil(t, err)
	})

	t.Run("Test data output from a sensor that has not seen any data before", func(t *testing.T) {
		update := createFakeAtcResult(bluetooth.MACAddress{}, "NEW_NAME", []byte{0xaa, 0xbb, 0xcc, 0x11, 0x22, 0x34, 0x01, 0x14, 0x29, 0x59, 0x0b, 0xc2, 0x13})
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

	t.Run("Test data output with a repeat packet does not generate a change notification", func(t *testing.T) {
		update := createFakeAtcResult(bluetooth.MACAddress{}, "NEW_NAME", []byte{0xaa, 0xbb, 0xcc, 0x11, 0x22, 0x34, 0x01, 0x14, 0x29, 0x59, 0x0b, 0xc2, 0x13})
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

	t.Run("Test data output with a same data input data, but ignored data not generate a change notification", func(t *testing.T) {
		update1 := createFakeAtcResult(bluetooth.MACAddress{}, "NEW_NAME", []byte{0xaa, 0xbb, 0xcc, 0x11, 0x22, 0x34, 0x02, 0x24, 0x39, 0x49, 0x11, 0xc2, 0x13})
		update2 := createFakeAtcResult(bluetooth.MACAddress{}, "NEW_NAME", []byte{0xaa, 0xbb, 0xcc, 0x11, 0x22, 0x34, 0x02, 0x24, 0x39, 0x49, 0xfa, 0xfb, 0xfc})
		expectedPacket := AtcPacket{
			Temperature: 54.8,
			Humidity:    57,
			Battery:     73,
		}
		sensor := NewATCSensor(mac)
		response1, err1 := sensor.UpdateDevice(update1)
		response2, err2 := sensor.UpdateDevice(update2)

		actual := sensor.Packet()
		assert.Equal(t, true, response1)
		assert.Equal(t, false, response2)
		assert.Nil(t, err1)
		assert.Nil(t, err2)
		assert.Equal(t, expectedPacket, actual)
	})
}

func createFakeAtcResult(mac bluetooth.MACAddress, name string, payload []byte) *bluetooth.ScanResult {
	return &bluetooth.ScanResult{
		Address: mac,
		RSSI:    50,
		AdvertisementPayload: &bluetooth.InternalAdvertisementFields{
			AdvertisementFields: bluetooth.AdvertisementFields{
				LocalName:    name,
				ServiceUUIDs: []bluetooth.UUID{},
				ServiceData: map[string][]byte{
					"0000181a-0000-1000-8000-00805f9b34fb": payload,
				},
			},
		},
	}
}
