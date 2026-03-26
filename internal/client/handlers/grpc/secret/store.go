package secret

import (
	"context"
	"fmt"
	"os"

	pb "github.com/acya-skulskaya/yandex-practicum-gophkeeper/api/gophkeeper"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/auth"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/config"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) Store(ctx context.Context, name string, secretType string, secretMetadata []byte, binaryFilePath string) (string, error) {
	pbSecretType := models.MatchSecretTypeWithPBEnum(secretType)
	data := &pb.SecretData{}
	data.SetText(string(secretMetadata))
	if binaryFilePath != "" {
		fileData, err := os.ReadFile(binaryFilePath)
		if err != nil {
			return "", fmt.Errorf("unable to read file %s: %w", binaryFilePath, err)
		}
		data.SetData(fileData)
	}

	request := pb.SecretStoreRequest_builder{
		Type: &pbSecretType,
		Name: proto.String(name),
		Data: data,
	}.Build()

	token, err := auth.GetSessionToken()
	if err != nil {
		return "", fmt.Errorf("could not get session token: %w", err)
	}

	ctx = metadata.AppendToOutgoingContext(ctx, config.MetadataAuthKeyName, token)

	resp, err := h.Client.Store(ctx, request)
	if err != nil {
		return "", fmt.Errorf("could not handle store: %w", err)
	}

	id := resp.GetId()

	return id, nil
}
