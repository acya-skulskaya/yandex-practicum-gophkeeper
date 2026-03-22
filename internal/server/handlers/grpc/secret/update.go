package secret

import (
	"context"
	"strconv"

	pb "github.com/acya-skulskaya/yandex-practicum-gophkeeper/api/gophkeeper"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/logger"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	authService "github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/service/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s SecretServer) Update(ctx context.Context, in *pb.SecretUpdateRequest) (*pb.SecretResponse, error) {
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
	name := in.GetName()
	data := in.GetData().GetData()
	text := in.GetData().GetText()

	secret, err := s.ServiceSecret.Update(ctx, userID, uint(uintID), name, data, text)
	if err != nil {
		//nolint:wrapcheck
		return nil, status.Errorf(codes.Internal, "could not update secret: %v", err)
	}

	var versions []*pb.SecretVersion
	version := &pb.SecretVersion{}
	versionIDStr := strconv.FormatUint(uint64(secret.Versions[0].ID), 10)
	version.SetVersionId(versionIDStr)
	version.SetCreatedAt(timestamppb.New(*secret.Versions[0].CreatedAt))
	versions = append(versions, version)

	secretType := models.MatchSecretTypeWithPBEnum(secret.Type)

	response := pb.SecretResponse_builder{
		Id:        &id,
		Type:      &secretType,
		Name:      proto.String(secret.Name),
		VersionId: &versionIDStr,
		CreatedAt: timestamppb.New(*secret.CreatedAt),
		UpdatedAt: timestamppb.New(*secret.UpdatedAt),
		Versions:  versions,
	}

	return response.Build(), nil
}
