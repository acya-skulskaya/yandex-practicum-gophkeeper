package secret

import (
	"context"
	"fmt"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/logger"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/helpers"
	secretRepo "github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/repository/secret"
	"go.uber.org/zap"
)

func (s *Service) Delete(ctx context.Context, userID uint, id uint, versionID uint) error {
	exists, err := s.Repo.Exists(ctx, userID, id)
	if err != nil {
		return fmt.Errorf("could not check secret %d existance: %w", id, err)
	}
	if !exists {
		return secretRepo.ErrSecretNotFound
	}

	if versionID == 0 {
		secret, err := s.Repo.GetWithVersionsList(ctx, userID, id)
		if err != nil {
			return fmt.Errorf("could not get secret %d with versions list: %w", id, err)
		}

		if err := s.Repo.Delete(ctx, id); err != nil {
			return fmt.Errorf("could not get secret %d with versions list: %w", id, err)
		}

		go func(versions []models.SecretVersion) {
			for _, version := range versions {
				filename, err := s.Repo.GetFilePath(userID, version.ID)
				if err != nil {
					logger.Log.Error("could not get file path", zap.String("filename", filename), zap.Error(err))
				}
				helpers.DeleteFile(filename)
			}
		}(secret.Versions)
	} else {
		exists, err = s.Repo.ExistsVersion(ctx, versionID)
		if err != nil {
			return fmt.Errorf("could not check secret %d version %d existance: %w", id, versionID, err)
		}
		if !exists {
			return secretRepo.ErrSecretVersionNotFound
		}

		count, err := s.Repo.VersionsCount(ctx, id)
		if err != nil {
			return fmt.Errorf("could not check secret %d versions count: %w", id, err)
		}

		if count <= 1 {
			return secretRepo.ErrSecretHasOnlyOneOrLessVersions
		}

		if err := s.Repo.DeleteVersion(ctx, versionID); err != nil {
			return fmt.Errorf("could not get secret %d with version %d: %w", id, versionID, err)
		}

		filename, err := s.Repo.GetFilePath(userID, versionID)
		if err != nil {
			return fmt.Errorf("could not get file path for secret %d with version %d: %w", id, versionID, err)
		}
		helpers.DeleteFile(filename)
	}
	return nil
}
