package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	dotEnvFileName    = ".env"
	MaxBinaryFileSize = 32 * 1024 * 1024 // 32 MB
)

type Config struct {
	Logging    Logging
	BuildInfo  BuildInfo
	DB         DB
	Secret     Secret
	Auth       Auth
	GRPCServer GRPCServer
	GRPCClient GRPCClient
	TLS        TLS
}

func GetConfigClient() (Config, error) {
	cfg := Config{}

	if _, err := os.Stat(dotEnvFileName); err == nil {
		err := godotenv.Load()
		if err != nil {
			return cfg, fmt.Errorf("error loading dot env file: %w", err)
		}
	}

	cfgGRPCClient := GRPCClient{}
	err := cleanenv.ReadEnv(&cfgGRPCClient)
	if err != nil {
		return cfg, fmt.Errorf("could not read gRPC Client config %w", err)
	}
	cfg.GRPCClient = cfgGRPCClient

	cfgTLS := TLS{}
	err = cleanenv.ReadEnv(&cfgTLS)
	if err != nil {
		return cfg, fmt.Errorf("could not read TLS config %w", err)
	}
	cfgTLS.AutoCert = true   // for development purposes
	cfgTLS.SkipVerify = true // for development purposes
	cfg.TLS = cfgTLS

	return cfg, nil
}

func GetServerConfig() (Config, error) {
	cfg := Config{}

	if _, err := os.Stat(dotEnvFileName); err == nil {
		err := godotenv.Load()
		if err != nil {
			return cfg, fmt.Errorf("error loading dot env file: %w", err)
		}
	}

	// logging
	cfgLogging := Logging{}
	err := cleanenv.ReadEnv(&cfgLogging)
	if err != nil {
		return cfg, fmt.Errorf("could not read Logging config %w", err)
	}
	cfg.Logging = cfgLogging

	// build info
	cfgBuildInfo := BuildInfo{}
	err = cleanenv.ReadEnv(&cfgBuildInfo)
	if err != nil {
		return cfg, fmt.Errorf("could not read BuildInfo config %w", err)
	}
	cfg.BuildInfo = cfgBuildInfo

	// db
	cfgDB := DB{}
	err = cleanenv.ReadEnv(&cfgDB)
	if err != nil {
		return cfg, fmt.Errorf("could not read DB config %w", err)
	}
	cfg.DB = cfgDB

	// auth
	var cfgAuth Auth
	err = cleanenv.ReadEnv(&cfgAuth)
	if err != nil {
		return cfg, fmt.Errorf("could not read Auth config %w", err)
	}
	if cfgAuth.SecretKey == "" {
		return cfg, fmt.Errorf("can't run gophkeeper server with empty auth secret key")
	}
	cfg.Auth = cfgAuth

	// secret
	var cfgSecret Secret
	err = cleanenv.ReadEnv(&cfgSecret)
	if err != nil {
		return cfg, fmt.Errorf("could not read Secret config %w", err)
	}
	if cfgSecret.SecretKey == "" {
		return cfg, fmt.Errorf("can't run gophkeeper server with empty secret key")
	}
	cfg.Secret = cfgSecret

	// gRPC server
	cfgGRPCServer := GRPCServer{}
	err = cleanenv.ReadEnv(&cfgGRPCServer)
	if err != nil {
		return cfg, fmt.Errorf("could not read gRPC Server config %w", err)
	}
	cfg.GRPCServer = cfgGRPCServer

	// tls
	cfgTLS := TLS{}
	err = cleanenv.ReadEnv(&cfgTLS)
	if err != nil {
		return cfg, fmt.Errorf("could not read TLS config %w", err)
	}
	if cfgTLS.AutoCert {
		return cfg, fmt.Errorf("can't use autocert for gophkeeper server. please specify file for certificate and key")
	}
	cfg.TLS = cfgTLS

	return cfg, nil
}
