package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	urls := "nats://cluster1:cluster1@localhost:4222,nats://cluster1:cluster1@localhost:4223,nats://cluster1:cluster1@localhost:4224"
	// urls := "nats://@localhost:4222"
	nc, err := nats.Connect(urls)
	if err != nil {
		log.Fatal(err)
	}

	runConsumer(nc)
}

func runConsumer(conn *nats.Conn) {
	js, err := conn.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	sub, err := js.PullSubscribe("demo-subject", "", nats.DeliverNew())
	if err != nil {
		log.Fatal(err)
	}

	for {
		msgs, err := sub.Fetch(10, nats.MaxWait(5*time.Second))
		if err != nil {
			log.Printf("Error fetching messages: %v", err)
		}

		fmt.Printf("Received %d messages, subValid: %v\n", len(msgs), sub.IsValid())

		for _, msg := range msgs {
			fmt.Printf("Received a message: %s\n", string(msg.Data))
			msg.Ack()
		}
	}
}
