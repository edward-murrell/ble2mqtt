BLE 2 MQTT adapter

For use with home assistant and Xiaomi Mijia BT sensors.

Home Assistant Configuration file:
```yaml
sensor:
  - platform: mqtt
    name: "Temperature"
    state_topic: "home-assistant/sensor/ABC/state"
    unit_of_measurement: "°C"
    value_template: "{{ value_json.temperature }}"
  - platform: mqtt
    name: "Humidity"
    state_topic: "home-assistant/sensor/ABC/state"
    unit_of_measurement: "%"
    value_template: "{{ value_json.humidity }}
```