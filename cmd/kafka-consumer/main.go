package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/akhripko/dummy/src/kafka/consumer"
	options "github.com/akhripko/dummy/src/options/kafka"
	"github.com/akhripko/dummy/src/srv/info"

	log "github.com/sirupsen/logrus"
)

func main() {
	// read service config from os env
	config := options.ReadKafkaConsumerEnv()

	// init logger
	initLogger(config)

	log.Info("begin...")

	// prepare main context
	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)
	var wg = &sync.WaitGroup{}

	client, err := consumer.New(config.Consumer, []string{config.TopicName}, nil)
	if err != nil {
		log.Error("kafka consumer init error:", err.Error())
		os.Exit(1)
	}

	// build info srv
	infoSrv := info.New(
		config.InfoPort,
		client.HealthCheck,
	)

	// run server
	infoSrv.Run(ctx, wg)
	client.Run(ctx, wg)

	// wait while services work
	wg.Wait()
	log.Info("end")
}

func initLogger(config *options.KafkaConsumerConfig) {
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
