BLE 2 MQTT adapter

For use with home assistant and Xiaomi Mijia BT sensors running the ATC firmwar.

# Running
The first argument is address of the MQTT server, followed by the MAC addresses of any ATC devices.

`ble2mqtt tcp://192.168.1.55:1883 A4:C1:38:DB:4D:3C A4:C1:38:A0:30:EC`

# Home Assistant Configuration file:

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
