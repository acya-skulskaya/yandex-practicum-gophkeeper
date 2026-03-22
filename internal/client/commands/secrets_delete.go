package commands

import (
	"fmt"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/handlers/grpc/secret"
	"github.com/spf13/cobra"
)

func GetDeleteCmd(secretsHandler *secret.Handler) *cobra.Command {
	var secretID, secretVersionID string

	cmd := &cobra.Command{
		GroupID: CommandGroupSecrets,
		Use:     `delete --id=SECRET_ID --version-id=SECRET_VERSION_ID`,
		Short:   "delete a secret by id",
		RunE: func(cmd *cobra.Command, args []string) error {
			for secretID == "" {
				fmt.Print("secret id is required, enter id: ")
				_, err := fmt.Scanln(&secretID)
				if err != nil {
					fmt.Printf("error: could not read name: %s\n", err.Error())
				}
			}

			if secretVersionID == "" {
				fmt.Print("version id is empty, you can specify it or leave empty: ")
				//nolint:errcheck,gosec // can be left empty
				fmt.Scanln(&secretID)
			}

			err := secretsHandler.Delete(cmd.Context(), secretID, secretVersionID)
			if err != nil {
				return fmt.Errorf("could not delete: %w", err)
			}

			if secretVersionID != "" {
				fmt.Printf("secret's %s version %s was deleted\n", secretID, secretVersionID)
			} else {
				fmt.Printf("secret %s was deleted\n", secretID)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&secretID, "id", "", "secret id")
	cmd.Flags().StringVar(&secretVersionID, "version-id", "", "secret version id")

	return cmd
}
