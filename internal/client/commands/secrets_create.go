package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/handlers/grpc/secret"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/repository/secret_data"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	"github.com/spf13/cobra"
)

func GetCreateCmd(secretsHandler *secret.Handler) *cobra.Command {
	var name, secretType, binaryFilePath string

	secretTypes := []string{models.SecretTypeBinary, models.SecretTypeText, models.SecretTypeBankCard, models.SecretTypeLoginPassword}

	cmd := &cobra.Command{
		GroupID: CommandGroupSecrets,
		Use: `create 
	--name=SECRET_NAME 
	--type=SECRET_TYPE`,
		Short: "create a new secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			// name
			for name == "" {
				fmt.Print("name is required, enter name: ")
				_, err := fmt.Scanln(&name)
				if err != nil {
					fmt.Printf("error: could not read name: %s\n", err.Error())
				}
			}

			// type
			for secretType != models.SecretTypeBinary && secretType != models.SecretTypeText && secretType != models.SecretTypeBankCard && secretType != models.SecretTypeLoginPassword {
				fmt.Print("type is required, enter one of types: " + strings.Join(secretTypes, ", ") + ": ")
				_, err := fmt.Scanln(&secretType)
				if err != nil {
					fmt.Printf("error: could not read type: %s\n", err.Error())
				}
			}

			var secretDataJSON []byte
			var err error

			switch secretType {
			case models.SecretTypeText:
				secretData := secret_data.SecretTypeTextData{}
				secretDataJSON, _ = json.Marshal(secretData)
				repo := secret_data.NewStandardInputOutputRepo[*secret_data.SecretTypeTextData]()
				secretDataJSON, binaryFilePath, err = repo.GetJSON(string(secretDataJSON), true)
				if err != nil {
					return fmt.Errorf("could not get JSON: %w", err)
				}
			case models.SecretTypeBinary:
				secretData := secret_data.SecretTypeBinaryData{}
				secretDataJSON, _ = json.Marshal(secretData)
				repo := secret_data.NewStandardInputOutputRepo[*secret_data.SecretTypeBinaryData]()
				secretDataJSON, binaryFilePath, err = repo.GetJSON(string(secretDataJSON), true)
				if err != nil {
					return fmt.Errorf("could not get JSON: %w", err)
				}
			case models.SecretTypeLoginPassword:
				secretData := secret_data.SecretTypeLoginPasswordData{}
				secretDataJSON, _ = json.Marshal(secretData)
				repo := secret_data.NewStandardInputOutputRepo[*secret_data.SecretTypeLoginPasswordData]()
				secretDataJSON, binaryFilePath, err = repo.GetJSON(string(secretDataJSON), true)
				if err != nil {
					return fmt.Errorf("could not get JSON: %w", err)
				}
			case models.SecretTypeBankCard:
				secretData := secret_data.SecretTypeBankCardData{}
				secretDataJSON, _ = json.Marshal(secretData)
				repo := secret_data.NewStandardInputOutputRepo[*secret_data.SecretTypeBankCardData]()
				secretDataJSON, binaryFilePath, err = repo.GetJSON(string(secretDataJSON), true)
				if err != nil {
					return fmt.Errorf("could not get JSON: %w", err)
				}
			}

			id, err := secretsHandler.Store(cmd.Context(), name, secretType, secretDataJSON, binaryFilePath)
			if err != nil {
				return fmt.Errorf("could not store: %w", err)
			}

			fmt.Println("secret was saved with ID=" + id)

			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "secret name")
	cmd.Flags().StringVar(&secretType, "type", "", "type of secret: "+strings.Join(secretTypes, ", "))
	cmd.Flags().StringVar(&binaryFilePath, "path", "", "path to binary file")

	return cmd
}
