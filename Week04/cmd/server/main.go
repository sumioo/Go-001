package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	pb "hello.com/api"
	"hello.com/internal/data"
	"hello.com/internal/service"
)

const (
	port = ":50051"
)

// import "fmt"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	repo := data.NewArticleRepo()
	srv := service.NewArticleServer(repo)
	pb.RegisterArticleServer(s, srv)
	log.Println("start server at", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
