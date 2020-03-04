package main

import (
	"google.golang.org/grpc"
)

func main() {
	clientConn, err := grpc.Dial(":8080", nil)
	_, _ = clientConn, err
}
