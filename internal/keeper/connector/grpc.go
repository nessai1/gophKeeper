package connector

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/nessai1/gophkeeper/api/proto"
)

type GRPCServiceConnector struct {
	connection *grpc.ClientConn
	client     pb.KeeperServiceClient
}

func CreateGRPCConnector(serviceAddr string) (*GRPCServiceConnector, error) {
	conn, err := grpc.Dial(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("cannot create gRPC connector: %w", err)
	}

	return &GRPCServiceConnector{
		connection: conn,
		client:     pb.NewKeeperServiceClient(conn),
	}, nil
}

func (c *GRPCServiceConnector) Ping(ctx context.Context) (answer string, err error) {
	response, err := c.client.Ping(ctx, &pb.PingRequest{Message: "ping"})
	if err != nil {
		return "", fmt.Errorf("connector error while ping service: %w", err)
	}

	return response.Answer, nil
}
