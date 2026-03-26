package secret

import (
	"context"
	"strconv"

	pb "github.com/acya-skulskaya/yandex-practicum-gophkeeper/api/gophkeeper"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/logger"
	authService "github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/service/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s SecretServer) Download(ctx context.Context, in *pb.SecretShowRequest) (*pb.SecretDataResponse, error) {
	userID, ok := ctx.Value(authService.AuthContextKey(authService.AuthContextKeyUserID)).(uint)
	if !ok {
		logger.Log.Debug("could not get userID from context")
		//nolint:wrapcheck
		return nil, status.Error(codes.Unauthenticated, "could not get userID from context")
	}

	id := in.GetId()
	uintID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		//nolint:wrapcheck
		return nil, status.Errorf(codes.Internal, "could not convert id to uint: %v", err)
	}
	versionID := in.GetVersionId()
	var uintVersionID uint64
	if versionID != "" {
		uintVersionID, err = strconv.ParseUint(versionID, 10, 64)
		if err != nil {
			//nolint:wrapcheck
			return nil, status.Errorf(codes.Internal, "could not convert version id to uint: %v", err)
		}
	}

	secret, binary, err := s.ServiceSecret.Download(ctx, userID, uint(uintID), uint(uintVersionID))
	if err != nil {
		//nolint:wrapcheck
		return nil, status.Errorf(codes.Internal, "could not get secret: %v", err)
	}

	data := &pb.SecretData{}
	data.SetData(binary)
	data.SetText(secret.Versions[0].Data)
	versionIDStr := secret.Versions[0].GetStrID()

	response := pb.SecretDataResponse_builder{
		Type:      nil,
		VersionId: &versionIDStr,
		Data:      data,
		CreatedAt: nil,
	}

	return response.Build(), nil
}
