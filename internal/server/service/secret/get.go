package secret

import (
	"context"
	"fmt"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
)

func (s *Service) Get(ctx context.Context, userID uint, id uint, versionID uint, getLatestVersion bool) (models.Secret, error) {
	if versionID == 0 && !getLatestVersion {
		secret, err := s.Repo.GetWithVersionsList(ctx, userID, id)
		if err != nil {
			return models.Secret{}, fmt.Errorf("could not get secret %d with versions list: %w", id, err)
		}
		return secret, nil
	} else {
		secret, err := s.Repo.GetWithVersion(ctx, userID, id, versionID)
		if err != nil {
			return models.Secret{}, fmt.Errorf("could not get secret %d with version %d: %w", id, versionID, err)
		}

		unhexed, err := s.Crypto.FromHexString(secret.Versions[0].Data)
		if err != nil {
			return models.Secret{}, fmt.Errorf("could not decode version %d from hex: %w", id, err)
		}
		decryptedText, err := s.Crypto.Decrypt(unhexed)
		if err != nil {
			return models.Secret{}, fmt.Errorf("could not decrypt secret data: %w", err)
		}
		secret.Versions[0].Data = string(decryptedText)

		return secret, nil
	}
}
