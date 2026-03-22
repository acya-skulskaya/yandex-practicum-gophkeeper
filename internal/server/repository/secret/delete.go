package secret

import (
	"context"
	"fmt"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
)

func (repo *SecretRepository) Delete(ctx context.Context, id uint) error {
	tx, err := repo.DBPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not start transaction to delete secret %d: %w", id, err)
	}
	//nolint:errcheck // ignore err
	defer tx.Rollback(ctx)

	sql := "DELETE FROM " + models.SecretTable + " WHERE id = $1"
	_, err = tx.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("could not delete secret %d: %w", id, err)
	}

	sql = "DELETE FROM " + models.SecretVersionTable + " WHERE secret_id = $1"
	_, err = tx.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("could not delete secret %d versions: %w", id, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}
