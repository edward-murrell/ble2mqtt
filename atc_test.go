package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"tinygo.org/x/bluetooth"
)

func TestAtcSensor(t *testing.T) {
	mac, _ := bluetooth.ParseMAC("AA:BB:CC:11:22:34")

	t.Run("Test Name updated", func(t *testing.T) {
		update := fakeScanResult{
			Address:      bluetooth.MACAddress{},
			RSSI:         50,
			localname:    "NEW_NAME",
			ServiceUUIDs: []bluetooth.UUID{},
			ServiceData:  map[string][]byte{
				"0000181a-0000-1000-8000-00805f9b34fb": {0xaa, 0xbb, 0xcc, 0x11, 0x22, 0x34, 0x01, 0x18, 0x29, 0x59, 0x0b, 0xc2, 0x13},
			},
		}
		sensor := NewATCSensor(mac)

		response, err := sensor.UpdateDevice(&update)
		actual := sensor.Name()

		assert.Equal(t, true, response)
		assert.Equal(t, "NEW_NAME", actual)
		assert.Nil(t, err)
	})
}

// bluetooth.ScanResult has an unexported type (advertisementFields) which makes testing difficult
type fakeScanResult struct {
	Address bluetooth.Addresser
	RSSI int16
	localname string
	ServiceUUIDs []bluetooth.UUID
	ServiceData map[string][]byte
}

func (f *fakeScanResult) LocalName() string {
	return f.localname
}
func (f *fakeScanResult) HasServiceUUID(uuid bluetooth.UUID) bool {
	for _, u := range f.ServiceUUIDs {
		if u == uuid {
			return true
		}
	}
	return false
}

// Return the service data.
func (f *fakeScanResult) GetServiceData(uuid bluetooth.UUID) ([]byte, error) {
	data, ok := f.ServiceData[uuid.String()]
	if ok {
		return data, nil
	}
	return nil, errors.New("service key does not exist")
}

func (f *fakeScanResult) Bytes() []byte {
	return nil
}