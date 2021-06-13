package main

import (
	"context"
	"log"
	"net"

	h "github.com/Msesto/cloudnative-go/src/handlers"
	pb "github.com/Msesto/cloudnative-go/src/keyvalue"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedKeyValueServer
}

func (s *server) Get(ctx context.Context, r *pb.GetRequest) (*pb.GetResponse, error) {

	log.Printf("Received GET key=%v", r.Key)

	value, err := h.Get(r.Key)

	return &pb.GetResponse{Value: value}, err
}

func main() {
	s := grpc.NewServer()
	pb.RegisterKeyValueServer(s, &server{})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
