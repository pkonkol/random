package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "grpc.test/src/proto"
)

// const defaultName = "world"

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", "huj123", "prob unneeded")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTestClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Talk(ctx, &pb.TestRequest{Str: "asdf", Counter: 1})
	if err != nil {
		log.Fatalf("could not talk: %v", err)
	}
	log.Printf("Response is: %v - %v\n", r.GetCounter(), r.GetStr())
}
