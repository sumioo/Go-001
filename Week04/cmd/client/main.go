// Package main implements a client for Greeter service.
package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "hello.com/api"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewArticleClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	r, err := c.CreateArticle(ctx, &pb.CreateArticleRequest{Title: "天冷了", Content: "天气冷了，多穿衣服"})
	if err != nil {
		log.Fatalf("could not create article: %v", err)
	}
	log.Printf("create success: %d %s", r.GetId(), r.GetTitle())
}
