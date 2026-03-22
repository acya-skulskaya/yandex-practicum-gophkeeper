package secret

import (
	"context"
	"fmt"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
)

func (repo *SecretRepository) VersionsCount(ctx context.Context, id uint) (uint, error) {
	var count uint

	sql := "SELECT count(id) FROM " + models.SecretVersionTable + " WHERE secret_id = $1"
	err := repo.DBPool.QueryRow(ctx, sql, id).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("could not check number of secret %d versions: %w", id, err)
	}

	return count, nil
}
