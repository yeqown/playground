package main

import (
	"fmt"
	"log"
	"net"

	errorhandling "github.com/playground/golang/grpc/error-handling"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/playground/golang/grpc/codegen"
)

// protoc -I ../ ./hello.proto --go_out=plugins=grpc:./codegen

type HelloServer struct{}

func (s *HelloServer) SayHello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	return &pb.HelloResp{Result: fmt.Sprintf("Hey, %s!", req.GetName())}, nil
}

func (s *HelloServer) SayHelloStrict(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	if len(req.GetName()) >= 10 {
		return nil, status.Errorf(codes.InvalidArgument,
			"Length of `Name` cannot be more than 10 characters")
	}

	return &pb.HelloResp{Result: fmt.Sprintf("Hey, %s!", req.GetName())}, nil
}

//func (s *HelloServer) ErrorExample(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
//	err := errors.New("this normal error")
//	return nil, err
//}
//
// client got error:
// err=rpc error: code = Unknown desc = this normal error%

func (s *HelloServer) ErrorExample(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	err := errorhandling.New(100, "custom error")
	st := errorhandling.WrapErrorAsStatus(err)
	return nil, st.Err()
}

func Serve() {
	addr := fmt.Sprintf(":%d", 50051)
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Cannot listen to address %s", addr)
	}
	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &HelloServer{})
	if err := s.Serve(conn); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	Serve()
}
