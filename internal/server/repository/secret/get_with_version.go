package secret

import (
	"context"
	"errors"
	"fmt"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	"github.com/jackc/pgx/v5"
)

func (repo *SecretRepository) GetWithVersion(ctx context.Context, userID uint, id uint, versionID uint) (models.Secret, error) {
	secret := models.Secret{
		UserID: userID,
		ID:     id,
	}

	sql := "SELECT type, name, created_at, updated_at FROM " + models.SecretTable + " WHERE id = $1 AND user_id = $2"
	row := repo.DBPool.QueryRow(ctx, sql, id, userID)
	err := row.Scan(&secret.Type, &secret.Name, &secret.CreatedAt, &secret.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Secret{}, ErrSecretNotFound
	} else if err != nil {
		return models.Secret{}, fmt.Errorf("could not get secret %d: %w", id, err)
	}

	secretVersion := models.SecretVersion{
		SecretID: id,
		ID:       versionID,
	}

	if versionID == 0 {
		sql = "SELECT id, data, created_at FROM " + models.SecretVersionTable + " WHERE secret_id = $1 ORDER BY ID DESC LIMIT 1"
		row = repo.DBPool.QueryRow(ctx, sql, id)
		err = row.Scan(&secretVersion.ID, &secretVersion.Data, &secretVersion.CreatedAt)
		if err != nil {
			return models.Secret{}, fmt.Errorf("could not get latest secret %d version: %w", id, err)
		}
	} else {
		sql = "SELECT data, created_at FROM " + models.SecretVersionTable + " WHERE id = $1 AND secret_id = $2"
		row = repo.DBPool.QueryRow(ctx, sql, versionID, id)
		err = row.Scan(&secretVersion.Data, &secretVersion.CreatedAt)
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Secret{}, ErrSecretVersionNotFound
		} else if err != nil {
			return models.Secret{}, fmt.Errorf("could not get secret %d version: %w", id, err)
		}
	}

	secret.Versions = []models.SecretVersion{
		secretVersion,
	}

	return secret, nil
}
