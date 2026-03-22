package secret

import (
	"context"
	"errors"
	"fmt"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	"github.com/jackc/pgx/v5"
)

func (repo *SecretRepository) GetWithVersionsList(ctx context.Context, userID uint, id uint) (models.Secret, error) {
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

	rows, err := repo.DBPool.Query(ctx, "SELECT id, created_at FROM "+models.SecretVersionTable+" WHERE secret_id = $1 ORDER BY created_at DESC", id)
	if err != nil {
		return models.Secret{}, fmt.Errorf("could not get secret %d versions: %w", id, err)
	}
	defer rows.Close()

	for rows.Next() {
		secretVersion := models.SecretVersion{
			SecretID: id,
		}
		err = rows.Scan(&secretVersion.ID, &secretVersion.CreatedAt)
		if err != nil {
			return models.Secret{}, fmt.Errorf("could not scan secret %d version row: %w", id, err)
		}

		secret.Versions = append(secret.Versions, secretVersion)
	}

	err = rows.Err()
	if err != nil {
		return models.Secret{}, fmt.Errorf("could not get secret %d version rows: %w", id, err)
	}

	return secret, nil
}
