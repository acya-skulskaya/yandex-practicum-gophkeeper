package commands

import (
	"fmt"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/auth"
	"github.com/spf13/cobra"
)

func GetLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		GroupID: CommandGroupAuth,
		Use:     "logout",
		Short:   "logout current user",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := auth.DeleteSession(); err != nil {
				return fmt.Errorf("could not logout session: %w", err)
			}
			fmt.Printf("successfully logged out\n")
			return nil
		},
	}

	return cmd
}
