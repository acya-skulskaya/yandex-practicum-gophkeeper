package commands

import (
	"fmt"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/handlers/grpc/secret"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/repository/secret_data"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	"github.com/spf13/cobra"
)

func GetUpdateCmd(secretsHandler *secret.Handler) *cobra.Command {
	var secretID, name, binaryFilePath string

	cmd := &cobra.Command{
		GroupID: CommandGroupSecrets,
		Use: `update 
	--id=SECRET_ID
	--name=SECRET_NAME`,
		Short: "update an existing secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			for secretID == "" {
				fmt.Print("secret id is required, enter id: ")
				_, err := fmt.Scanln(&secretID)
				if err != nil {
					fmt.Printf("error: could not read name: %s\n", err.Error())
				}
			}

			secret, err := secretsHandler.Show(cmd.Context(), secretID, "0")
			if err != nil {
				return fmt.Errorf("could not get secret: %w", err)
			}

			if name == "" {
				fmt.Print("enter a new secret name if you want to update it or leave empty: ")
				//nolint:errcheck,gosec // can be left empty
				fmt.Scanln(&name)
			}
			secretType := secret.Type

			var secretDataJSON []byte

			switch secretType {
			case models.SecretTypeText:
				repo := secret_data.NewStandardInputOutputRepo[*secret_data.SecretTypeTextData]()
				// text should be always filled for this secret type
				secretDataJSON, binaryFilePath, err = repo.GetJSON(secret.Versions[0].Data, true)
				if err != nil {
					return fmt.Errorf("could not get JSON: %w", err)
				}
			case models.SecretTypeBinary:
				repo := secret_data.NewStandardInputOutputRepo[*secret_data.SecretTypeBinaryData]()
				// new file path should be always filled for this secret type
				secretDataJSON, binaryFilePath, err = repo.GetJSON(secret.Versions[0].Data, true)
				if err != nil {
					return fmt.Errorf("could not get JSON: %w", err)
				}
			case models.SecretTypeLoginPassword:
				repo := secret_data.NewStandardInputOutputRepo[*secret_data.SecretTypeLoginPasswordData]()
				secretDataJSON, binaryFilePath, err = repo.GetJSON(secret.Versions[0].Data, false)
				if err != nil {
					return fmt.Errorf("could not get JSON: %w", err)
				}
			case models.SecretTypeBankCard:
				repo := secret_data.NewStandardInputOutputRepo[*secret_data.SecretTypeBankCardData]()
				secretDataJSON, binaryFilePath, err = repo.GetJSON(secret.Versions[0].Data, false)
				if err != nil {
					return fmt.Errorf("could not get JSON: %w", err)
				}
			}

			if name == "" {
				name = secret.Name
			}

			err = secretsHandler.Update(cmd.Context(), secretID, name, secretDataJSON, binaryFilePath)
			if err != nil {
				return fmt.Errorf("could not update: %w", err)
			}

			fmt.Printf("secret %s was updated\n", secretID)

			return nil
		},
	}

	cmd.Flags().StringVar(&secretID, "id", "", "secret id")
	cmd.Flags().StringVar(&name, "name", "", "secret name")

	return cmd
}
