package secret

import (
	"context"
	"fmt"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
)

func (repo *SecretRepository) DeleteVersion(ctx context.Context, id uint) error {
	sql := "DELETE FROM " + models.SecretVersionTable + " WHERE id = $1"
	rows, err := repo.DBPool.Query(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("could not delete secret version: %w", err)
	}
	defer rows.Close()

	return nil
}
