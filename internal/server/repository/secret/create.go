package secret

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (repo *SecretRepository) Create(ctx context.Context, userID uint, name string, secretType string, text string, data []byte) (models.Secret, error) {
	tx, err := repo.DBPool.Begin(ctx)
	if err != nil {
		return models.Secret{}, fmt.Errorf("could not start transaction to create secret: %w", err)
	}
	//nolint:errcheck // ignore err
	defer tx.Rollback(ctx)

	sql := "INSERT INTO " + models.SecretTable + " (user_id, type, name, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id, created_at, updated_at"
	var id uint
	var createdAt time.Time
	var updatedAt time.Time
	err = tx.QueryRow(ctx, sql, userID, secretType, name).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return models.Secret{}, ErrNameAndTypeAlreadyExist
		} else {
			return models.Secret{}, fmt.Errorf("could not create secret: %w", err)
		}
	}

	secret := models.Secret{
		UpdatedAt: &updatedAt,
		CreatedAt: &createdAt,
		Name:      name,
		Type:      secretType,
		UserID:    userID,
		ID:        id,
	}

	sql = "INSERT INTO " + models.SecretVersionTable + " (secret_id, data, created_at) VALUES ($1, $2, NOW()) RETURNING id, created_at"
	var secretVersionID uint
	var secretVersionCreatedAt time.Time
	err = tx.QueryRow(ctx, sql, secret.ID, text).Scan(&secretVersionID, &secretVersionCreatedAt)
	if err != nil {
		return models.Secret{}, fmt.Errorf("could not create secret version: %w", err)
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

	if secretType == models.SecretTypeBinary {
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
