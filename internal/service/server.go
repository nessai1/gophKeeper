package service

import (
	"github.com/nessai1/gophkeeper/internal/service/config"
	"github.com/nessai1/gophkeeper/internal/service/mediastorage"
	"github.com/nessai1/gophkeeper/internal/service/plainstorage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"

	"github.com/nessai1/gophkeeper/internal/logger"
	"github.com/nessai1/gophkeeper/internal/service/mediastorage/s3storage"

	pb "github.com/nessai1/gophkeeper/api/proto"
)

var unauthorizedMethods = []string{
	pb.KeeperService_Ping_FullMethodName,
	pb.KeeperService_Register_FullMethodName,
	pb.KeeperService_Login_FullMethodName,
}

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

	log.Println("Load media storage (s3)")
	if c.S3Config == nil {
		log.Fatalf("Cannot find S3 config fields")
	}

	ms, err := s3storage.NewStorage(*c.S3Config, l)

	if err != nil {
		log.Fatalf("Cannot build S3 storage")
	}

	listen, err := net.Listen("tcp", c.Address)
	if err != nil {
		log.Fatalf("Cannot listen '%s' address for service: %s", c.Address, err)
	}

	if c.PlainStorageConfig == nil || c.PlainStorageConfig.PSQLStorage == nil {
		log.Fatalf("No one plain storage configured!")
	}

	log.Println("Load plain storage (postgres)")

	s, err := plainstorage.NewPSQLPlainStorage(*c.PlainStorageConfig.PSQLStorage, l)
	if err != nil {
		log.Fatalf("Cannot build plain psql storage: %s", err.Error())
	}

	server := Server{
		mediaStorage: ms,
		plainStorage: s,
		logger:       l,
		config:       c,
	}
	gRPCServer := grpc.NewServer(grpc.UnaryInterceptor(server.unaryAuthInterceptor), grpc.StreamInterceptor(server.streamAuthInterceptor))
	pb.RegisterKeeperServiceServer(gRPCServer, &server)

	log.Printf("Service started at %s", c.Address)
	if err := gRPCServer.Serve(listen); err != nil {
		log.Fatalf("Error while run gRPC server: %s", err.Error())
	}
}

type Server struct {
	plainStorage plainstorage.PlainStorage
	mediaStorage mediastorage.MediaStorage
	logger       *zap.Logger
	config       config.Config

	pb.UnimplementedKeeperServiceServer
}
