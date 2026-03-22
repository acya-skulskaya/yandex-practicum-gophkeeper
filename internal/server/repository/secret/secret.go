package secret

import (
	"context"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SecretRepositoryInterface interface {
	Create(ctx context.Context, userID uint, name string, secretType string, text string, data []byte) (models.Secret, error)
	Update(ctx context.Context, userID uint, id uint, name string, text string, data []byte) (models.Secret, error)
	SaveBinary(ctx context.Context, userID uint, secretVersionID uint, data []byte) error
	GetWithVersionsList(ctx context.Context, userID uint, id uint) (models.Secret, error)
	GetWithVersion(ctx context.Context, userID uint, id uint, versionID uint) (models.Secret, error)
	GetLatestVersionID(ctx context.Context, id uint) (uint, error)
	Delete(ctx context.Context, id uint) error
	DeleteVersion(ctx context.Context, versionID uint) error
	Exists(ctx context.Context, userID uint, id uint) (bool, error)
	ExistsVersion(ctx context.Context, versionID uint) (bool, error)
	VersionsCount(ctx context.Context, id uint) (uint, error)
	List(ctx context.Context, userID uint) ([]models.Secret, error)
	GetFilePath(userID uint, versionID uint) (string, error)
}

type SecretRepository struct {
	DBPool   *pgxpool.Pool
	FilesDir string
}

func NewSecretRepository(dbPool *pgxpool.Pool, dir string) *SecretRepository {
	return &SecretRepository{
		DBPool:   dbPool,
		FilesDir: dir,
	}
}
