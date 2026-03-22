package secret

import (
	"context"
	"fmt"
	"strconv"

	pb "github.com/acya-skulskaya/yandex-practicum-gophkeeper/api/gophkeeper"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/auth"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/config"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) Show(ctx context.Context, secretID string, secretVersionID string) (models.Secret, error) {
	request := pb.SecretShowRequest_builder{
		Id:        proto.String(secretID),
		VersionId: proto.String(secretVersionID),
	}.Build()

	token, err := auth.GetSessionToken()
	if err != nil {
		return models.Secret{}, fmt.Errorf("could not get session token: %w", err)
	}

	ctx = metadata.AppendToOutgoingContext(ctx, config.MetadataAuthKeyName, token)

	resp, err := h.Client.Show(ctx, request)
	if err != nil {
		return models.Secret{}, fmt.Errorf("could not handle show: %w", err)
	}

	pbVersions := resp.GetVersions()
	versions := make([]models.SecretVersion, 0, len(pbVersions))
	for _, pbVersion := range pbVersions {
		createdAt := (pbVersion.GetCreatedAt()).AsTime()
		uintID, err := strconv.ParseUint(pbVersion.GetVersionId(), 10, 64)
		if err != nil {
			return models.Secret{}, fmt.Errorf("could not get uint secret version ID: %w", err)
		}
		data := ""
		if secretVersionID != "" {
			data = resp.GetData().GetText()
		}
		version := models.SecretVersion{
			CreatedAt: &createdAt,
			Data:      data,
			SecretID:  0,
			ID:        uint(uintID),
		}

		versions = append(versions, version)
	}

	updatedAt := (resp.GetUpdatedAt()).AsTime()
	createdAt := (resp.GetCreatedAt()).AsTime()
	uintID, err := strconv.ParseUint(resp.GetId(), 10, 64)
	if err != nil {
		return models.Secret{}, fmt.Errorf("could not get uint secret ID: %w", err)
	}
	secret := models.Secret{
		UpdatedAt: &updatedAt,
		CreatedAt: &createdAt,
		Name:      resp.GetName(),
		Type:      models.MatchPBEnumWIthSecretType(resp.GetType().String()),
		ID:        uint(uintID),
		Versions:  versions,
	}

	return secret, nil
}
