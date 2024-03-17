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
	"io"
)

func (s *Server) UploadMediaSecret(stream pb.KeeperService_UploadMediaSecretServer) error {
	req, err := stream.Recv()

	if err != nil {
		s.logger.Error("Client media upload error", zap.Error(err))

		return status.Error(codes.DataLoss, "cannot recover media metadata package")
	}

	metadata := req.GetMetadata()
	if metadata == nil {
		s.logger.Error("Client sends empty media metadata by first package")

		return status.Error(codes.InvalidArgument, "first argument must be media metadata")
	}

	userCtxVal := stream.Context().Value(UserContextKey)
	user := userCtxVal.(*plainstorage.User)

	s.logger.Info("User sends new media", zap.String("login", user.Login), zap.String("filename", metadata.Name))

	dbMetadata, err := s.plainStorage.AddSecretMetadata(stream.Context(), user.UUID, metadata.Name, plainstorage.SecretTypeMedia)
	if err != nil {
		s.logger.Error("Error while save metadata of secret", zap.Error(err))

		return status.Error(codes.Internal, "cannot save media metadata")
	}

	upload, err := s.mediaStorage.StartUpload(stream.Context(), dbMetadata.UUID)
	if err != nil {
		s.logger.Error("Cannot start media upload", zap.Error(err))
	}

	cancelUpload := func(err error) {
		s.logger.Error("Error while upload media", zap.Error(err))
		err = upload.Abort(context.TODO())
		if err != nil {
			s.logger.Error("Error while abort media upload", zap.Error(err))
		}

		err = s.plainStorage.RemoveSecretByUUID(context.TODO(), dbMetadata.UUID)
		if err != nil {
			s.logger.Error("Error while remove secret metadata", zap.Error(err))
		}
	}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			cancelUpload(err)

			return status.Error(codes.DataLoss, "error while receive media data")
		}

		data := req.GetData()
		if data == nil {
			cancelUpload(errors.New("user sends no data"))

			return status.Error(codes.InvalidArgument, "packages must contain data after metadata package")
		}

		err = upload.Upload(stream.Context(), data.Chunk)
		if err != nil {
			cancelUpload(err)

			return status.Error(codes.Internal, "server cannot save data-chunk")
		}
	}

	err = upload.Complete(stream.Context())
	if err != nil {
		cancelUpload(fmt.Errorf("error while complete upload user media: %w", err))

		return status.Error(codes.Internal, "server cannot complete data upload")
	}

	return stream.SendAndClose(&pb.UploadMediaSecretResponse{
		Uuid: dbMetadata.UUID,
		Name: dbMetadata.Name,
	})
}
