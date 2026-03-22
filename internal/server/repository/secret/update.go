package secret

import (
	"context"
	"fmt"
	"time"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
)

func (repo *SecretRepository) Update(ctx context.Context, userID uint, id uint, name string, text string, data []byte) (models.Secret, error) {
	tx, err := repo.DBPool.Begin(ctx)
	if err != nil {
		return models.Secret{}, fmt.Errorf("could not start transaction to update secret %d: %w", id, err)
	}
	//nolint:errcheck // ignore err
	defer tx.Rollback(ctx)

	secret := models.Secret{
		Name:   name,
		UserID: userID,
		ID:     id,
	}

	sql := "UPDATE " + models.SecretTable + " SET name = $1, updated_at = NOW() WHERE id = $2 AND user_id = $3  RETURNING type, created_at, updated_at"
	err = tx.QueryRow(ctx, sql, name, id, userID).Scan(&secret.Type, &secret.CreatedAt, &secret.UpdatedAt)
	if err != nil {
		return models.Secret{}, fmt.Errorf("could not update secret %d: %w", id, err)
	}

	sql = "INSERT INTO " + models.SecretVersionTable + " (secret_id, data, created_at) VALUES ($1, $2, NOW()) RETURNING id, created_at"
	var secretVersionID uint
	var secretVersionCreatedAt time.Time
	err = tx.QueryRow(ctx, sql, secret.ID, text).Scan(&secretVersionID, &secretVersionCreatedAt)
	if err != nil {
		return models.Secret{}, fmt.Errorf("could not create secret version to update secret %d: %w", id, err)
	}
	secretVersion := models.SecretVersion{
		CreatedAt: &secretVersionCreatedAt,
		SecretID:  secret.ID,
		Data:      text,
		ID:        secretVersionID,
	}

	secret.Versions = []models.SecretVersion{
		secretVersion,
	}

	if secret.Type == models.SecretTypeBinary {
		err = repo.SaveBinary(ctx, userID, secretVersion.ID, data)
		if err != nil {
			return models.Secret{}, fmt.Errorf("could not save binary: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return models.Secret{}, fmt.Errorf("could not commit transaction: %w", err)
	}

	return secret, nil
}
