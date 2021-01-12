package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/akhripko/dummy/src/kafka/producer"
	options "github.com/akhripko/dummy/src/options/kafka"
	"github.com/akhripko/dummy/src/srv/info"

	log "github.com/sirupsen/logrus"
)

func main() {
	// read service config from os env
	config := options.ReadKafkaProducerEnv()

	// init logger
	initLogger(config)

	log.Info("begin...")
	log.Debugf("config: %+v\n", *config)

	// prepare main context
	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)
	var wg = &sync.WaitGroup{}

	client, err := producer.New(ctx, wg, config.TopicName, config.Producer)
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

	go generateData(ctx, client)

	// wait while services work
	wg.Wait()
	log.Info("end")
}

func initLogger(config *options.ProducerConfig) {
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

func generateData(ctx context.Context, client *producer.Producer) {
	var n int
	for {
		select {
		case <-ctx.Done():
			log.Debug("stop generation: canceled context")
			return
		case <-time.After(time.Second):
			n++
			key := time.Now().UTC().Format(time.RFC3339)
			data := fmt.Sprintf("recNo:%d time:%s", n, time.Now().UTC().Format(time.RFC3339))
			log.Debugf("generate data: key=%s data=%s", key, data)
			if err := client.ProduceSync([]byte(key), []byte(data)); err != nil {
				log.Errorf("failed to produce data: %s", err.Error())
			}
			log.Debug("generate data")
		}
	}
}
