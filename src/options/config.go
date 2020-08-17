package options

import (
	//"github.com/akhripko/dummy/src/kafka/consumer"
	//"github.com/akhripko/dummy/src/kafka/producer"
	"github.com/akhripko/dummy/src/providers/grpc/hellosrv"
	"github.com/akhripko/dummy/src/storage/postgres"
)

type Config struct {
	LogLevel       string
	HTTPPort       int
	GraphqlPort    int
	GRPCPort       int
	InfoPort       int
	PrometheusPort int
	Storage        postgres.Config
	CacheAddr      string
	HelloSrvConf   hellosrv.Config
	//KafkaTopic      KafkaTopic
	//KafkaProducer   producer.Config
	//KafkaConsumer   consumer.Config
}
