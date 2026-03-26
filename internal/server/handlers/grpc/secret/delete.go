package secret

import (
	"context"
	"errors"
	"strconv"

	pb "github.com/acya-skulskaya/yandex-practicum-gophkeeper/api/gophkeeper"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/logger"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/repository/secret"
	authService "github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/service/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s SecretServer) Delete(ctx context.Context, in *pb.SecretDeleteRequest) (*pb.SecretDeleteResponse, error) {
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

	err = s.ServiceSecret.Delete(ctx, userID, uint(uintID), uint(uintVersionID))
	if err != nil {
		if errors.Is(err, secret.ErrSecretHasOnlyOneOrLessVersions) {
			//nolint:wrapcheck
			return nil, status.Error(codes.Canceled, "could not delete secret's version: secret has only one version, you should delete it fully")
		}
		//nolint:wrapcheck
		return nil, status.Errorf(codes.Internal, "could not delete secret: %v", err)
	}

	response := pb.SecretDeleteResponse_builder{}

	return response.Build(), nil
}
