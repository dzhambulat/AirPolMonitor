package mqtt

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Message is a received MQTT message.
type Message struct {
	Topic   string
	Payload []byte
}

// SubscriberConfig configures the MQTT subscriber.
type SubscriberConfig struct {
	BrokerURL string // e.g. "tcp://localhost:1883"
	ClientID  string
	Topic     string
	ShareName string // non-empty enables shared subscription: $share/ShareName/Topic
	Username  string // optional; often from MQTT_USERNAME env
	Password  string // optional; often from MQTT_PASSWORD env
	QoS       byte   // 0, 1, or 2
}

// Subscriber reads messages from MQTT, optionally using a shared subscription.
type Subscriber struct {
	cfg    SubscriberConfig
	client mqtt.Client
	mu     sync.Mutex
}

// NewSubscriber creates a subscriber (call Subscribe to connect and receive).
func NewSubscriber(cfg SubscriberConfig) *Subscriber {
	return &Subscriber{cfg: cfg}
}

// subscribeTopic returns the topic filter; with ShareName it becomes a shared subscription.
func (s *Subscriber) subscribeTopic() string {
	if s.cfg.ShareName != "" {
		return fmt.Sprintf("$share/%s/%s", s.cfg.ShareName, s.cfg.Topic)
	}
	return s.cfg.Topic
}

// Subscribe connects to the broker, subscribes (with shared sub if ShareName set), and
// returns a channel of messages until ctx is cancelled. Caller must consume or cancel
// to avoid blocking the client's delivery flow.
func (s *Subscriber) Subscribe(ctx context.Context) (<-chan Message, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(s.cfg.BrokerURL)
	opts.SetClientID(s.cfg.ClientID)
	if s.cfg.Username != "" {
		opts.SetUsername(s.cfg.Username)
		opts.SetPassword(s.cfg.Password)
	} else if u := os.Getenv("MQTT_USERNAME"); u != "" {
		opts.SetUsername(u)
		opts.SetPassword(os.Getenv("MQTT_PASSWORD"))
	}

	msgCh := make(chan Message, 32)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		select {
		case msgCh <- Message{Topic: msg.Topic(), Payload: msg.Payload()}:
		case <-ctx.Done():
		}
	})

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("mqtt: connection lost: %v", err)
	})

	s.mu.Lock()
	s.client = mqtt.NewClient(opts)
	s.mu.Unlock()

	token := s.client.Connect()
	if !token.WaitTimeout(connectTimeout) {
		return nil, fmt.Errorf("mqtt: connect timeout to %s", s.cfg.BrokerURL)
	}
	if err := token.Error(); err != nil {
		return nil, fmt.Errorf("mqtt: connect: %w", err)
	}

	topic := s.subscribeTopic()
	log.Printf("mqtt: subscribing (shared=%v) to %q", s.cfg.ShareName != "", topic)
	token = s.client.Subscribe(topic, s.cfg.QoS, nil)
	if !token.WaitTimeout(subscribeTimeout) {
		s.client.Disconnect(250)
		return nil, fmt.Errorf("mqtt: subscribe timeout")
	}
	if err := token.Error(); err != nil {
		s.client.Disconnect(250)
		return nil, fmt.Errorf("mqtt: subscribe: %w", err)
	}

	go func() {
		<-ctx.Done()
		s.mu.Lock()
		c := s.client
		s.mu.Unlock()
		if c != nil {
			c.Unsubscribe(topic)
			c.Disconnect(250)
		}
		close(msgCh)
	}()

	return msgCh, nil
}

// Disconnect disconnects the client. Safe to call after Subscribe's context is cancelled.
func (s *Subscriber) Disconnect() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.client != nil {
		s.client.Disconnect(250)
		s.client = nil
	}
}

const connectTimeout = 10 * time.Second
const subscribeTimeout = 5 * time.Second
