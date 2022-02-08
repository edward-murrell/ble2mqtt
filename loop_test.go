package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"tinygo.org/x/bluetooth"
)

func Test_scanLoop(t *testing.T) {
	adapter := &bluetooth.Adapter{}
	mqtt := &FakeMqtt{}
	app := &appLoop{
		config: &Config{
			Sensors: []SensorConfig{},
			MQTT: MqttConfig{
				Path: "sensor/%s/state",
			},
		},
		sensors: NewSensorStack("AA:BB:CC:11:22:34"),
		mqttAdaptor: mqtt,
	}

	t.Run("test single update", func(t *testing.T) {
		blePacket := bluetooth.ScanResult{
			Address: &bluetooth.Address{
				MACAddress: bluetooth.MACAddress{MAC: [6]byte{0x34, 0x22, 0x11, 0xcc, 0xbb, 0xaa}},
			},
			RSSI: 50,
			AdvertisementPayload: &bluetooth.InternalAdvertisementFields{
				AdvertisementFields: bluetooth.AdvertisementFields{
					LocalName:    "TEST_SENSOR",
					ServiceUUIDs: []bluetooth.UUID{},
					ServiceData: map[string][]byte{
						"0000181a-0000-1000-8000-00805f9b34fb": {0x34, 0x22, 0x11, 0xcc, 0xbb, 0xaa, 0x01, 0x18, 0x29, 0x59, 0x0b, 0xc2, 0x13},
					},
				},
			},
		}

		app.handlePacket(adapter, blePacket)

		assert.Equal(t, &PubStore{"sensor/TEST_SENSOR/state", []byte(`{"temperature":28,"humidity":41,"battery":89}`)}, mqtt.Publishes[0])
	})
}

type FakeMqtt struct {
	Publishes []*PubStore
}
type PubStore struct {
	Topic   string
	message []byte
}

func (m *FakeMqtt) Publish(topic string, message []byte) bool {
	m.Publishes = append(m.Publishes, &PubStore{topic, message})
	return true
}

func NewSensorStack(macs... string) sensorStack  {
	sensors := make(sensorStack)
	for _, mac := range macs {
		parsedMac, _ := bluetooth.ParseMAC(mac)
		sensors[parsedMac] = *NewATCSensor(parsedMac)
	}
	return sensors
}