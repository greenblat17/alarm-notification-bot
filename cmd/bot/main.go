package main

import (
	"flag"
	"log"

	tgClient "github.com/greenblat17/alarm-notification-bot/internal/clients/telegram"
	eventconsumer "github.com/greenblat17/alarm-notification-bot/internal/consumer/event-consumer"
	"github.com/greenblat17/alarm-notification-bot/internal/events/telegram"
	"github.com/greenblat17/alarm-notification-bot/internal/storage/files"
)

const (
	storagePath = "storage"
	batchSize   = 100
)

func main() {
	eventProcessor := telegram.New(
		tgClient.New(mustHost(), mustToken()),
		files.New(storagePath),
	)

	log.Println("service started")

	consumer := eventconsumer.New(eventProcessor, eventProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatalf("service stopped: %v", err)
	}
}

func mustToken() string {
	var token string
	flag.StringVar(&token, "token", "", "token for access to telegram bot")

	flag.Parse()

	if token == "" {
		log.Fatal("token is not specified")
	}

	return token
}

func mustHost() string {
	var host string
	flag.StringVar(&host, "host", "", "host for telegram bot API")

	flag.Parse()

	if host == "" {
		log.Fatal("host is not specified")
	}

	return host
}
