# BLE to MQTT daemon

This daemon listens to BLE broadcasts from Xiaomi Mijia BT sensors running the [ATC](https://github.com/atc1441/ATC_MiThermometer) firmware 
 and sends the received data to an MQTT server. This is ideal for use with Home Assistant.

This program will put your Bluetooth controller into scan mode, and may interfere with active bluetooth connections.
If you have long-running active bluetooth connections, it's advisable to have more than one bluetooth adapter.

## 1. Configuration
Copy the example YAML configuration file and modify to suit your installation.

`cp ble2mqtt.conf.example ble2mqtt.conf`

A minimal configuration requires a single sensor input, and will connect to the MQTT server on localhost, port 1883 if
 no other configuration is specified. This default configuration will log error messages only to stderr.

```yaml
sensors:
  - mac: A4:C1:38:DB:12:34 # MAC address of the sensor
    type: atc              # Currently, only the type atc is supported.
```

## 2. Home Assistant Configuration file:

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

## 3. Installing

### Debian and Ubuntu package
The package installs a systemd service that reads `/etc/ble2mqtt/ble2mqtt.conf`.
The service stays stopped until that file exists.

```
sudo apt install ble2mqtt
sudo cp /usr/share/doc/ble2mqtt/examples/ble2mqtt.conf.example /etc/ble2mqtt/ble2mqtt.conf
sudo editor /etc/ble2mqtt/ble2mqtt.conf
sudo systemctl start ble2mqtt
```

### Docker
Images for amd64, arm64 and arm/v7 are on GHCR.
The container talks to BlueZ on the host, so mount the host's D-Bus socket.

```
docker run -d --name ble2mqtt --restart unless-stopped \
  -v /run/dbus:/run/dbus:ro \
  -v $PWD/ble2mqtt.conf:/etc/ble2mqtt/ble2mqtt.conf:ro \
  ghcr.io/edward-murrell/ble2mqtt:latest
```

## 4. Building
```
go build
./ble2mqtt
```

Set `BLE2MQTT_CONFIG` or pass `--config` to read a config file other than `./ble2mqtt.conf`.
