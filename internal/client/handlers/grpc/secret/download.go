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

func (h *Handler) Download(ctx context.Context, secretID string, secretVersionID string) (string, string, []byte, error) {
	request := pb.SecretShowRequest_builder{
		Id:        proto.String(secretID),
		VersionId: proto.String(secretVersionID),
	}.Build()

	token, err := auth.GetSessionToken()
	if err != nil {
		return "", "", []byte{}, fmt.Errorf("could not get session token: %w", err)
	}

	ctx = metadata.AppendToOutgoingContext(ctx, config.MetadataAuthKeyName, token)

	resp, err := h.Client.Download(ctx, request)
	if err != nil {
		return "", "", []byte{}, fmt.Errorf("could not handle download: %w", err)
	}

	versionID := resp.GetVersionId()
	data := resp.GetData()

	return versionID, data.GetText(), data.GetData(), nil
}
