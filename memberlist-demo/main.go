package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/memberlist"
)

type EventDelegate struct {
	name string
}

func (e *EventDelegate) NotifyJoin(node *memberlist.Node) {
	fmt.Printf("[%s] Node joined: %s (%s:%d)\n", e.name, node.Name, node.Addr, node.Port)
}

func (e *EventDelegate) NotifyLeave(node *memberlist.Node) {
	fmt.Printf("[%s] Node left: %s (%s:%d)\n", e.name, node.Name, node.Addr, node.Port)
}

func (e *EventDelegate) NotifyUpdate(node *memberlist.Node) {
	fmt.Printf("[%s] Node updated: %s (%s:%d)\n", e.name, node.Name, node.Addr, node.Port)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <node-id> [join-address]")
	}

	nodeID := os.Args[1]
	port := 7946 + parseInt(nodeID)

	config := memberlist.DefaultLocalConfig()
	config.Name = "node-" + nodeID
	config.BindPort = port
	config.AdvertisePort = port
	config.Events = &EventDelegate{name: config.Name}

	list, err := memberlist.Create(config)
	if err != nil {
		log.Fatal("Failed to create memberlist: ", err)
	}

	if len(os.Args) > 2 {
		joinAddr := os.Args[2]
		_, err := list.Join([]string{joinAddr})
		if err != nil {
			log.Fatal("Failed to join cluster: ", err)
		}
	}

	fmt.Printf("Node %s started on port %d\n", config.Name, port)
	fmt.Printf("Members: %d\n", list.NumMembers())

	for {
		time.Sleep(5 * time.Second)
		members := list.Members()
		fmt.Printf("[%s] Current members (%d):\n", config.Name, len(members))
		for _, member := range members {
			fmt.Printf("  - %s (%s:%d)\n", member.Name, member.Addr, member.Port)
		}
		fmt.Println()
	}
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
