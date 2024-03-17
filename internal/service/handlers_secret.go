package service

import (
	"context"
	"fmt"
	pb "github.com/nessai1/gophkeeper/api/proto"
	"github.com/nessai1/gophkeeper/internal/service/plainstorage"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) SecretList(ctx context.Context, request *pb.SecretListRequest) (*pb.SecretListResponse, error) {
	userCtxVal := ctx.Value(UserContextKey)
	user := userCtxVal.(*plainstorage.User)

	s.logger.Info("User list secrets", zap.String("login", user.Login), zap.String("secret_type", request.SecretType.String()))
	secretType, err := translateGRPCSecretTypeToSecretType(request.GetSecretType())
	if err != nil {
		s.logger.Info("User sends invalid secret type", zap.String("login", user.Login))

		return nil, status.Error(codes.InvalidArgument, "invalid secret type got")
	}

	secrets, err := s.plainStorage.GetUserSecretsMetadataByType(ctx, user.UUID, secretType)
	if err != nil {
		s.logger.Error("Cannot get list of secrets for user", zap.String("login", user.Login))

		return nil, status.Error(codes.Internal, "cannot list user secrets")
	}

	outputSecrets := make([]*pb.Secret, len(secrets))
	for i, v := range secrets {
		outputSecrets[i] = &pb.Secret{
			SecretType:      request.SecretType,
			Name:            v.Name,
			CreateTimestamp: v.Created.Unix(),
			UpdateTimestamp: v.Updated.Unix(),
			Content:         nil,
		}
	}

	return &pb.SecretListResponse{
		Secrets: outputSecrets,
	}, nil
}

func translateGRPCSecretTypeToSecretType(secretType pb.SecretType) (plainstorage.SecretType, error) {
	switch secretType {
	case pb.SecretType_CREDENTIALS:
		return plainstorage.SecretTypeCredentials, nil
	case pb.SecretType_TEXT:
		return plainstorage.SecretTypeText, nil
	case pb.SecretType_MEDIA:
		return plainstorage.SecretTypeMedia, nil
	case pb.SecretType_CREDIT_CARD:
		return plainstorage.SecretTypeCard, nil
	default:
		return 0, fmt.Errorf("got undefined secret type to translate")
	}
}
