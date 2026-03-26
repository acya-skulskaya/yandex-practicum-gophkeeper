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
)

func (h *Handler) List(ctx context.Context) ([]models.Secret, error) {
	request := pb.SyncRequest_builder{}.Build()

	token, err := auth.GetSessionToken()
	if err != nil {
		return []models.Secret{}, fmt.Errorf("could not get session token: %w", err)
	}

	ctx = metadata.AppendToOutgoingContext(ctx, config.MetadataAuthKeyName, token)

	resp, err := h.Client.Sync(ctx, request)
	if err != nil {
		return []models.Secret{}, fmt.Errorf("could not handle store: %w", err)
	}

	pbSecrets := resp.GetSecrets()
	secrets := make([]models.Secret, 0, len(pbSecrets))
	//var secrets []models.Secret
	for _, pbSecret := range pbSecrets {
		updatedAt := (pbSecret.GetUpdatedAt()).AsTime()
		createdAt := (pbSecret.GetCreatedAt()).AsTime()
		uintID, err := strconv.ParseUint(pbSecret.GetId(), 10, 64)
		if err != nil {
			return []models.Secret{}, fmt.Errorf("could not get uint secret ID: %w", err)
		}
		secret := models.Secret{
			UpdatedAt: &updatedAt,
			CreatedAt: &createdAt,
			Name:      pbSecret.GetName(),
			Type:      models.MatchPBEnumWIthSecretType(pbSecret.GetType().String()),
			ID:        uint(uintID),
		}

		pbVersions := pbSecret.GetVersions()
		versions := make([]models.SecretVersion, 0, len(pbVersions))
		for _, pbVersion := range pbVersions {
			createdAt = (pbVersion.GetCreatedAt()).AsTime()
			uintID, err = strconv.ParseUint(pbVersion.GetVersionId(), 10, 64)
			if err != nil {
				return []models.Secret{}, fmt.Errorf("could not get uint secret version ID: %w", err)
			}
			version := models.SecretVersion{
				CreatedAt: &createdAt,
				ID:        uint(uintID),
			}

			versions = append(versions, version)
		}

		secret.Versions = versions

		secrets = append(secrets, secret)
	}

	return secrets, nil
}
