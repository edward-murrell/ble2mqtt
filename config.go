package main

import "github.com/asppj/goconfig"

type Config struct {
	Controller string `short:"c" default:"hci0" desc:"Controller to use for listening to BLE packets"`
	Sensors    []SensorConfig
	MQTT       MqttConfig `id:"mqtt"`
	Logging    LoggingConfig `id:"logs"`
	Config     string `id:"config"`
}

type SensorConfig struct {
	MAC  string
	Type string `default:"atc"`
}

type LoggingConfig struct {
	File      string `default:"stderr"`
	Timestamp bool   `default:"true"`
	Level     string `short:"l" default:"debug"`
}

type MqttConfig struct {
	Host string `default:"localhost"`
	Port int16  `default:"1883"`
	Path string `default:"sensor/%s/state"`
}

func getConfig() (*Config, error) {
	config := &Config{}
	err := goconfig.Load(config, goconfig.Conf{
		ConfigFileVariable:  "config",
		FileDefaultFilename: "ble2mqtt.conf",
		FileDecoder:         goconfig.DecoderYAML,
		EnvPrefix:           "BLE2MQTT_",
	})

	return config, err
}
