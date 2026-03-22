package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/acya-skulskaya/yandex-practicum-gophkeeper/api/gophkeeper"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/config"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/logger"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/crypto"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/db"
	authGRPCHandler "github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/handlers/grpc/auth"
	secretGRPCHandler "github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/handlers/grpc/secret"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/interceptors"
	secretRepo "github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/repository/secret"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/repository/user"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/service/auth"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/service/secret"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	shutdownPeriod        = 20 * time.Second
	buildInfoDefaultValue = "N/A"
)

var (
	buildInfoVersion    = buildInfoDefaultValue
	buildInfoDate       = buildInfoDefaultValue
	buildInfoCommitHash = buildInfoDefaultValue
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.GetServerConfig()
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	// init logger
	if err = logger.Init(cfg.Logging); err != nil {
		return fmt.Errorf("could not init logging: %w", err)
	}
	logger.Log.Info("config loaded", zap.Any("cfg", cfg))

	// print build info
	cfg.BuildInfo.Version, cfg.BuildInfo.Date, cfg.BuildInfo.CommitHash = config.GetBuildInfo(buildInfoDefaultValue, buildInfoVersion, buildInfoDate, buildInfoCommitHash, cfg.BuildInfo)
	logger.Log.Info("build info", zap.String("version", cfg.BuildInfo.Version), zap.String("date", cfg.BuildInfo.Date), zap.String("commit_hash", cfg.BuildInfo.CommitHash))

	// init DB
	dbPool, err := db.New(ctx, cfg.DB)
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	// setup signal context
	rootCtx, stopRootCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stopRootCtx()

	// gRPC server
	listen, err := net.Listen("tcp", cfg.GRPCServer.Address)
	if err != nil {
		//nolint:gocritic // ignore
		log.Fatalf("failed listen: %v", err)
	}
	tlsCreds, err := config.LoadTLSCredentials(cfg.TLS)
	if err != nil {
		logger.Log.Debug("could not get TLS credentials", zap.Error(err))
	}
	var grpcServer *grpc.Server
	authInterceptor := interceptors.NewAuthInterceptor(cfg.Auth)
	chainUnaryInterceptor := grpc.ChainUnaryInterceptor(
		interceptors.LoggingUnaryInterceptor,
		authInterceptor.Unary(),
	)
	grpcServer = grpc.NewServer(
		grpc.Creds(tlsCreds),
		chainUnaryInterceptor,
	)

	// auth server
	userRepository := user.NewUserRepository(dbPool)
	authService := auth.New(cfg.Auth.SecretKey, userRepository)
	pb.RegisterAuthServiceServer(grpcServer, &authGRPCHandler.AuthServer{
		UnimplementedAuthServiceServer: pb.UnimplementedAuthServiceServer{},
		ServiceAuth:                    authService,
	})

	// secret server
	cry, err := crypto.NewCrypto(cfg.Secret.SecretKey, cfg.TLS.KeyFilePath)
	if err != nil {
		return fmt.Errorf("could not create crypto: %w", err)
	}
	secretRepository := secretRepo.NewSecretRepository(dbPool, cfg.Secret.FilesDir)
	secretService := secret.New(cfg.Secret.SecretKey, secretRepository, *cry)
	pb.RegisterSecretServiceServer(grpcServer, &secretGRPCHandler.SecretServer{
		UnimplementedSecretServiceServer: pb.UnimplementedSecretServiceServer{},
		ServiceSecret:                    secretService,
	})

	go func() {
		logger.Log.Info("starting gRPC server",
			zap.String("GRPCServerAddress", cfg.GRPCServer.Address),
		)
		if err := grpcServer.Serve(listen); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatalf("failed to listen and serve grpc: %v", err)
		}
	}()

	// Wait for signal
	<-rootCtx.Done()
	stopRootCtx()
	logger.Log.Info("received shutdown signal, shutting down")

	shutdownCtx, cancelShutdownCtx := context.WithTimeout(context.Background(), shutdownPeriod)
	defer cancelShutdownCtx()

	grpcServerStop := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(grpcServerStop)
	}()

	select {
	case <-grpcServerStop:
		logger.Log.Info("gRPC server gracefully stopped")
	case <-shutdownCtx.Done():
		logger.Log.Debug("gRPC graceful shutdown timed out, forcing shutdown...")
		grpcServer.Stop()
	}

	return nil
}
