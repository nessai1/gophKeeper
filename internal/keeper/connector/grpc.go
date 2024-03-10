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

func (c *GRPCServiceConnector) Register(ctx context.Context, login, password string) (token string, err error) {
	req := pb.UserCredentialsRequest{
		Login:    login,
		Password: password,
	}

	response, err := c.client.Register(ctx, &req)
	if err != nil {
		return "", fmt.Errorf("error while register (service error: %w)", err)
	}

	return response.Token, nil
}

func (c *GRPCServiceConnector) Login(ctx context.Context, login, password string) (token string, err error) {
	req := pb.UserCredentialsRequest{
		Login:    login,
		Password: password,
	}

	response, err := c.client.Login(ctx, &req)
	if err != nil {
		return "", fmt.Errorf("error while login (service error: %w)", err)
	}

	return response.Token, nil
}
