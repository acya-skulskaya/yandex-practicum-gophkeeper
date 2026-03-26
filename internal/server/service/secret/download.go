package secret

import (
	"context"
	"fmt"
	"os"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
)

func (s *Service) Download(ctx context.Context, userID uint, id uint, versionID uint) (models.Secret, []byte, error) {
	var secret models.Secret
	var err error

	// if version is not specified, latest version will be found
	secret, err = s.Repo.GetWithVersion(ctx, userID, id, versionID)
	if err != nil {
		return models.Secret{}, nil, fmt.Errorf("could not get secret %d with version %d: %w", id, versionID, err)
	}

	if secret.Type != models.SecretTypeBinary {
		return models.Secret{}, nil, fmt.Errorf("invalid secret type: %s", secret.Type)
	}

	unhexed, err := s.Crypto.FromHexString(secret.Versions[0].Data)
	if err != nil {
		return models.Secret{}, nil, fmt.Errorf("could not decode version data from hex: %w", err)
	}
	decryptedText, err := s.Crypto.Decrypt(unhexed)
	if err != nil {
		return models.Secret{}, nil, fmt.Errorf("could not decrypt secret data: %w", err)
	}

	secret.Versions[0].Data = string(decryptedText)

	// get file
	secretVersionID := secret.Versions[0].ID
	filename, err := s.Repo.GetFilePath(userID, secretVersionID)
	if err != nil {
		return models.Secret{}, nil, fmt.Errorf("could not get file path: %w", err)
	}
	//nolint:gosec // ignore
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0o666)
	if err != nil {
		return models.Secret{}, nil, fmt.Errorf("error opening file %s: %w", filename, err)
	}
	//nolint:errcheck // ignore err
	defer file.Close()

	fileinfo, _ := file.Stat()
	bs := make([]byte, fileinfo.Size())

	_, err = file.Read(bs)
	if err != nil {
		return models.Secret{}, nil, fmt.Errorf("error reading file %s: %w", filename, err)
	}

	decryptedData, err := s.Crypto.Decrypt(bs)
	if err != nil {
		return models.Secret{}, nil, fmt.Errorf("error decrypting file %s: %w", filename, err)
	}
	return secret, decryptedData, nil
}
