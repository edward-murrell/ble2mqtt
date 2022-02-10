package main

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"tinygo.org/x/bluetooth"
)

func Test_scanLoop(t *testing.T) {
	adapter := &bluetooth.Adapter{}
	mqtt := &FakeMqtt{}
	logger, logBuffer := NewFakeLogger()

	mac := bluetooth.MACAddress{MAC: [6]byte{0x34, 0x22, 0x11, 0xcc, 0xbb, 0xaa}}

	app := &appLoop{
		config: &Config{
			Sensors: []SensorConfig{},
			MQTT: MqttConfig{
				Path: "sensor/%s/state",
			},
		},
		logger:      logger,
		sensors:     NewSensorStack(mac.String()),
		mqttAdaptor: mqtt,
	}

	t.Run("test single update", func(t *testing.T) {
		blePacket := createFakeAtcResult(mac, "TEST_SENSOR", []byte{0x34, 0x22, 0x11, 0xcc, 0xbb, 0xaa, 0x01, 0x18, 0x29, 0x59, 0x0b, 0xc2, 0x13})

		app.handlePacket(adapter, *blePacket)

		assert.Equal(t, &PubStore{"sensor/TEST_SENSOR/state", []byte(`{"temperature":28,"humidity":41,"battery":89}`)}, mqtt.Publishes[0])
		assert.Equal(t, `level=info msg="Published to topic sensor/TEST_SENSOR/state, data {\"temperature\":28,\"humidity\":41,\"battery\":89}"` + "\n", logBuffer.String())
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

func NewSensorStack(macs ...string) sensorStack {
	sensors := make(sensorStack)
	for _, mac := range macs {
		parsedMac, _ := bluetooth.ParseMAC(mac)
		sensors[mac] = *NewATCSensor(parsedMac)
	}
	return sensors
}

func NewFakeLogger() (*log.Logger, *bytes.Buffer) {
	output := new(bytes.Buffer)
	return &log.Logger{
		Out:          output,
		Formatter:    &log.TextFormatter{
			DisableTimestamp: true,
			FullTimestamp: false,
		},
		Hooks:        make(log.LevelHooks),
		Level:        log.DebugLevel,
		ExitFunc:     fakeExit,
		ReportCaller: false,
	}, output
}

func fakeExit(_ int) {}