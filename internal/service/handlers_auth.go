package service

import (
	"context"
	"errors"
	pb "github.com/nessai1/gophkeeper/api/proto"
	"github.com/nessai1/gophkeeper/internal/service/plainstorage"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Ping(ctx context.Context, request *pb.PingRequest) (*pb.PingResponse, error) {
	s.logger.Info("Got ping")

	resp := pb.PingResponse{}
	resp.Answer = "pong!"

	return &resp, nil
}

func (s *Server) Register(ctx context.Context, request *pb.UserCredentialsRequest) (*pb.UserCredentialsResponse, error) {
	user, err := s.plainStorage.CreateUser(ctx, request.Login, hashPassword(request.Password))
	if err != nil && errors.Is(plainstorage.ErrEntityAlreadyExists, err) {
		s.logger.Info("User try to register existing login", zap.String("login", request.Login))

		return nil, status.Error(codes.AlreadyExists, "User already exists")
	} else if err != nil {
		s.logger.Error("Got unexpected error while create user", zap.Error(err))

		return nil, status.Error(codes.Internal, "Unexpected error while create user")
	}

	sign, err := generateSign(user.UUID, s.config.SecretToken)
	if err != nil {
		s.logger.Error("Cannot generate sign for user (register)", zap.Error(err), zap.String("user_uuid", user.UUID))

		return nil, status.Error(codes.Internal, "Error while generate sign")
	}

	s.logger.Info("New user registered", zap.String("login", user.Login), zap.String("uuid", user.UUID))

	return &pb.UserCredentialsResponse{Token: sign}, nil
}

func (s *Server) Login(ctx context.Context, request *pb.UserCredentialsRequest) (*pb.UserCredentialsResponse, error) {
	user, err := s.plainStorage.GetUserByLogin(ctx, request.Login)
	if err != nil && !errors.Is(plainstorage.ErrEntityNotFound, err) {
		s.logger.Error("Error while get user for log-in", zap.Error(err))

		return nil, status.Error(codes.Internal, "Unexpected error while get user for log-in")
	} else if errors.Is(plainstorage.ErrEntityNotFound, err) {
		s.logger.Info("User try to log-in not exists account", zap.String("login", request.Login))

		return nil, status.Error(codes.NotFound, "Account doesn't exists")
	}

	if user.PasswordHash != hashPassword(request.Password) {
		s.logger.Info("User send wrong password for login", zap.String("login", request.Login))

		return nil, status.Error(codes.InvalidArgument, "Incorrect password")
	}

	sign, err := generateSign(user.UUID, s.config.SecretToken)
	if err != nil {
		s.logger.Error("Cannot generate sign for user (login)", zap.Error(err), zap.String("user_uuid", user.UUID))

		return nil, status.Error(codes.Internal, "Error while generate sign")
	}

	s.logger.Info("User log-in", zap.String("login", user.Login), zap.String("uuid", user.UUID))

	return &pb.UserCredentialsResponse{Token: sign}, nil
}
