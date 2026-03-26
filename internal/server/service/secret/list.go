package secret

import (
	"context"
	"fmt"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
)

func (s *Service) List(ctx context.Context, userID uint) ([]models.Secret, error) {
	secrets, err := s.Repo.List(ctx, userID)
	if err != nil {
		return []models.Secret{}, fmt.Errorf("could not get a list of user's %d secrets: %w", userID, err)
	}

	return secrets, nil
}
