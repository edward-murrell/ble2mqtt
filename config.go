package main

import "github.com/asppj/goconfig"

type Config struct {
	Controller string         `short:"c" default:"hci0" desc:"Controller to use for listening to BLE packets"`
	Sensors    []SensorConfig
	MQTT       MqttConfig     `id:"mqtt"`
}

type SensorConfig struct {
	MAC string
	Type string `default:"atc"`
}

type MqttConfig struct {
	Host string `default:"localhost"`
	Port int16  `default:"1883"`
	Path string `default:"sensor/%s/state"`
}

func getConfig() (*Config, error) {
	config := &Config{}
	err := goconfig.Load(config, goconfig.Conf{
		//ConfigFileVariable: "config", // enables passing --configfile myfile.conf

		FileDefaultFilename: "ble2mqtt.conf",
		// The default decoder will try TOML, YAML and JSON.
		FileDecoder: goconfig.DecoderYAML,

		EnvPrefix: "BLE2MQTT_",
	})
	return config, err
}


