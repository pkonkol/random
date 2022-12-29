package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "grpc.test/src/proto"
)

var port = flag.Int("port", 50051, "The server port")

type server struct {
	pb.UnimplementedTestServer
}

func (s *server) SendRequest(ctx context.Context, in *pb.TestRequest) (*pb.TestReply, error) {
	log.Printf("received message: %v - %v", in.GetCounter(), in.GetMessage())
	return &pb.TestReply{
			Message: fmt.Sprintf("reply from golang for %s: %d", in.GetMessage(), in.GetCounter()),
			Counter: in.GetCounter() + 100},
		nil
}

func (s *server) ClientStream(pb.Test_ClientStreamServer) error {
	return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTestServer(s, &server{})
	log.Println("server listening at: %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
