package secret

import (
	"context"
	"fmt"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
)

func (s *Service) Store(ctx context.Context, userID uint, name string, secretType string, data []byte, text string) (models.Secret, error) {
	encryptedText, err := s.Crypto.Encrypt([]byte(text))
	if err != nil {
		return models.Secret{}, fmt.Errorf("could not encrypt data %w", err)
	}
	hexText := s.Crypto.ToHexString(encryptedText)

	encryptedData, err := s.Crypto.Encrypt(data)
	if err != nil {
		return models.Secret{}, fmt.Errorf("could not encrypt data %w", err)
	}

	secret, err := s.Repo.Create(ctx, userID, name, secretType, hexText, encryptedData)
	if err != nil {
		return models.Secret{}, fmt.Errorf("could not create secret: %w", err)
	}

	return secret, nil
}
