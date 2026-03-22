package auth

import (
	pb "github.com/acya-skulskaya/yandex-practicum-gophkeeper/api/gophkeeper"
)

type Handler struct {
	Client pb.AuthServiceClient
}

func New(client pb.AuthServiceClient) *Handler {
	return &Handler{
		Client: client,
	}
}
