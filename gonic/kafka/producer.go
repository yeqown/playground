package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
)

// Sarama configuration options
var (
	brokers = ""
	version = ""
	// group   = ""
	topic = ""
	// assignor = ""
	// oldest   = true
	verbose = false
)

func init() {
	flag.StringVar(&brokers, "brokers", "10.111.153.97:9092,10.111.153.96:9092,10.111.153.98:9092", "Kafka bootstrap brokers to connect to, as a comma separated list")
	// flag.StringVar(&group, "group", "", "Kafka consumer group definition")
	flag.StringVar(&version, "version", "2.1.1", "Kafka cluster version")
	flag.StringVar(&topic, "topic", "", "Kafka topics to be consumed, as a comma separated list")
	// flag.StringVar(&assignor, "assignor", "roundrobin", "Consumer group partition assignment strategy (range, roundrobin, sticky)")
	// flag.BoolVar(&oldest, "oldest", true, "Kafka consumer consume initial offset from oldest")
	flag.BoolVar(&verbose, "verbose", false, "Sarama logging")
	flag.Parse()

	if len(brokers) == 0 {
		panic("no Kafka bootstrap brokers defined, please set the -brokers flag")
	}

	if len(topic) == 0 {
		panic("no topics given to be consumed, please set the -topics flag")
	}

	// if len(group) == 0 {
	// 	panic("no Kafka consumer group defined, please set the -group flag")
	// }
}

func main() {

	if verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	//version, err := sarama.ParseKafkaVersion(version)
	//if err != nil {
	//	log.Panicf("Error parsing Kafka version: %v", err)
	//}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	//config := sarama.NewConfig()
	//config.Version = version
	//config.Producer.Return.Successes = true
	//config.Producer.Timeout = 5 * time.Second

	config := sarama.NewConfig()
	config.Producer.Retry.Max = 3
	config.Producer.Retry.Backoff = 1000 * time.Millisecond
	config.Producer.Return.Successes = true

	sarama.Logger.Println(brokers)
	producer, err := sarama.NewSyncProducer(strings.Split(brokers, ","), config)
	if err != nil {
		panic(err)
	}

	sarama.Logger.Println("procuder starting")
	ctx, cancel := context.WithCancel(context.Background())
	cnt := 0
	quit := make(chan struct{})

	go func() {
		for {
			cnt++
			time.Sleep(250 * time.Millisecond)
			partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
				Topic: topic,
				Key:   sarama.StringEncoder(strconv.Itoa(cnt)),
				Value: sarama.StringEncoder(fmt.Sprintf("%d-%d", cnt, time.Now().Unix())),
				// Headers: nil,
				// Metadata:  nil,
				// Offset:    0,
				// Partition: 0,
				Timestamp: time.Now(),
			})
			sarama.Logger.Printf("SendMessage, ERR=%v, part=%d, offset=%d\n", err, partition, offset)

			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				sarama.Logger.Printf("ERROR: %v", ctx.Err())
				quit <- struct{}{}
				return
			}
		}
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		cancel()
		log.Println("terminating: via signal")
	case <-quit:
		break
	}

	if err = producer.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}
