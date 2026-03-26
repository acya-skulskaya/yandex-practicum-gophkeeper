package secret

import (
	"context"

	pb "github.com/acya-skulskaya/yandex-practicum-gophkeeper/api/gophkeeper"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/logger"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	authService "github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/service/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s SecretServer) Store(ctx context.Context, in *pb.SecretStoreRequest) (*pb.SecretResponse, error) {
	userID, ok := ctx.Value(authService.AuthContextKey(authService.AuthContextKeyUserID)).(uint)
	if !ok {
		logger.Log.Debug("could not get userID from context")
		//nolint:wrapcheck
		return nil, status.Error(codes.Unauthenticated, "could not get userID from context")
	}

	name := in.GetName()
	secretType := models.MatchPBEnumWIthSecretType(in.GetType().String())
	data := in.GetData().GetData()
	text := in.GetData().GetText()

	secret, err := s.ServiceSecret.Store(ctx, userID, name, secretType, data, text)
	if err != nil {
		//nolint:wrapcheck
		return nil, status.Errorf(codes.Internal, "could not store secret: %v", err)
	}

	var versions []*pb.SecretVersion
	version := &pb.SecretVersion{}
	versionIDStr := secret.Versions[0].GetStrID()
	version.SetVersionId(versionIDStr)
	version.SetCreatedAt(timestamppb.New(*secret.Versions[0].CreatedAt))
	versions = append(versions, version)

	secretID := secret.GetStrID()

	response := pb.SecretResponse_builder{
		Id:        &secretID,
		Type:      in.GetType().Enum(),
		Name:      proto.String(secret.Name),
		VersionId: &versionIDStr,
		CreatedAt: timestamppb.New(*secret.CreatedAt),
		UpdatedAt: timestamppb.New(*secret.UpdatedAt),
		Versions:  versions,
	}

	return response.Build(), nil
}
