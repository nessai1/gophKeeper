package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	tokenTTL = time.Hour * 24 * 30
)

const UserContextKey UserContextKeyType = "UserContext"

type UserContextKeyType string

type claims struct {
	jwt.RegisteredClaims
	UserUUID string
}

// ContextAuthKey type for store user UUID in request context
type ContextAuthKey string

func (s *Server) unaryAuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	s.logger.Debug("Got new unary request", zap.String("method", info.FullMethod))

	inWhiteList := false
	for _, val := range unauthorizedMethods {
		if val == info.FullMethod {
			inWhiteList = true
			break
		}
	}

	if inWhiteList {
		resp, err = handler(ctx, req)

		return resp, err
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s.logger.Debug("User request doesn't contain metadata for auth request", zap.String("method", info.FullMethod))

		return nil, status.Error(codes.Unauthenticated, "This method require auth metadata")
	}

	tokenArr := md.Get("jwt")
	if tokenArr == nil || len(tokenArr) == 0 {
		return nil, status.Error(codes.Unauthenticated, "This method require auth metadata (got empty jwt field in metadata)")
	}

	userUUID, err := fetchUUID(tokenArr[0], s.config.SecretToken)
	if err != nil {
		s.logger.Info("User sends invalid to parse jwt", zap.Error(err))

		return nil, status.Error(codes.Unauthenticated, "Cannot parse jwt token for authorize")
	}

	user, err := s.plainStorage.GetUserByUUID(ctx, userUUID)
	if err != nil {
		s.logger.Error("Cannot get user by UUID", zap.Error(err))

		return nil, status.Error(codes.Unauthenticated, "Cannot get user by given credentials")
	}

	ctx = context.WithValue(ctx, UserContextKey, user)
	resp, err = handler(ctx, req)

	return resp, err
}

func fetchUUID(sign, secretToken string) (string, error) {
	c := &claims{}
	_, err := jwt.ParseWithClaims(sign, c, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretToken), nil
	})

	if err != nil {
		return "", fmt.Errorf("cannot parse jwt: %w", err)
	}

	return c.UserUUID, nil
}

func (s *Server) streamAuthInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	s.logger.Info("Got new stream request", zap.String("method", info.FullMethod))
	return nil
}

func generateSign(UUID, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		},
		UserUUID: UUID,
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func hashPassword(pwd string) string {
	hash := sha256.Sum256([]byte(pwd + "myGophkeeperSaltlol"))

	return fmt.Sprintf("%x", hash)
}
