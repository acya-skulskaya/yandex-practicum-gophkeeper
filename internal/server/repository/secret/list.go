package secret

import (
	"context"
	"fmt"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
)

func (repo *SecretRepository) List(ctx context.Context, userID uint) ([]models.Secret, error) {
	var secrets []models.Secret
	var secretsIDs []uint

	sql := "SELECT id, type, name, created_at, updated_at FROM " + models.SecretTable + " WHERE user_id = $1 ORDER BY created_at DESC"

	rows, err := repo.DBPool.Query(ctx, sql, userID)
	if err != nil {
		return []models.Secret{}, fmt.Errorf("could not get list of secretts for user %d: %w", userID, err)
	}
	defer rows.Close()
	for rows.Next() {
		secret := models.Secret{
			UserID: userID,
		}

		err = rows.Scan(&secret.ID, &secret.Type, &secret.Name, &secret.CreatedAt, &secret.UpdatedAt)
		if err != nil {
			return []models.Secret{}, fmt.Errorf("could not scan a row with secret to get user's %d secrets: %w", userID, err)
		}

		secrets = append(secrets, secret)
		secretsIDs = append(secretsIDs, secret.ID)
	}

	if len(secrets) == 0 {
		return []models.Secret{}, nil
	}

	rows, err = repo.DBPool.Query(ctx, "SELECT "+versionsIteratorSelectColumns+" FROM "+models.SecretVersionTable+" WHERE secret_id = ANY($1) ORDER BY created_at DESC", secretsIDs)
	if err != nil {
		return []models.Secret{}, fmt.Errorf("could not get list of secret versions for user %d: %w", userID, err)
	}
	defer rows.Close()

	secretVersions := make(map[uint][]models.SecretVersion)
	for secretVersion, err := range versionsIterator(rows) {
		if err != nil {
			if err != nil {
				return []models.Secret{}, fmt.Errorf("could not scan a row with secret versions to get user's %d secrets: %w", userID, err)
			}
			break // Stop on error
		}
		secretVersions[secretVersion.SecretID] = append(secretVersions[secretVersion.SecretID], secretVersion)
	}

	for i := range secrets {
		secrets[i].Versions = secretVersions[secrets[i].ID]
	}

	return secrets, nil
}
