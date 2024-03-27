package connector

import (
	"context"
	"fmt"
	pb "github.com/nessai1/gophkeeper/api/proto"
	"github.com/nessai1/gophkeeper/internal/keeper/secret"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"io"
	"os"
	"time"
)

type GRPCServiceConnector struct {
	authToken string

	connection *grpc.ClientConn
	client     pb.KeeperServiceClient
}

const uploadBlockSize = 256 * 524288 // ~0.5mb

func (c *GRPCServiceConnector) UploadMedia(ctx context.Context, name string, reader io.Reader, replace bool) (string, error) {
	stream, err := c.client.UploadMediaSecret(ctx)
	if err != nil {
		return "", fmt.Errorf("cannot open stream for upload media sercret: %w", err)
	}

	md := pb.MediaSecretMetadata{Name: name, Overwrite: replace}

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

func CreateGRPCConnector(serviceAddr, CertificateTLSFilePath string) (*GRPCServiceConnector, error) {
	server := &GRPCServiceConnector{}

	var creds credentials.TransportCredentials
	if CertificateTLSFilePath == "" {
		creds = insecure.NewCredentials()
	} else {
		var err error
		creds, err = buildTLSCredentials(CertificateTLSFilePath)
		if err != nil {
			return nil, fmt.Errorf("cannot build GRPC connector")
		}
	}

	conn, err := grpc.Dial(
		serviceAddr,
		grpc.WithTransportCredentials(creds),
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

func buildTLSCredentials(crtPath string) (credentials.TransportCredentials, error) {
	creds, err := credentials.NewClientTLSFromFile(crtPath, "")
	if err != nil {
		return nil, fmt.Errorf("cannot build TLS credentials from certificate file: %w", err)
	}

	return creds, nil
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

func (c *GRPCServiceConnector) ListSecret(ctx context.Context, secretType secret.SecretType) ([]secret.Secret, error) {
	translatedType, err := translateSecretTypeTypeToGRPCType(secretType)
	if err != nil {
		return nil, fmt.Errorf("cannot list secrets: %w", err)
	}

	resp, err := c.client.SecretList(ctx, &pb.SecretListRequest{SecretType: translatedType})
	if err != nil {
		return nil, fmt.Errorf("error while list secrets from external service: %w", err)
	}

	secrets := make([]secret.Secret, len(resp.Secrets))
	for i, v := range resp.Secrets {
		secrets[i] = secret.Secret{
			SecretType: secretType,
			Name:       v.Name,
			Created:    time.Unix(v.CreateTimestamp, 0),
			Updated:    time.Unix(v.UpdateTimestamp, 0),
		}
	}

	return secrets, nil
}

func (c *GRPCServiceConnector) SetSecret(ctx context.Context, name string, secretType secret.SecretType, data []byte) error {
	translatedType, err := translateSecretTypeTypeToGRPCType(secretType)
	if err != nil {
		return fmt.Errorf("cannot set secret: %w", err)
	}

	_, err = c.client.SecretSet(ctx, &pb.SecretSetRequest{
		SecretType: translatedType,
		Name:       name,
		Content:    data,
	})

	if err != nil {
		return fmt.Errorf("cannot set secret: %w", err)
	}

	return nil
}

func (c *GRPCServiceConnector) UpdateSecret(ctx context.Context, name string, secretType secret.SecretType, data []byte) error {
	translatedType, err := translateSecretTypeTypeToGRPCType(secretType)
	if err != nil {
		return fmt.Errorf("cannot update secret: %w", err)
	}

	_, err = c.client.SecretUpdate(ctx, &pb.SecretUpdateRequest{
		SecretType: translatedType,
		Name:       name,
		Content:    data,
	})

	if err != nil {
		return fmt.Errorf("cannot update secret: %w", err)
	}
	return nil
}

func (c *GRPCServiceConnector) RemoveSecret(ctx context.Context, name string, secretType secret.SecretType) error {
	translatedType, err := translateSecretTypeTypeToGRPCType(secretType)
	if err != nil {
		return fmt.Errorf("cannot remove secret: %w", err)
	}

	_, err = c.client.SecretDelete(ctx, &pb.SecretDeleteRequest{
		SecretType: translatedType,
		SecretName: name,
	})
	if err != nil {
		return fmt.Errorf("cannot remove secret: %w", err)
	}

	return nil
}

func (c *GRPCServiceConnector) GetSecret(ctx context.Context, name string, secretType secret.SecretType) ([]byte, error) {
	translatedType, err := translateSecretTypeTypeToGRPCType(secretType)
	if err != nil {
		return nil, fmt.Errorf("cannot get secret: %w", err)
	}

	resp, err := c.client.SecretGet(ctx, &pb.SecretGetRequest{
		SecretType: translatedType,
		Name:       name,
	})

	if err != nil {
		return nil, fmt.Errorf("cannot get secret: %w", err)
	}

	return resp.Secret.Content, nil
}

func translateSecretTypeTypeToGRPCType(keeperSecret secret.SecretType) (pb.SecretType, error) {
	switch keeperSecret {
	case secret.SecretTypeCredentials:
		return pb.SecretType_CREDENTIALS, nil
	case secret.SecretTypeCard:
		return pb.SecretType_CREDIT_CARD, nil
	case secret.SecretTypeText:
		return pb.SecretType_TEXT, nil
	case secret.SecretTypeMedia:
		return pb.SecretType_MEDIA, nil
	default:
		return 0, fmt.Errorf("undefined secret type translated")
	}
}
