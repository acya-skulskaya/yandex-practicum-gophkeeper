package secret

import (
	"context"
	"fmt"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
)

func (repo *SecretRepository) ExistsVersion(ctx context.Context, versionID uint) (bool, error) {
	var exists bool

	sql := "SELECT EXISTS(SELECT 1 FROM " + models.SecretVersionTable + " WHERE id = $1)"
	err := repo.DBPool.QueryRow(ctx, sql, versionID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("could not check that secret version %d exists: %w", versionID, err)
	}

	return exists, nil
}
