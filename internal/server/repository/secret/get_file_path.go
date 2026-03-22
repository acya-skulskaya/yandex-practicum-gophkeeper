package secret

import (
	"fmt"
	"os"
)

func (repo *SecretRepository) GetFilePath(userID uint, versionID uint) (string, error) {
	dir := fmt.Sprintf("%s%d", repo.FilesDir, userID)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		//nolint:gosec // ignore
		err2 := os.MkdirAll(dir, 0o766)
		if err2 != nil {
			return "", fmt.Errorf("error creating directory %s: %w", dir, err2)
		}
	}

	return fmt.Sprintf("%s%d/%d.bin", repo.FilesDir, userID, versionID), nil
}
