package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

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
	log.Debugf("config: %+v\n", *config)

	// prepare main context
	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)
	var wg = &sync.WaitGroup{}

	handler := &Handler{}

	client, err := consumer.New(config.Consumer, []string{config.TopicName}, handler)
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

	go dowork("http://kafka-service:3030")
	go dowork("http://kafka-service:9092")

	// wait while services work
	wg.Wait()
	log.Info("end")
}

func initLogger(config *options.ConsumerConfig) {
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

type Handler struct {
}

func (h *Handler) Handle(ctx context.Context, key, value []byte, timestamp time.Time) error {
	log.Printf("key:%s, value:%s", string(key), string(value))
	return nil
}

func dowork(url string) {
	log.Println("dowork: url:", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(url, " err:", err.Error())
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(url, " err:", err.Error())
		return
	}
	log.Println(url, " res:", string(data))
}
