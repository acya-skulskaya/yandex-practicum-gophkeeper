package secret

import (
	"context"
	"errors"
	"fmt"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	"github.com/jackc/pgx/v5"
)

func (repo *SecretRepository) GetLatestVersionID(ctx context.Context, id uint) (uint, error) {
	sql := "SELECT id FROM " + models.SecretVersionTable + " WHERE secret_id = $1 ORDER BY created_at DESC LIMIT 1"
	row := repo.DBPool.QueryRow(ctx, sql, id)
	var versionID uint
	err := row.Scan(&versionID)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrSecretNotFound
	} else if err != nil {
		return 0, fmt.Errorf("could not get secret %d version id: %w", id, err)
	}

	return versionID, nil
}
