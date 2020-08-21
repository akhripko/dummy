package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/akhripko/dummy/src/cache/redis"
	"github.com/akhripko/dummy/src/metrics"
	"github.com/akhripko/dummy/src/options"
	"github.com/akhripko/dummy/src/service"
	"github.com/akhripko/dummy/src/srv/info"
	"github.com/akhripko/dummy/src/srv/prometheus"
	"github.com/akhripko/dummy/src/srv/srvgql"
	"github.com/akhripko/dummy/src/srv/srvgrpc"
	"github.com/akhripko/dummy/src/srv/srvhttp"
	"github.com/akhripko/dummy/src/storage/postgres"

	log "github.com/sirupsen/logrus"
)

func main() {
	// read service config from os env
	config := options.ReadEnv()

	// init logger
	initLogger(config)

	log.Info("begin...")
	// register metrics
	metrics.Register()

	// prepare main context
	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)
	var wg = &sync.WaitGroup{}

	// build storage
	storage, err := postgres.New(ctx, config.Postgres)
	if err != nil {
		log.Error("sql db init error:", err.Error())
		os.Exit(1)
	}
	// build cache
	cache, err := redis.New(ctx, config.CacheAddr)
	if err != nil {
		log.Error("cache init error:", err.Error())
		os.Exit(1)
	}

	//p, err := producer.New(ctx, wg, config.KafkaTopic.Hello, config.KafkaProducer)
	//if err != nil {
	//	log.Error("kafka producer init error:", err.Error())
	//	os.Exit(1)
	//}

	//helloConsumer, err := consumer.New(config.KafkaConsumer, []string{config.KafkaTopic.Hello}, nil)
	//if err != nil {
	//	log.Error("kafka consumer init error:", err.Error())
	//	os.Exit(1)
	//}

	//hellosrvClient, err := hellosrv.New(ctx, config.HelloSrvConf)
	//if err != nil {
	//	log.Error("hellosrv client init error:", err.Error())
	//	os.Exit(1)
	//}

	// build main service
	srv, err := service.New(storage, cache)
	if err != nil {
		log.Error("service init error:", err.Error())
		os.Exit(1)
	}

	http, err := srvhttp.New(config.HTTPPort, srv)
	if err != nil {
		log.Error("http service init error:", err.Error())
		os.Exit(1)
	}

	grpc, err := srvgrpc.New(config.GRPCPort, srv)
	if err != nil {
		log.Error("grpc service init error:", err.Error())
		os.Exit(1)
	}

	gql, err := srvgql.New(config.GraphqlPort, srv)
	if err != nil {
		log.Error("graphql service init error:", err.Error())
		os.Exit(1)
	}

	// build prometheus srv
	prometheusSrv := prometheus.New(config.PrometheusPort)

	// build info srv
	infoSrv := info.New(
		config.InfoPort,
		storage.Check,
		cache.Check,
		prometheusSrv.HealthCheck,
		http.HealthCheck,
		gql.HealthCheck,
		grpc.HealthCheck,
		//helloConsumer.HealthCheck,
	)

	// run server
	http.Run(ctx, wg)
	grpc.Run(ctx, wg)
	gql.Run(ctx, wg)
	infoSrv.Run(ctx, wg)
	prometheusSrv.Run(ctx, wg)
	//helloConsumer.Run(ctx, wg)

	// wait while services work
	wg.Wait()
	log.Info("end")
}

func initLogger(config *options.Config) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)

	switch strings.ToLower(config.LogLevel) {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}

func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		log.Info("got Interrupt signal")
		stop()
	}()
}
