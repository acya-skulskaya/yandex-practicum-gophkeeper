package secret

import (
	pb "github.com/acya-skulskaya/yandex-practicum-gophkeeper/api/gophkeeper"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/service/secret"
)

type SecretServer struct {
	pb.UnimplementedSecretServiceServer

	ServiceSecret *secret.Service
}
