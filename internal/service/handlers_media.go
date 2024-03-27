package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/nessai1/gophkeeper/api/proto"
	"github.com/nessai1/gophkeeper/internal/service/plainstorage"
	"github.com/nessai1/gophkeeper/pkg/bytesize"
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

	mediaUUID := uuid.New().String()
	var completeUpload func() error
	secret, getSecretErr := s.plainStorage.GetUserSecretByName(stream.Context(), user.UUID, metadata.Name, plainstorage.SecretTypeMedia)
	if getSecretErr != nil && !errors.Is(plainstorage.ErrEntityNotFound, getSecretErr) {
		s.logger.Error("Cannot check existing media secret", zap.Error(getSecretErr), zap.String("login", user.Login), zap.String("filename", metadata.Name))

		return status.Errorf(codes.Internal, "Cannot check existing media secret '%s'", metadata.Name)
	}

	if metadata.Overwrite {
		if getSecretErr != nil && errors.Is(plainstorage.ErrEntityNotFound, getSecretErr) {
			s.logger.Info("Cannot overwrite not existing media", zap.String("login", user.Login), zap.String("filename", metadata.Name))

			return status.Errorf(codes.NotFound, "Cannot overwrite not existing media '%s'", metadata.Name)
		}

		s.logger.Info("Start remove old media secret", zap.String("filename", metadata.Name))

		completeUpload = func() error {
			err := s.plainStorage.UpdateSecretMetadataUUID(stream.Context(), user.UUID, secret.Metadata.UUID, mediaUUID, plainstorage.SecretTypeMedia)
			if err != nil {
				s.logger.Error("Cannot update plain storage to new media UUID", zap.Error(err), zap.String("filename", metadata.Name), zap.String("login", user.Login))

				return status.Errorf(codes.Internal, "cannot update plain storage for new media UUID")
			}

			go func() {
				err := s.mediaStorage.Delete(context.TODO(), secret.Metadata.UUID)
				if err != nil {
					s.logger.Error("Cannot delete old media", zap.Error(err), zap.String("filename", metadata.Name), zap.String("login", user.Login))
				}
			}()

			return nil
		}
	} else {
		if !errors.Is(plainstorage.ErrEntityNotFound, getSecretErr) {
			s.logger.Info("User try to create existing media", zap.String("login", user.Login), zap.String("filename", metadata.Name))

			return status.Errorf(codes.AlreadyExists, "Cannot create existing media secret '%s'", metadata.Name)
		}

		completeUpload = func() error {
			_, err := s.plainStorage.AddSecretMetadata(stream.Context(), user.UUID, mediaUUID, metadata.Name, plainstorage.SecretTypeMedia)
			if err != nil {
				return status.Error(codes.Internal, "Cannot add secret media to plain storage")
			}

			return nil
		}
	}

	upload, err := s.mediaStorage.StartUpload(stream.Context(), mediaUUID)
	if err != nil {
		s.logger.Error("Cannot start media upload", zap.Error(err))
	}

	cancelUpload := func(err error) {
		s.logger.Error("Error while upload media", zap.Error(err))
		err = upload.Abort(context.TODO())
		if err != nil {
			s.logger.Error("Error while abort media upload", zap.Error(err))
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

		s.logger.Debug("Send media part to storage", zap.String("content_size", fmt.Sprintf("%f KB", float64(len(data.Chunk))/float64(bytesize.KB))))
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

	err = completeUpload()
	if err != nil {
		return err
	}

	return stream.SendAndClose(&pb.UploadMediaSecretResponse{
		Uuid: mediaUUID,
		Name: metadata.Name,
	})
}

func (s *Server) DownloadMediaSecret(req *pb.DownloadMediaSecretRequest, stream pb.KeeperService_DownloadMediaSecretServer) error {
	userCtxVal := stream.Context().Value(UserContextKey)
	user := userCtxVal.(*plainstorage.User)

	s.logger.Info("User try to get media", zap.String("login", user.Login), zap.String("filename", req.SecretName))

	secret, err := s.plainStorage.GetUserSecretByName(stream.Context(), user.UUID, req.SecretName, plainstorage.SecretTypeMedia)
	if errors.Is(plainstorage.ErrEntityNotFound, err) {
		s.logger.Info("User try to get not existing media-secret", zap.String("login", user.UUID), zap.String("filaname", req.SecretName))

		return status.Errorf(codes.NotFound, "cannot found secret with name %s", req.SecretName)
	} else if err != nil {
		s.logger.Error("Cannot get secret for media download", zap.String("login", user.UUID), zap.String("filaname", req.SecretName), zap.Error(err))

		return status.Error(codes.Internal, "cannot get secret info from DB")
	}

	rc, err := s.mediaStorage.StartDownload(stream.Context(), secret.Metadata.UUID)
	if err != nil {
		s.logger.Error("Cannot start media download from storage", zap.String("media_uuid", secret.Metadata.UUID), zap.Error(err))

		return status.Error(codes.Internal, "cannot start media file download from storage")
	}

	var b []byte
	for {
		b = make([]byte, bytesize.MB)
		n, err := rc.Read(b)
		if errors.Is(io.EOF, err) {
			break
		}

		if n != int(bytesize.MB) {
			b = b[:n]
		}

		err = stream.Send(&pb.DownloadMediaSecretResponse{SecretPart: &pb.MediaSecret{Chunk: b}})
		if err != nil {
			s.logger.Error("Cannot send media part to client", zap.String("login", user.Login), zap.String("media_uuid", secret.Metadata.UUID), zap.Error(err))

			return status.Error(codes.DataLoss, "error while sending media data to client")
		}
	}

	return nil
}
