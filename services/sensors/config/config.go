package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	MqttBrokerURL  string `json:"mqtt_broker_url"`
	MqttClientID   string `json:"mqtt_client_id"`
	MqttTopic      string `json:"mqtt_topic"`
	KafkaBrokerURL string `json:"kafka_broker_url"`
	KafkaTopic     string `json:"kafka_topic"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
