#!/usr/bin/env python3
"""
Publish sample air-pollution data to MQTT for testing the app → Kafka pipeline.
Usage:
  pip install paho-mqtt
  python scripts/mqtt_publish_test.py [--broker localhost] [--port 1883] [--topic air-pollution] [--count 5] [--interval 2]
"""

import argparse
import json
import time
import random

try:
    import paho.mqtt.client as mqtt
except ImportError:
    print("Install with: pip install paho-mqtt")
    raise

DEFAULT_BROKER = "localhost"
DEFAULT_PORT = 1883
DEFAULT_TOPIC = "air-pollution"
SENSOR_IDS = ("sensor-01", "sensor-02", "sensor-03")
POLLUTANTS = ("pm25", "pm10", "no2", "o3")


def on_connect(client, userdata, flags, reason_code, properties=None):
    if reason_code != 0:
        print(f"Failed to connect: {reason_code}")
    else:
        print(f"Connected to {userdata.get('broker', 'broker')}:{userdata.get('port', 1883)}")


def on_publish(client, userdata, mid, reason_code, properties=None):
    if reason_code != 0:
        print(f"Publish failed mid={mid} reason_code={reason_code}")


def make_payload(sensor_id: str, value: float, timestamp: int, pollutant: str) -> str:
    return json.dumps({
        "sensor_id": sensor_id,
        "value": round(value, 2),
        "timestamp": timestamp,
        "pollutant": pollutant,
    })


def main():
    p = argparse.ArgumentParser(description="Send test air-pollution data via MQTT")
    p.add_argument("--broker", default=DEFAULT_BROKER, help="MQTT broker host")
    p.add_argument("--port", type=int, default=DEFAULT_PORT, help="MQTT broker port")
    p.add_argument("--topic", default=DEFAULT_TOPIC, help="MQTT topic to publish to")
    p.add_argument("--count", type=int, default=5, help="Number of messages to send (0 = infinite)")
    p.add_argument("--interval", type=float, default=2.0, help="Seconds between messages")
    p.add_argument("--username", default=None, help="MQTT username (if required)")
    p.add_argument("--password", default=None, help="MQTT password (if required)")
    args = p.parse_args()

    client = mqtt.Client(
        callback_api_version=mqtt.CallbackAPIVersion.VERSION2,
        client_id="mqtt-test-publisher",
    )
    client.user_data_set({"broker": args.broker, "port": args.port})
    client.on_connect = on_connect
    client.on_publish = on_publish

    if args.username:
        client.username_pw_set(args.username, args.password)

    try:
        client.connect(args.broker, args.port, 60)
    except Exception as e:
        print(f"Cannot connect to {args.broker}:{args.port}: {e}")
        print("Ensure MQTT broker is running (e.g. Mosquitto).")
        return 1

    client.loop_start()
    sent = 0

    try:
        while True:
            sensor_id = random.choice(SENSOR_IDS)
            pollutant = random.choice(POLLUTANTS)
            value = round(random.uniform(5.0, 150.0), 2)
            timestamp = int(time.time())
            payload = make_payload(sensor_id, value, timestamp, pollutant)

            result = client.publish(args.topic, payload, qos=0)
            result.wait_for_publish()
            sent += 1
            print(f"[{sent}] {args.topic} <- {payload}")

            if args.count > 0 and sent >= args.count:
                break
            time.sleep(args.interval)
    except KeyboardInterrupt:
        print("\nStopped by user")
    finally:
        client.loop_stop()
        client.disconnect()

    print(f"Sent {sent} message(s).")
    return 0


if __name__ == "__main__":
    exit(main())
