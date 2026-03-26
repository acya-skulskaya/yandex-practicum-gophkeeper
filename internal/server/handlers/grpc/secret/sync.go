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

func (s SecretServer) Sync(ctx context.Context, in *pb.SyncRequest) (*pb.SyncResponse, error) {
	userID, ok := ctx.Value(authService.AuthContextKey(authService.AuthContextKeyUserID)).(uint)
	if !ok {
		logger.Log.Debug("could not get userID from context")
		//nolint:wrapcheck
		return nil, status.Error(codes.Unauthenticated, "could not get userID from context")
	}

	secretsList, err := s.ServiceSecret.List(ctx, userID)
	if err != nil {
		//nolint:wrapcheck
		return nil, status.Errorf(codes.Internal, "could not get list of all user's secrets: %v", err)
	}

	secrets := make([]*pb.SecretResponse, len(secretsList))
	for i, secret := range secretsList {
		versions := make([]*pb.SecretVersion, len(secret.Versions))
		for i, version := range secret.Versions {
			pbVersion := &pb.SecretVersion{}
			pbVersion.SetVersionId(version.GetStrID())
			pbVersion.SetCreatedAt(timestamppb.New(*version.CreatedAt))
			versions[i] = pbVersion
		}

		sr := &pb.SecretResponse{}
		sr.SetId(secret.GetStrID())
		secretType := models.MatchSecretTypeWithPBEnum(secret.Type)
		sr.SetType(secretType)
		sr.SetName(*proto.String(secret.Name))
		if len(secret.Versions) > 0 {
			sr.SetVersionId(secret.Versions[0].GetStrID())
		}
		sr.SetCreatedAt(timestamppb.New(*secret.CreatedAt))
		sr.SetUpdatedAt(timestamppb.New(*secret.UpdatedAt))
		sr.SetVersions(versions)

		secrets[i] = sr
	}

	response := pb.SyncResponse_builder{
		Secrets: secrets,
	}

	return response.Build(), nil
}
