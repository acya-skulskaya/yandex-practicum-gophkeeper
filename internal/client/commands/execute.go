package commands

import (
	"context"
	"fmt"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/handlers/grpc/auth"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/handlers/grpc/secret"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/config"
	"github.com/spf13/cobra"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context, cfg *config.Config, authHandler *auth.Handler, secretsHandler *secret.Handler) error {
	rootCmd := GetRootCmd(cfg)

	rootCmd.AddGroup(&cobra.Group{ID: CommandGroupAuth, Title: "user management commands"})
	rootCmd.AddGroup(&cobra.Group{ID: CommandGroupSecrets, Title: "secrets management commands"})
	rootCmd.AddGroup(&cobra.Group{ID: CommandGroupApp, Title: "app info commands"})

	// auth
	rootCmd.AddCommand(GetRegisterCmd(authHandler))
	rootCmd.AddCommand(GetLoginCmd(authHandler))
	rootCmd.AddCommand(GetLogoutCmd())

	// secrets
	rootCmd.AddCommand(GetCreateCmd(secretsHandler))
	rootCmd.AddCommand(GetUpdateCmd(secretsHandler))
	rootCmd.AddCommand(GetGetCmd(secretsHandler))
	rootCmd.AddCommand(GetDownloadCmd(secretsHandler))
	rootCmd.AddCommand(GetListCmd(secretsHandler))
	rootCmd.AddCommand(GetDeleteCmd(secretsHandler))

	// app
	rootCmd.AddCommand(GetVersionCmd(cfg.BuildInfo))

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		return fmt.Errorf("could not execute context %w", err)
	}

	return nil
}
