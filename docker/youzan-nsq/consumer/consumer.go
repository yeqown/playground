package main

import (
	"log"

	"github.com/nsqio/go-nsq"
)

type customConsumer struct {
	consumer *nsq.Consumer
	Idx      int
	Channel  string
	Topics   string
	Error    error
}

func main() {
	cfg := nsq.NewConfig()
	var consumers = make([]*customConsumer, 4)

	// initial
	consumers[0] = &customConsumer{Idx: 0, Channel: "channel1", Topics: "order_created"}
	consumers[1] = &customConsumer{Idx: 1, Channel: "channel1", Topics: "order_created"}
	consumers[2] = &customConsumer{Idx: 2, Channel: "channel1", Topics: "order_created"}
	consumers[3] = &customConsumer{Idx: 3, Channel: "channel2", Topics: "order_created"}

	// build connection
	// one message(Topic: order_created) for "channel1" consumers only one can recv
	// but for "channel2" consumer can always recved (there's only one consumer)
	// channel1 consumers
	consumers[0].consumer, _ = nsq.NewConsumer("order_created", "channel1", cfg)
	consumers[1].consumer, _ = nsq.NewConsumer("order_created", "channel1", cfg)
	consumers[2].consumer, _ = nsq.NewConsumer("order_created", "channel1", cfg)
	// channel2
	consumers[3].consumer, _ = nsq.NewConsumer("order_created", "channel2", cfg)

	// bind handlers
	for idx := range consumers {
		consumers[idx].consumer.AddConcurrentHandlers(handler{idx: idx}, 1)
		// connect
		if err := consumers[idx].consumer.ConnectToNSQD("127.0.0.1:4150"); err != nil {
			log.Fatal(err)
		}
	}

	// block main goroutine
	var blocked chan struct{}
	<-blocked
}

type handler struct {
	idx int
}

func (h handler) HandleMessage(msg *nsq.Message) error {
	log.Printf("consumer[%d] recv msg(%s)", h.idx, msg.Body)
	return nil
}
