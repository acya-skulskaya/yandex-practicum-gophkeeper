package secret

import (
	"context"
	"fmt"

	pb "github.com/acya-skulskaya/yandex-practicum-gophkeeper/api/gophkeeper"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/auth"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/config"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) Delete(ctx context.Context, secretID string, secretVersionID string) error {
	request := pb.SecretDeleteRequest_builder{
		Id:        proto.String(secretID),
		VersionId: proto.String(secretVersionID),
	}.Build()

	token, err := auth.GetSessionToken()
	if err != nil {
		return fmt.Errorf("could not get session token: %w", err)
	}

	ctx = metadata.AppendToOutgoingContext(ctx, config.MetadataAuthKeyName, token)

	_, err = h.Client.Delete(ctx, request)
	if err != nil {
		return fmt.Errorf("could not handle store: %w", err)
	}

	return nil
}
