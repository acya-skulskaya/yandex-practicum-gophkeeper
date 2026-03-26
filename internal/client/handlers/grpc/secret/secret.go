package secret

import (
	pb "github.com/acya-skulskaya/yandex-practicum-gophkeeper/api/gophkeeper"
)

type Handler struct {
	Client pb.SecretServiceClient
}

func New(client pb.SecretServiceClient) *Handler {
	return &Handler{
		Client: client,
	}
}
