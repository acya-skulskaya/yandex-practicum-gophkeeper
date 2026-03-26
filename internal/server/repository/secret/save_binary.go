package secret

import (
	"bufio"
	"context"
	"fmt"
	"os"
)

func (repo *SecretRepository) SaveBinary(ctx context.Context, userID uint, secretVersionID uint, data []byte) error {
	filename, err := repo.GetFilePath(userID, secretVersionID)
	if err != nil {
		return fmt.Errorf("could not get file path: %w", err)
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", filename, err)
	}
	//nolint:errcheck // ignore err
	defer file.Close()

	writer := bufio.NewWriter(file)

	// записываем событие в буфер
	if _, err = writer.Write(data); err != nil {
		return fmt.Errorf("error writing file %s: %w", filename, err)
	}

	// записываем буфер в файл
	return writer.Flush()
}
