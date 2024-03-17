package connector

import (
	"context"
	"fmt"
	pb "github.com/nessai1/gophkeeper/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"io"
	"os"
)

type GRPCServiceConnector struct {
	authToken string

	connection *grpc.ClientConn
	client     pb.KeeperServiceClient
}

const uploadBlockSize = 256 * 524288 // ~0.5mb

func (c *GRPCServiceConnector) UploadMedia(ctx context.Context, name string, reader io.Reader) (string, error) {
	stream, err := c.client.UploadMediaSecret(ctx)
	if err != nil {
		return "", fmt.Errorf("cannot open stream for upload media sercret: %w", err)
	}

	md := pb.MediaSecretMetadata{Name: name}

	err = stream.Send(&pb.UploadMediaSecretRequest{Request: &pb.UploadMediaSecretRequest_Metadata{Metadata: &md}})
	if err != nil {
		return "", fmt.Errorf("cannot send metadata in media secret upload stream: %w", err)
	}

	var b []byte
	for {
		b = make([]byte, uploadBlockSize)
		n, err := reader.Read(b)
		if err == io.EOF {
			break
		}

		if err != nil {
			return "", fmt.Errorf("cannot read media for upload: %w", err)
		}

		if n != uploadBlockSize {
			b = b[:n]
		}

		err = stream.Send(&pb.UploadMediaSecretRequest{Request: &pb.UploadMediaSecretRequest_Data{Data: &pb.MediaSecret{Chunk: b}}})
		if err != nil {
			return "", fmt.Errorf("cannot send media chunk: %w", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return "", fmt.Errorf("cannot close stream for upload media secret: %w", err)
	}

	return res.Uuid, nil
}

func CreateGRPCConnector(serviceAddr string) (*GRPCServiceConnector, error) {
	server := &GRPCServiceConnector{}
	conn, err := grpc.Dial(
		serviceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(server.unaryAuthInterceptor),
		grpc.WithStreamInterceptor(server.streamAuthInterceptor),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create gRPC connector: %w", err)
	}

	server.connection = conn
	server.client = pb.NewKeeperServiceClient(conn)

	return server, nil
}

func (c *GRPCServiceConnector) unaryAuthInterceptor(ctx context.Context, method string, req interface{},
	reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {

	if c.authToken != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "jwt", c.authToken)
	}

	err := invoker(ctx, method, req, reply, cc, opts...)

	return err
}

func (c *GRPCServiceConnector) streamAuthInterceptor(ctx context.Context, desc *grpc.StreamDesc,
	cc *grpc.ClientConn, method string, streamer grpc.Streamer,
	opts ...grpc.CallOption) (grpc.ClientStream, error) {

	if c.authToken != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "jwt", c.authToken)
	}

	stream, err := streamer(ctx, desc, cc, method, opts...)

	return stream, err
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

func (c *GRPCServiceConnector) SetAuthToken(token string) {
	c.authToken = token
}

func (c *GRPCServiceConnector) DownloadMedia(ctx context.Context, name string, dest string) (*os.File, error) {
	stream, err := c.client.DownloadMediaSecret(ctx, &pb.DownloadMediaSecretRequest{SecretName: name})
	if err != nil {
		return nil, fmt.Errorf("cannot start media download: %w", err)
	}

	f, err := os.Create(dest)
	if err != nil {
		return nil, fmt.Errorf("cannot create file to dest %s: %w", dest, err)
	}

	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			f.Close()

			return nil, fmt.Errorf("cannot receive package of media file: %w", err)
		}

		_, err = f.Write(recv.GetSecretPart().Chunk)
		if err != nil {
			f.Close()

			return nil, fmt.Errorf("cannot write media chunk to file while download: %w", err)
		}
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("cannot seek destination file: %w", err)
	}

	return f, nil
}
