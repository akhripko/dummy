package kafka

import (
	"github.com/akhripko/dummy/src/kafka/consumer"
	"github.com/akhripko/dummy/src/kafka/producer"
)

type KafkaConsumerConfig struct {
	LogLevel  string
	InfoPort  int
	TopicName string
	Consumer  consumer.Config
}

type KafkaProducerConfig struct {
	LogLevel  string
	InfoPort  int
	TopicName string
	Producer  producer.Config
}
