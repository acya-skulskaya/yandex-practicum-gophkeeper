package commands

import (
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/config"
	"github.com/spf13/cobra"
)

func GetRootCmd(cfg *config.Config) *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	cmd := &cobra.Command{
		Short: "Gophkeeper CLI - a client for secrets keeper Gophkeeper",
		Long:  config.GetBuildInfoStr(cfg.BuildInfo),
	}

	return cmd
}
