package service

import (
	"context"
	"github.com/nessai1/gophkeeper/internal/service/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"

	"github.com/nessai1/gophkeeper/internal/logger"
	"github.com/nessai1/gophkeeper/internal/service/storage"
	"github.com/nessai1/gophkeeper/internal/service/storage/s3storage"

	pb "github.com/nessai1/gophkeeper/api/proto"
)

func Run() {
	log.Println("Starting keeper service")
	log.Println("Load configuration")
	c, err := config.FetchConfig()
	if err != nil {
		log.Fatalf("Cannot fetch config for service: %s", err.Error())
	}

	l, err := logger.BuildLogger(logger.LevelDev, os.Stdout)
	if err != nil {
		log.Fatalf("Cannot build logger: %s", err.Error())
	}

	if c.S3Config == nil {
		log.Fatalf("Cannot find S3 config fields")
	}

	s, err := s3storage.NewStorage(*c.S3Config, l)

	if err != nil {
		log.Fatalf("Cannot build S3 storage")
	}

	listen, err := net.Listen("tcp", c.Address)
	if err != nil {
		log.Fatalf("Cannot listen '%s' address for service: %s", c.Address, err)
	}

	server := Server{
		storage: s,
		logger:  l,
		config:  c,
	}
	gRPCServer := grpc.NewServer()
	pb.RegisterKeeperServiceServer(gRPCServer, &server)

	log.Printf("Service started at %s", c.Address)
	if err := gRPCServer.Serve(listen); err != nil {
		log.Fatalf("Error while run gRPC server: %s", err.Error())
	}
}

type Server struct {
	storage storage.Storage
	logger  *zap.Logger
	config  config.Config

	pb.UnimplementedKeeperServiceServer
}

func (s *Server) Ping(ctx context.Context, request *pb.PingRequest) (*pb.PingResponse, error) {
	s.logger.Info("Got ping")
	resp := pb.PingResponse{}
	resp.Answer = "pong!"
	return &resp, nil
}
