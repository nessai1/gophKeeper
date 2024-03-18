package service

import (
	"context"
	"errors"
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

func (s *Server) SecretDelete(ctx context.Context, request *pb.SecretDeleteRequest) (*pb.SecretDeleteResponse, error) {
	userCtxVal := ctx.Value(UserContextKey)
	user := userCtxVal.(*plainstorage.User)

	s.logger.Info("User try delete secret", zap.String("login", user.Login), zap.String("secret_name", request.GetSecretName()), zap.String("secret_type", request.GetSecretType().String()))

	secretType, err := translateGRPCSecretTypeToSecretType(request.GetSecretType())
	if err != nil {
		s.logger.Error("User send invalid secret type", zap.String("login", user.Login))

		return nil, status.Error(codes.InvalidArgument, "invalid secret type got")
	}

	userSecret, err := s.plainStorage.GetUserSecretByName(ctx, user.UUID, request.SecretName, secretType)
	if err != nil && errors.Is(plainstorage.ErrSecretNotFound, err) {
		s.logger.Info("Cannot find secret by name", zap.String("login", user.Login), zap.String("secret_name", request.SecretName))

		return nil, status.Errorf(codes.NotFound, "secret '%s' not found", request.SecretName)
	} else if err != nil {
		s.logger.Error("Error while get secret", zap.Error(err), zap.String("login", user.Login), zap.String("secret_name", request.SecretName))

		return nil, status.Error(codes.Internal, "internal error while get secret")
	}

	if request.GetSecretType() == pb.SecretType_MEDIA {
		err = s.mediaStorage.Delete(ctx, userSecret.Metadata.UUID)
		if err != nil {
			s.logger.Error("Error while delete media secret", zap.Error(err), zap.String("login", user.Login), zap.String("secret_name", request.SecretName))

			return nil, status.Error(codes.Internal, "internal error while remove media secret")
		}
	}

	err = s.plainStorage.RemoveSecretByUUID(ctx, userSecret.Metadata.UUID)
	if err != nil {
		s.logger.Error("Cannot remove secret from DB", zap.Error(err), zap.String("login", user.Login), zap.String("secret_name", request.SecretName))

		return nil, status.Error(codes.Internal, "cannot remove secret from DB")
	}

	return &pb.SecretDeleteResponse{}, nil
}

func (s *Server) SecretSet(ctx context.Context, request *pb.SecretSetRequest) (*pb.SecretSetResponse, error) {
	userCtxVal := ctx.Value(UserContextKey)
	user := userCtxVal.(*plainstorage.User)

	s.logger.Info("User try set plain secret", zap.String("login", user.Login), zap.String("secret_name", request.GetName()), zap.String("secret_type", request.GetSecretType().String()))

	secretType, err := translateGRPCSecretTypeToSecretType(request.GetSecretType())
	if err != nil {
		s.logger.Error("User send invalid secret type", zap.String("login", user.Login))

		return nil, status.Error(codes.InvalidArgument, "invalid secret type got")
	}

	if secretType == plainstorage.SecretTypeMedia {
		s.logger.Error("User try to set media secret as plain", zap.String("login", user.Login))

		return nil, status.Error(codes.InvalidArgument, "cannot set media secret to plain storage")
	}

	_, err = s.plainStorage.AddPlainSecret(ctx, user.UUID, request.Name, secretType, request.Content)
	if err != nil {
		s.logger.Error("Cannot add plain secret", zap.String("login", user.Login), zap.Error(err))

		return nil, status.Error(codes.Internal, "internal error while save secret")
	}

	return &pb.SecretSetResponse{}, nil
}

func (s *Server) SecretGet(ctx context.Context, request *pb.SecretGetRequest) (*pb.SecretGetResponse, error) {
	userCtxVal := ctx.Value(UserContextKey)
	user := userCtxVal.(*plainstorage.User)

	s.logger.Info("User try get plain secret", zap.String("login", user.Login), zap.String("secret_name", request.GetName()), zap.String("secret_type", request.GetSecretType().String()))

	secretType, err := translateGRPCSecretTypeToSecretType(request.GetSecretType())
	if err != nil {
		s.logger.Error("User send invalid secret type", zap.String("login", user.Login))

		return nil, status.Error(codes.InvalidArgument, "invalid secret type got")
	}

	if secretType == plainstorage.SecretTypeMedia {
		s.logger.Error("User try to get media secret as plain", zap.String("login", user.Login))

		return nil, status.Error(codes.InvalidArgument, "cannot get media secret to plain storage")
	}

	secret, err := s.plainStorage.GetUserSecretByName(ctx, user.UUID, request.GetName(), secretType)
	if err != nil && errors.Is(plainstorage.ErrSecretNotFound, err) {
		s.logger.Info("Cannot find secret by name", zap.String("login", user.Login), zap.String("secret_name", request.GetName()))

		return nil, status.Errorf(codes.NotFound, "secret '%s' not found", request.GetName())
	} else if err != nil {
		s.logger.Error("Error while get secret", zap.Error(err), zap.String("login", user.Login), zap.String("secret_name", request.GetName()))

		return nil, status.Error(codes.Internal, "internal error while get secret")
	}

	return &pb.SecretGetResponse{
		Secret: &pb.Secret{
			SecretType:      request.GetSecretType(),
			Name:            secret.Metadata.Name,
			CreateTimestamp: secret.Metadata.Created.Unix(),
			UpdateTimestamp: secret.Metadata.Updated.Unix(),
			Content:         secret.Data,
		},
	}, nil
}

func (s *Server) SecretUpdate(ctx context.Context, request *pb.SecretUpdateRequest) (*pb.SecretUpdateResponse, error) {
	userCtxVal := ctx.Value(UserContextKey)
	user := userCtxVal.(*plainstorage.User)

	s.logger.Info("User try update plain secret", zap.String("login", user.Login), zap.String("secret_name", request.GetName()), zap.String("secret_type", request.GetSecretType().String()))

	secretType, err := translateGRPCSecretTypeToSecretType(request.GetSecretType())
	if err != nil {
		s.logger.Error("User send invalid secret type", zap.String("login", user.Login))

		return nil, status.Error(codes.InvalidArgument, "invalid secret type got")
	}

	if secretType == plainstorage.SecretTypeMedia {
		s.logger.Error("User try to get media secret as plain", zap.String("login", user.Login))

		return nil, status.Error(codes.InvalidArgument, "cannot get media secret to plain storage")
	}

	err = s.plainStorage.UpdatePlainSecretByName(ctx, user.UUID, request.GetName(), request.GetContent())
	if err != nil {
		s.logger.Error("Cannot update plain secret", zap.Error(err), zap.String("login", user.Login))

		return nil, status.Error(codes.Internal, "cannot update plain secret")
	}

	return &pb.SecretUpdateResponse{}, nil
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
