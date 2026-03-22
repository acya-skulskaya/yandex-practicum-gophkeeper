package secret

import (
	"context"
	"fmt"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
)

func (repo *SecretRepository) Exists(ctx context.Context, userID uint, id uint) (bool, error) {
	var exists bool

	sql := "SELECT EXISTS(SELECT 1 FROM " + models.SecretTable + " WHERE id = $1 AND user_id = $2)"
	err := repo.DBPool.QueryRow(ctx, sql, id, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("could not check that secret %d exists: %w", id, err)
	}

	return exists, nil
}
