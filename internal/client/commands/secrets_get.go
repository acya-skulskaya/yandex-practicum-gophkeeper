package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/handlers/grpc/secret"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/repository/secret_data"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	"github.com/spf13/cobra"
)

func GetGetCmd(secretsHandler *secret.Handler) *cobra.Command {
	var secretID, secretVersionID string

	cmd := &cobra.Command{
		GroupID: CommandGroupSecrets,
		Use:     `get --id=SECRET_ID --version-id=SECRET_VERSION_ID`,
		Short:   "show a secret by id",
		Long:    `Will return secret's data (if it's not of type binary) for a specified version or a list of versions`,
		RunE: func(cmd *cobra.Command, args []string) error {
			for secretID == "" {
				fmt.Print("secret id is required, enter id: ")
				_, err := fmt.Scanln(&secretID)
				if err != nil {
					fmt.Printf("error: could not read name: %s\n", err.Error())
				}
			}

			if secretVersionID == "" {
				fmt.Print("secret version id is empty, you can specify version id, write 0 to get latest version data or leave empty to get a list of versions: ")
				//nolint:errcheck,gosec // can be left empty
				fmt.Scanln(&secretVersionID)
			}

			secret, err := secretsHandler.Show(cmd.Context(), secretID, secretVersionID)
			if err != nil {
				return fmt.Errorf("could not get: %w", err)
			}

			emdash := "—"
			dash := "–"
			format := "%-5s | %-15s | %-20s | %-25s | %-25s | %s\n"

			fmt.Println(strings.Repeat(emdash, 115))
			fmt.Printf(format, "ID", "TYPE", "NAME", "CREATED", "UPDATED", "NUM.VER.")
			fmt.Println(strings.Repeat(dash, 115))
			fmt.Printf(format,
				strconv.Itoa(int(secret.ID)),
				secret.Type,
				secret.Name,
				secret.CreatedAt.Format(time.DateTime),
				secret.UpdatedAt.Format(time.DateTime),
				strconv.Itoa(len(secret.Versions)))
			fmt.Println(strings.Repeat(emdash, 115))

			if secretVersionID == "" {
				fmt.Println("VERSIONS:")
				format = "%-5s | %s\n"
				fmt.Println(strings.Repeat(emdash, 30))
				fmt.Printf(format, "ID", "CREATED")
				fmt.Println(strings.Repeat(dash, 30))
				for _, s := range secret.Versions {
					fmt.Printf(format,
						strconv.Itoa(int(s.ID)),
						s.CreatedAt.Format(time.DateTime),
					)
				}
				fmt.Println(strings.Repeat(emdash, 30))
			} else {
				fmt.Printf("VERSION ID: %d\n", secret.Versions[0].ID)
				fmt.Printf("VERSION CREATE DATE: %s\n", secret.Versions[0].CreatedAt.Format(time.DateTime))

				switch secret.Type {
				case models.SecretTypeBinary:
					repo := secret_data.NewStandardInputOutputRepo[*secret_data.SecretTypeBinaryData]()
					repo.Print(secret.Versions[0].Data)
				case models.SecretTypeText:
					repo := secret_data.NewStandardInputOutputRepo[*secret_data.SecretTypeTextData]()
					repo.Print(secret.Versions[0].Data)
				case models.SecretTypeLoginPassword:
					repo := secret_data.NewStandardInputOutputRepo[*secret_data.SecretTypeLoginPasswordData]()
					repo.Print(secret.Versions[0].Data)
				case models.SecretTypeBankCard:
					repo := secret_data.NewStandardInputOutputRepo[*secret_data.SecretTypeBankCardData]()
					repo.Print(secret.Versions[0].Data)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&secretID, "id", "", "secret id")
	cmd.Flags().StringVar(&secretVersionID, "version-id", "", "secret version id")

	return cmd
}
