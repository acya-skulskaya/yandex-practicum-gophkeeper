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

// Show Will return secret's data (if it's not of type binary) for a specified version and a list of versions
func (s SecretServer) Show(ctx context.Context, in *pb.SecretShowRequest) (*pb.SecretResponse, error) {
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

	getLatestVersion := false
	var uintVersionID uint64
	if versionID == "0" {
		getLatestVersion = true
	} else if versionID != "" {
		uintVersionID, err = strconv.ParseUint(versionID, 10, 64)
		if err != nil {
			//nolint:wrapcheck
			return nil, status.Errorf(codes.Internal, "could not convert version id to uint: %v", err)
		}
	}

	secret, err := s.ServiceSecret.Get(ctx, userID, uint(uintID), uint(uintVersionID), getLatestVersion)
	if err != nil {
		//nolint:wrapcheck
		return nil, status.Errorf(codes.Internal, "could not get secret: %v", err)
	}

	versions := make([]*pb.SecretVersion, len(secret.Versions))
	for i, version := range secret.Versions {
		pbVersion := &pb.SecretVersion{}
		versionIDStr := strconv.FormatUint(uint64(version.ID), 10)
		pbVersion.SetVersionId(versionIDStr)
		pbVersion.SetCreatedAt(timestamppb.New(*version.CreatedAt))
		if version.Data != "" {
			pbSecretData := &pb.SecretData{}
			pbSecretData.SetText(version.Data)
			pbVersion.SetData(pbSecretData)
		}
		versions[i] = pbVersion
	}

	secretData := &pb.SecretData{}
	if len(secret.Versions) > 0 {
		secretData = versions[0].GetData()
	}

	versionIDStr := strconv.FormatUint(uint64(secret.Versions[0].ID), 10)
	secretType := models.MatchSecretTypeWithPBEnum(secret.Type)

	response := pb.SecretResponse_builder{
		Id:        &id,
		Type:      &secretType,
		Name:      proto.String(secret.Name),
		VersionId: &versionIDStr,
		Data:      secretData,
		CreatedAt: timestamppb.New(*secret.CreatedAt),
		UpdatedAt: timestamppb.New(*secret.UpdatedAt),
		Versions:  versions,
	}

	return response.Build(), nil
}
