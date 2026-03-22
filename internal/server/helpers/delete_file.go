package helpers

import (
	"os"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/logger"
	"go.uber.org/zap"
)

func DeleteFile(file string) {
	go func(file string) {
		err := os.Remove(file)
		if err != nil {
			logger.Log.Error("Cannot delete file: "+file, zap.Error(err))
		}
		logger.Log.Debug("File was removed: " + file)
	}(file)
}
