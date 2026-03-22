package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/handlers/grpc/secret"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/config"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func GetUpdateCmd(secretsHandler *secret.Handler) *cobra.Command {
	var secretID, name, binaryFilePath, text, lpLogin, lpPassword, bcNumber, bcHolderName, bcCVV, bcExpiry, metadata string

	cmd := &cobra.Command{
		GroupID: CommandGroupSecrets,
		Use: `update 
	--id=SECRET_ID
	--name=SECRET_NAME
	--path=PATH_TO_BINARY_FILE 
	--metadata=METADATA 
	--login=LOGIN_TO_SAVE 
	--password=PASSWORD_TO_SAVE 
	--bc-number=BANK_CARD_NUMBER 
	--bc-cvv=BANK_CARD_CVV 
	--bc-expiry=BANK_CARD_EXPIRY_DATE 
	--bc-holder-name=BANK_CARD_HOLDER_NAME`,
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

			if metadata == "" {
				fmt.Print("enter new metadata or leave empty: ")
				//nolint:errcheck,gosec // can be left empty
				fmt.Scanln(&metadata)
			}

			var secretDataJSON []byte

			switch secretType {
			case models.SecretTypeText:
				if text == "" {
					fmt.Print("enter new secret text if you want to update it or leave empty: ")
					//nolint:errcheck,gosec // can be left empty
					fmt.Scanln(&text)
				}
				secretData := models.SecretTypeTextData{Text: secret.Versions[0].Data}
				if text != "" {
					secretData.Text = text
				}
				if metadata != "" {
					secretData.Metadata = metadata
				}

				secretDataJSON, err = json.Marshal(secretData)
				if err != nil {
					return fmt.Errorf("could not marshal secretData: %w", err)
				}
			case models.SecretTypeBinary:
				var fileInfo os.FileInfo
				for binaryFilePath == "" {
					fmt.Print("enter new path to the binary file: ")
					_, err = fmt.Scanln(&binaryFilePath)
					if err != nil {
						fmt.Printf("error: could not read path: %s\n", err.Error())
					}
					fileInfo, err = os.Stat(binaryFilePath)
					if err != nil {
						fmt.Printf("error: could not read path: %s\n", err.Error())
						binaryFilePath = ""
					} else if fileInfo.IsDir() {
						fmt.Printf("error: binary file path can't be a directory\n")
						binaryFilePath = ""
					}
				}

				if fileInfo.Size() > config.MaxBinaryFileSize {
					return fmt.Errorf("error: binary file size is too large (max size is %dMB)", config.MaxBinaryFileSize/1024/1024)
				}

				var secretData models.SecretTypeBinaryData
				err = json.Unmarshal([]byte(secret.Versions[0].Data), &secretData)
				if err != nil {
					return fmt.Errorf("could not unmarshal secret data: %w", err)
				}
				secretData.OriginalFilePath = binaryFilePath
				secretData.FileName = fileInfo.Name()
				secretData.FileSize = fileInfo.Size()
				if metadata != "" {
					secretData.Metadata = metadata
				}

				secretDataJSON, err = json.Marshal(secretData)
				if err != nil {
					return fmt.Errorf("could not marshal secret data: %w", err)
				}
			case models.SecretTypeLoginPassword:
				for lpLogin == "" {
					fmt.Print("enter new login if you want to update it or leave empty: ")
					//nolint:errcheck,gosec // can be left empty
					fmt.Scanln(&lpLogin)
				}

				for lpPassword == "" {
					fmt.Print("enter new password if you want to update it or leave empty: ")
					bytePassword, _ := term.ReadPassword(syscall.Stdin)
					lpPassword = string(bytePassword)
					fmt.Println()
				}
				var secretData models.SecretTypeLoginPasswordData
				err = json.Unmarshal([]byte(secret.Versions[0].Data), &secretData)
				if err != nil {
					return fmt.Errorf("could not unmarshal secretData: %w", err)
				}
				if lpLogin != "" {
					secretData.Login = lpLogin
				}
				if lpPassword != "" {
					secretData.Password = lpPassword
				}
				if metadata != "" {
					secretData.Metadata = metadata
				}

				secretDataJSON, err = json.Marshal(secretData)
				if err != nil {
					return fmt.Errorf("could not marshal secretData: %w", err)
				}
			case models.SecretTypeBankCard:
				if bcNumber == "" {
					fmt.Print("enter new card number if you want to update it or leave empty: ")
					//nolint:errcheck,gosec // can be left empty
					fmt.Scanln(&bcNumber)
				}

				if bcHolderName == "" {
					fmt.Print("enter new card holder name if you want to update it or leave empty: ")
					//nolint:errcheck,gosec // can be left empty
					fmt.Scanln(&bcHolderName)
				}

				if bcCVV == "" {
					fmt.Print("enter new card CVV if you want to update it or leave empty: ")
					//nolint:errcheck,gosec // can be left empty
					fmt.Scanln(&bcCVV)
				}

				if bcExpiry == "" {
					fmt.Print("enter new card expiry date if you want to update it or leave empty: ")
					//nolint:errcheck,gosec // can be left empty
					fmt.Scanln(&bcExpiry)
				}
				var secretData models.SecretTypeBankCardData
				err = json.Unmarshal([]byte(secret.Versions[0].Data), &secretData)
				if err != nil {
					return fmt.Errorf("could not unmarshal secret data: %w", err)
				}

				if bcNumber != "" {
					secretData.Number = bcNumber
				}
				if bcHolderName != "" {
					secretData.HolderName = bcHolderName
				}
				if bcCVV != "" {
					secretData.CVV = bcCVV
				}
				if bcExpiry != "" {
					secretData.Expiry = bcExpiry
				}
				if metadata != "" {
					secretData.Metadata = metadata
				}

				secretDataJSON, err = json.Marshal(secretData)
				if err != nil {
					return fmt.Errorf("could not marshal secret data: %w", err)
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
	cmd.Flags().StringVar(&binaryFilePath, "path", "", "path to binary file")
	cmd.Flags().StringVar(&text, "text", "", "text to save")
	cmd.Flags().StringVar(&lpLogin, "login", "", "login to save")
	cmd.Flags().StringVar(&lpPassword, "password", "", "password to save")
	cmd.Flags().StringVar(&bcNumber, "number", "", "bank card number")
	cmd.Flags().StringVar(&bcCVV, "cvv", "", "bank card CVV")
	cmd.Flags().StringVar(&bcExpiry, "expiry", "", "bank card expiry")
	cmd.Flags().StringVar(&bcHolderName, "holder", "", "bank card holder name")
	cmd.Flags().StringVar(&metadata, "metadata", "", "metadata")

	return cmd
}
