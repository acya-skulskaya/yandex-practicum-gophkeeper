package main

import (
	"context"
	"fmt"
	"os"

	pb "github.com/acya-skulskaya/yandex-practicum-gophkeeper/api/gophkeeper"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/commands"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/handlers/grpc/auth"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/handlers/grpc/secret"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/config"
	"google.golang.org/grpc"
)

const buildInfoDefaultValue = "N/A"

var (
	buildInfoVersion    = buildInfoDefaultValue
	buildInfoDate       = buildInfoDefaultValue
	buildInfoCommitHash = buildInfoDefaultValue
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Printf("\nError: %v\n", err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.GetConfigClient()
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	cfg.BuildInfo = config.BuildInfo{
		Version:    buildInfoVersion,
		Date:       buildInfoDate,
		CommitHash: buildInfoCommitHash,
	}

	tlsCreds, err := config.LoadTLSCredentials(cfg.TLS)
	if err != nil {
		return fmt.Errorf("could not load TLS credentials: %w", err)
	}

	// Устанавливаем соединение с сервером
	conn, err := grpc.NewClient(cfg.GRPCClient.Address, grpc.WithTransportCredentials(tlsCreds))
	if err != nil {
		return fmt.Errorf("could not create grpc client: %w", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("\nERROR: Could not close gRPC connection: %v\n", err.Error())
		}
	}(conn)

	secretsClient := pb.NewSecretServiceClient(conn)
	secretsHandler := secret.New(secretsClient)

	authClient := pb.NewAuthServiceClient(conn)
	authHandler := auth.New(authClient)

	err = commands.Execute(ctx, &cfg, authHandler, secretsHandler)
	if err != nil {
		return fmt.Errorf("could not execute command: %w", err)
	}
	return nil
}
