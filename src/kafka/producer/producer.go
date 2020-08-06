package producer

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Config struct {
	Idempotence      bool
	ReadEvents       bool
	FlushTimeoutMs   int
	BootstrapServers string
	SSL              SSLConfig
}

type SSLConfig struct {
	Enabled             bool
	KeyLocation         string
	CertificateLocation string
	CALocation          string
}

type Producer struct {
	producer       *kafka.Producer
	topic          string
	flushTimeoutMs int
}

func New(ctx context.Context, wg *sync.WaitGroup, topic string, config Config) (*Producer, error) {
	// Apply kafka producer settings
	configMap := kafka.ConfigMap{
		"bootstrap.servers":  config.BootstrapServers,
		"enable.idempotence": config.Idempotence,
	}

	if config.SSL.Enabled {
		configMap["security.protocol"] = "ssl"
		configMap["ssl.key.location"] = config.SSL.KeyLocation
		configMap["ssl.certificate.location"] = config.SSL.CertificateLocation
		configMap["ssl.ca.location"] = config.SSL.CALocation
	}

	producer, err := kafka.NewProducer(&configMap)
	if err != nil {
		return nil, err
	}

	p := &Producer{
		producer:       producer,
		topic:          topic,
		flushTimeoutMs: config.FlushTimeoutMs,
	}

	// Log broker async events
	if config.ReadEvents {
		go p.readEvents()
	}

	// Flush and wait for outstanding messages and requests to complete delivery
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		p.producer.Flush(p.flushTimeoutMs)
		p.producer.Close()
		log.Infof("kafka producer '%s': closed", p.topic)
	}()

	return p, nil
}

func (p *Producer) readEvents() {
	for event := range p.producer.Events() {
		switch e := event.(type) {
		case *kafka.Message:
			if e.TopicPartition.Error != nil {
				log.Infof("kafka producer '%s': failed to deliver message (key='%s', value='%s'): %v", p.topic, e.Key, e.Value, e.TopicPartition.Error)
			}
		case kafka.Error:
			if e.IsFatal() {
				log.Errorf("kafka producer '%s': fatal error: %v", p.topic, e)
			}
			log.Tracef("kafka producer '%s': error: %v", p.topic, e)
		}
	}
}

func (p *Producer) HealthCheck() error {
	return p.producer.GetFatalError()
}

// It is an example of how we can implement the message send to Kafka topic.
// Feel free to change this implementation if you need to add more options
// like headers, timestamp, partition etc.
func (p *Producer) Produce(key, value []byte) error {
	if err := p.producer.GetFatalError(); err != nil {
		return errors.Wrap(err, "the client instance raised a fatal error")
	}
	p.producer.ProduceChannel() <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Value: value,
		Key:   key,
	}
	return nil
}

// The same as Produce but imitates synchronous send
func (p *Producer) ProduceSync(key, value []byte) error {
	if err := p.producer.GetFatalError(); err != nil {
		return errors.Wrap(err, "the client instance raised a fatal error")
	}
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)
	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Value: value,
		Key:   key,
	}, deliveryChan)
	if err != nil {
		return err
	}
	event := <-deliveryChan
	msg := event.(*kafka.Message)
	if msg.TopicPartition.Error != nil {
		return errors.Wrap(err, "delivery failed")
	}
	return nil
}
