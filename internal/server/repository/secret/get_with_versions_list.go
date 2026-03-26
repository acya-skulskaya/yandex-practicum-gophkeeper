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

	rows, err := repo.DBPool.Query(ctx, "SELECT "+versionsIteratorSelectColumns+" FROM "+models.SecretVersionTable+" WHERE secret_id = $1 ORDER BY created_at DESC", id)
	if err != nil {
		return models.Secret{}, fmt.Errorf("could not get secret %d versions: %w", id, err)
	}
	defer rows.Close()

	for secretVersion, err := range versionsIterator(rows) {
		if err != nil {
			if err != nil {
				return models.Secret{}, fmt.Errorf("could not scan secret %d version row: %w", id, err)
			}
			break // Stop on error
		}
		secret.Versions = append(secret.Versions, secretVersion)
	}

	return secret, nil
}
