package main

import (
	"fmt"
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	broker = "212.64.36.250:1883"
	nodeID = "1d3c3d0452ce@172.17.0.2"
	// nodeID   = "1d3c3d0452ce"
	clientID = "D1C@@@1-2-3-4-6-2"
)

func main() {
	client := genClient()
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf(" client.Connect() error: %v\n", token.Error())
		panic(token.Error())
	}
	log.Println("client connected success!")
	onTopic := genOnTopic(nodeID)
	offTopic := genOffTopic(nodeID)
	println(onTopic, offTopic)

	if token := client.Subscribe(onTopic, 0, cb); token.Wait() && token.Error() != nil {
		log.Printf("client.Subscribe error: %v\n", token.Error())
		panic(token.Error())
	}

	if token := client.Subscribe(offTopic, 0, cb); token.Wait() && token.Error() != nil {
		log.Printf("client.Subscribe error: %v\n", token.Error())
		panic(token.Error())
	}

	if token := client.Subscribe("offTopic", 0, cb); token.Wait() && token.Error() != nil {
		log.Printf("client.Subscribe error: %v\n", token.Error())
		panic(token.Error())
	}

	println("Subscribe done")
	// blocked here
	select {}
}

func genOnTopic(nodeID string) string {
	// $SYS/brokers/${node}/clients/${clientid}/connected
	return fmt.Sprintf("$SYS/brokers/%s/clients/D1C@@@1-2-3-4-5-2/connected", nodeID)
}
func genOffTopic(nodeID string) string {
	return fmt.Sprintf("$SYS/brokers/%s/clients/D1C@@@1-2-3-4-5-2/disconnected", nodeID)
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

func cb(client MQTT.Client, msg MQTT.Message) {
	log.Printf("TOPIC: %s, MSG: %s, MSG_ID: %d\n", msg.Topic(), msg.Payload(), msg.MessageID())
}
