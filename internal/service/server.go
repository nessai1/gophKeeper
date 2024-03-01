package service

import (
	"context"
	pb "github.com/nessai1/gophkeeper/api/proto"
	"log"
)

func Run() {
	log.Println("Starting keeper service")
	log.Println("Load configuration")
	config, err := fetchConfig()
	if err != nil {
		log.Fatalf("Cannot fetch config for service: %s", err)
	}

	log.Printf("Service started at %s", config.Address)
}

type Server struct {
	pb.UnimplementedKeeperServiceServer
}

func (s *Server) Ping(ctx context.Context, request *pb.PingRequest) (*pb.PingResponse, error) {
	resp := pb.PingResponse{}
	resp.Answer = "pong!"
	return &pb.PingResponse{}, nil
}
