BLE 2 MQTT adapter

For use with home assistant and Xiaomi Mijia BT sensors running the ATC firmware.

This uses Bluetooth controller scanning, and may interfere with active connections. If you have long-running active
 bluetooth connections, it's advisable to have more than one bluetooth adapter.

# 1. Configuration
Copy the example YAML configuration file and modify to suit your installation.

`cp ble2mqtt.conf.example ble2mqtt.conf`

A minimal configuration requires a single sensor input, and will connect to the MQTT server on localhost, port 1883 if
 no other configuration is specified. This default configuration will log error messages only to stderr.

```yaml
sensors:
  - mac: A4:C1:38:DB:12:34 # MAC address of the sensor
    type: atc              # Currently, only the type atc is supported.
```

# 2. Home Assistant Configuration file:

Replace ATC_FA1234 with the name of your sensor.

```yaml
sensor:
  - platform: mqtt
    name: "Temperature"
    state_topic: "sensor/ATC_FA1234/state"
    unit_of_measurement: "°C"
    value_template: "{{ value_json.temperature }}"
  - platform: mqtt
    name: "Humidity"
    state_topic: "sensor/ATC_FA1234/state"
    unit_of_measurement: "%"
    value_template: "{{ value_json.humidity }}"
```

# 3. Building
```
go get
go build
```

# 4. Running
`./ble2mqtt`
