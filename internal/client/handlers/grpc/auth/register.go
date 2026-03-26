package auth

import (
	"context"
	"fmt"

	pb "github.com/acya-skulskaya/yandex-practicum-gophkeeper/api/gophkeeper"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) Register(ctx context.Context, login string, password string) (string, error) {
	request := pb.AuthRequest_builder{
		Username: proto.String(login),
		Password: proto.String(password),
	}.Build()

	resp, err := h.Client.Register(ctx, request)
	if err != nil {
		return "", fmt.Errorf("could not handle register: %w", err)
	}

	token := resp.GetToken()

	return token, nil
}
