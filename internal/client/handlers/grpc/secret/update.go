package secret

import (
	"context"
	"fmt"
	"os"

	pb "github.com/acya-skulskaya/yandex-practicum-gophkeeper/api/gophkeeper"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/auth"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/config"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) Update(ctx context.Context, secretID string, name string, secretMetadata []byte, binaryFilePath string) error {
	data := &pb.SecretData{}
	data.SetText(string(secretMetadata))
	if binaryFilePath != "" {
		fileData, err := os.ReadFile(binaryFilePath)
		if err != nil {
			return fmt.Errorf("unable to read file %s: %w", binaryFilePath, err)
		}
		data.SetData(fileData)
	}

	request := pb.SecretUpdateRequest_builder{
		Id:   proto.String(secretID),
		Name: proto.String(name),
		Data: data,
	}.Build()

	token, err := auth.GetSessionToken()
	if err != nil {
		return fmt.Errorf("could not get session token: %w", err)
	}

	ctx = metadata.AppendToOutgoingContext(ctx, config.MetadataAuthKeyName, token)

	_, err = h.Client.Update(ctx, request)
	if err != nil {
		return fmt.Errorf("could not handle update: %w", err)
	}

	return nil
}
