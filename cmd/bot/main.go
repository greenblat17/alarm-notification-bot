package main

import (
	"flag"
	"log"

	"github.com/greenblat17/alarm-notification-bot/internal/clients/telegram"
)

func main() {
	tgClient := telegram.New(mustHost(), mustToken())
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
