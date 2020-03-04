package main

import (
	"log"

	"github.com/nsqio/go-nsq"
)

func main() {
	cfg := nsq.NewConfig()

	// connect to nsqd based TCP:4150
	// configs refer to docker-compose.yaml
	p, err := nsq.NewProducer("127.0.0.1:4150", cfg)
	if err != nil {
		log.Fatal(err)
	}

	// check connection
	if err := p.Ping(); err != nil {
		log.Fatal(err)
	}

	// publish message
	if err := p.Publish("order_created", []byte("product id is 213172")); err != nil {
		log.Fatal(err)
	}
}
