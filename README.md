[![Go](https://github.com/dzhambulat/AirPolMonitor/actions/workflows/go.yml/badge.svg)](https://github.com/dzhambulat/AirPolMonitor/actions/workflows/go.yml)

# Air Pollution Monitor

The **core** app generates sensor-style readings (mock source), forwards them through the pipeline, and **produces JSON messages to Kafka** (`types.AirData`: `SensorID`, `Value`, `Timestamp`).

The **places** service (`services/places`) consumes that Kafka topic, keeps an in-memory sensor store keyed by place, and exposes **HTTP APIs** to create and edit places.

## Prerequisites

- **Go** (see `go.mod`; currently **1.25.x**)
- **Docker** and **Docker Compose** (for local Kafka, MQTT, observability, and optional places container)

## Configuration (core)

The core app reads `config.json` from the **repository root** (working directory when you run it).

```json
{
    "mqtt_broker_url": "tcp://localhost:1883",
    "mqtt_client_id": "air-pol-monitor",
    "mqtt_topic": "air-pollution",
    "kafka_broker_url": "localhost:9092",
    "kafka_topic": "air-pollution-data"
}
```

Optional MQTT authentication environment variables:

- `MQTT_USERNAME`
- `MQTT_PASSWORD`

The core app also listens for **Prometheus metrics** on `:9091` by default (override with `METRICS_PORT`, for example `9091` to listen on `:9091`).

## Run with Docker Compose (recommended for local stack)

From the repository root:

```bash
docker compose up -d
```

This starts **Mosquitto**, **Kafka**, **Prometheus**, **Grafana**, and the **places** service. Kafka is available at `localhost:9092`. The places HTTP API is mapped to **port 8082**.

Wait until Kafka is healthy (the `places` service waits on Kafka’s healthcheck). Check containers:

```bash
docker compose ps
```

To stop:

```bash
docker compose down
```

## Run the core app

With Kafka reachable at the URL in `config.json` (for Compose, use `localhost:9092` on the host):

```bash
go run .
```

The process reads mock sensor data and writes to the configured Kafka topic (`kafka_topic`, default `air-pollution-data`).

## Run the places service

### Option A: same machine as Compose (host talks to Kafka in Docker)

```bash
cd services/places
export KAFKA_BOOTSTRAP_SERVERS=localhost:9092
export KAFKA_TOPIC=air-pollution-data
export KAFKA_GROUP_ID=places-service
export PLACES_HTTP_ADDR=:8082
go mod download
go run .
```

### Option B: places already running in Compose

After `docker compose up -d`, the **places** container runs `go run .` inside the repo mount. Use **http://localhost:8082** on the host.

Environment variables used by places (defaults in parentheses):

| Variable | Purpose |
|----------|---------|
| `KAFKA_BOOTSTRAP_SERVERS` | Kafka brokers (`localhost:9092`) |
| `KAFKA_TOPIC` | Topic to consume (`air-pollution-data`) |
| `KAFKA_GROUP_ID` | Consumer group (`places-service`) |
| `PLACES_HTTP_ADDR` | HTTP listen address (`:8082`) |
| `PLACES_DEFAULT_PLACE_ID` | Default place id for routing when needed (`default-place`) |

## Build (core)

From the repository root:

```bash
go build -o AirPolMonitor .
```

Run the binary:

```bash
./AirPolMonitor
```

## Test end-to-end (core → Kafka → places)

1. Start the stack: `docker compose up -d`.
2. In another terminal, from the **repo root**, run the core app so it produces to Kafka:
   ```bash
   go run .
   ```
3. Confirm the **places** container or local `go run .` in `services/places` is running and connected (check logs for the consumer / HTTP listen message).

Optional: publish sample MQTT payloads (useful if you extend the core to use MQTT instead of the mock source):

```bash
pip install -r scripts/requirements-mqtt-test.txt
python scripts/mqtt_publish_test.py --broker localhost --port 1883 --topic air-pollution
```

## Test places HTTP API

With places listening on `localhost:8082`:

**Create a place** (`POST /places`):

```bash
curl -sS -X POST "http://localhost:8082/places" \
  -H "Content-Type: application/json" \
  -d '{"id":"","name":"Lab","coordinates":[40.7128,-74.0060]}'
```

**Edit a place** (`PUT /places/{id}` — replace `{id}` with the place id you use):

```bash
curl -sS -X PUT "http://localhost:8082/places/{id}" \
  -H "Content-Type: application/json" \
  -d '{"name":"Lab updated","coordinates":[40.7306,-73.9352]}'
```

Successful responses are JSON with an `id` field (HTTP 201 on create, 200 on update). Errors return JSON `{"error":"..."}` with a 4xx/5xx status.

## Sending test data over MQTT (reference)

```bash
mosquitto_pub -h localhost -p 1883 -t air-pollution -m '{"sensor_id": "sensor-123", "pm25": 10.5, "pm10": 25.2}'
```

## Observability (Compose)

- **Prometheus**: http://localhost:9090  
- **Grafana**: http://localhost:3000 (default admin/admin in Compose)
