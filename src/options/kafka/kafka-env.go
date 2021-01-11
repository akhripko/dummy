package kafka

import (
	"github.com/akhripko/dummy/src/kafka/consumer"
	"github.com/akhripko/dummy/src/kafka/producer"

	"github.com/spf13/viper"
)

func ReadKafkaConsumerEnv() *KafkaConsumerConfig {
	viper.AutomaticEnv()

	viper.SetEnvPrefix("APP")

	viper.SetDefault("LOG_LEVEL", "DEBUG")
	viper.SetDefault("INFO_PORT", 8888)

	return &KafkaConsumerConfig{
		LogLevel:  viper.GetString("LOG_LEVEL"),
		InfoPort:  viper.GetInt("INFO_PORT"),
		TopicName: viper.GetString("KAFKA_TOPIC_NAME"),
		Consumer: consumer.Config{
			BootstrapServers: viper.GetString("KAFKA_BOOTSTRAP_SERVERS"),
			PollTimeoutMs:    viper.GetInt("KAFKA_CONSUMER_POLL_TIMEOUT_MS"),
			Name:             viper.GetString("KAFKA_CONSUMER_NAME"),
			GroupID:          viper.GetString("KAFKA_CONSUMER_GROUP_ID"),
			SessionTimeoutMs: viper.GetString("KAFKA_CONSUMER_SESSION_TIMEOUT_MS"),
			AutoOffsetReset:  viper.GetString("KAFKA_CONSUMER_AUTO_OFFSET_RESET"),
			SSL: consumer.SSLConfig{
				Enabled:             viper.GetBool("KAFKA_SSL_ENABLED"),
				KeyLocation:         viper.GetString("KAFKA_SSL_KEY_LOCATION"),
				CertificateLocation: viper.GetString("KAFKA_SSL_CERTIFICATE_LOCATION"),
				CALocation:          viper.GetString("KAFKA_SSL_CA_LOCATION"),
			},
		},
	}
}

func ReadKafkaProducerEnv() *KafkaProducerConfig {
	viper.AutomaticEnv()

	viper.SetEnvPrefix("APP")

	viper.SetDefault("LOG_LEVEL", "DEBUG")
	viper.SetDefault("INFO_PORT", 8888)

	return &KafkaProducerConfig{
		LogLevel:  viper.GetString("LOG_LEVEL"),
		InfoPort:  viper.GetInt("INFO_PORT"),
		TopicName: viper.GetString("KAFKA_TOPIC_NAME"),
		Producer: producer.Config{
			Idempotence:      viper.GetBool("KAFKA_PRODUCER_IDEMPOTENCE"),
			ReadEvents:       viper.GetBool("KAFKA_PRODUCER_READ_EVENTS"),
			FlushTimeoutMs:   viper.GetInt("KAFKA_PRODUCER_FLUSH_TIMEOUT_MS"),
			BootstrapServers: viper.GetString("KAFKA_BOOTSTRAP_SERVERS"),
			SSL: producer.SSLConfig{
				Enabled:             viper.GetBool("KAFKA_SSL_ENABLED"),
				KeyLocation:         viper.GetString("KAFKA_SSL_KEY_LOCATION"),
				CertificateLocation: viper.GetString("KAFKA_SSL_CERTIFICATE_LOCATION"),
				CALocation:          viper.GetString("KAFKA_SSL_CA_LOCATION"),
			},
		},
	}
}
