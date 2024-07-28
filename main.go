package main

import (
	"flag"
	"log"
)

func main() {
	t := mustToken()
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
