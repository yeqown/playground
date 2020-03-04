package main

import (
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	broker   = "212.64.36.250:1883"
	nodeID   = "1d3c3d0452ce@172.17.0.2"
	clientID = "D1C@@@1-2-3-4-5-2"
)

func main() {
	client := genClient()
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf(" client.Connect() error: %v\n", token.Error())
		panic(token.Error())
	}
	log.Println("client connected success!")
	// blocked here
	select {}
}

func genClient() MQTT.Client {
	opts := MQTT.NewClientOptions()

	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetAutoReconnect(true)
	opts.SetCleanSession(false)
	opts.SetProtocolVersion(4)
	opts.OnConnect = onConnect

	client := MQTT.NewClient(opts)

	return client
}

func onConnect(c MQTT.Client) {
	log.Printf("Connected to MQTT broker, and status: %v, open status: %v\n",
		c.IsConnected(), c.IsConnectionOpen())
}
