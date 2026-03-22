package secret

import (
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/crypto"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/repository/secret"
)

type Service struct {
	Repo      secret.SecretRepositoryInterface
	Crypto    *crypto.Crypto
	SecretKey string
}

func New(secretKey string, repo secret.SecretRepositoryInterface, crpt crypto.Crypto) *Service {
	return &Service{
		SecretKey: secretKey,
		Repo:      repo,
		Crypto:    &crpt,
	}
}
