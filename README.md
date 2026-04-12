[![Go](https://github.com/dzhambulat/AirPolMonitor/actions/workflows/go.yml/badge.svg)](https://github.com/dzhambulat/AirPolMonitor/actions/workflows/go.yml)

# Air Pollution Monitor

This application connects to an MQTT broker to receive air pollution data and sends it to a Kafka topic.

## Prerequisites

- Go (version 1.18 or later)
- MQTT Broker (e.g., Mosquitto)
- Kafka Broker (e.g., Apache Kafka)

## Configuration

The application is configured using the `config.json` file.

```json
{
    "mqtt_broker_url": "tcp://localhost:1883",
    "mqtt_client_id": "air-pol-monitor",
    "mqtt_topic": "air-pollution",
    "kafka_broker_url": "localhost:9092",
    "kafka_topic": "air-pollution-data"
}
```

You can also set the following environment variables for MQTT authentication:

- `MQTT_USERNAME`
- `MQTT_PASSWORD`

## Build

To build the application, run the following command:

```bash
go build
```

## Run

To run the application, use the following command:

```bash
./AirPolMonitor
```

The application will connect to the MQTT and Kafka brokers and start listening for messages.

## Sending Test Data

You can use an MQTT client like `mosquitto_pub` to send test data to the `air-pollution` topic.

```bash
mosquitto_pub -h localhost -p 1883 -t air-pollution -m '{"sensor_id": "sensor-123", "pm25": 10.5, "pm10": 25.2}'
```
