package auth

import (
	pb "github.com/acya-skulskaya/yandex-practicum-gophkeeper/api/gophkeeper"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/service/auth"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer

	ServiceAuth *auth.Service
}
