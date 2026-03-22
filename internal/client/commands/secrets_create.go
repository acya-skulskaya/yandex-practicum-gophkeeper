package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/handlers/grpc/secret"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/config"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func GetCreateCmd(secretsHandler *secret.Handler) *cobra.Command {
	var name, secretType, binaryFilePath, text, lpLogin, lpPassword, bcNumber, bcHolderName, bcCVV, bcExpiry, metadata string

	secretTypes := []string{models.SecretTypeBinary, models.SecretTypeText, models.SecretTypeBankCard, models.SecretTypeLoginPassword}

	cmd := &cobra.Command{
		GroupID: CommandGroupSecrets,
		Use: `create 
	--name=SECRET_NAME 
	--type=SECRET_TYPE 
	--path=PATH_TO_BINARY_FILE 
	--metadata=METADATA 
	--login=LOGIN_TO_SAVE 
	--password=PASSWORD_TO_SAVE 
	--bc-number=BANK_CARD_NUMBER 
	--bc-cvv=BANK_CARD_CVV 
	--bc-expiry=BANK_CARD_EXPIRY_DATE 
	--bc-holder-name=BANK_CARD_HOLDER_NAME`,
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

			// metadata
			if metadata == "" {
				fmt.Print("enter metadata or leave empty: ")
				//nolint:errcheck,gosec // can be left empty
				fmt.Scanln(&metadata)
			}

			var secretDataJSON []byte

			switch secretType {
			case models.SecretTypeText:
				for text == "" {
					fmt.Print("secret text is required, enter text: ")
					_, err := fmt.Scanln(&text)
					if err != nil {
						fmt.Printf("error: could not read text: %s\n", err.Error())
					}
				}
				secretData := models.SecretTypeTextData{Text: text, Metadata: metadata}
				secretDataJSON, _ = json.Marshal(secretData)
			case models.SecretTypeBinary:
				var fileInfo os.FileInfo
				for binaryFilePath == "" {
					fmt.Print("file path to the binary file is required, enter path: ")
					_, err := fmt.Scanln(&binaryFilePath)
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

				secretData := models.SecretTypeBinaryData{
					OriginalFilePath: binaryFilePath,
					FileName:         fileInfo.Name(),
					FileSize:         fileInfo.Size(),
					Metadata:         metadata,
				}
				secretDataJSON, _ = json.Marshal(secretData)
			case models.SecretTypeLoginPassword:
				for lpLogin == "" {
					fmt.Print("login is required, enter login: ")
					_, err := fmt.Scanln(&lpLogin)
					if err != nil {
						fmt.Printf("error: could not read login: %s\n", err.Error())
					}
				}

				for lpPassword == "" {
					fmt.Print("password is required, enter password: ")
					bytePassword, err := term.ReadPassword(syscall.Stdin)
					if err != nil {
						fmt.Printf("error: could not read password: %s\n", err.Error())
					}
					lpPassword = string(bytePassword)
					fmt.Println()
				}
				secretData := models.SecretTypeLoginPasswordData{
					Password: lpPassword,
					Login:    lpLogin,
					Metadata: metadata,
				}
				secretDataJSON, _ = json.Marshal(secretData)
			case models.SecretTypeBankCard:
				for bcNumber == "" {
					fmt.Print("card number is required, enter card number: ")
					_, err := fmt.Scanln(&bcNumber)
					if err != nil {
						fmt.Printf("error: could not read card number: %s\n", err.Error())
					}
				}

				for bcHolderName == "" {
					fmt.Print("card holder name is required, enter holder name: ")
					_, err := fmt.Scanln(&bcHolderName)
					if err != nil {
						fmt.Printf("error: could not read holder name: %s\n", err.Error())
					}
				}

				for bcCVV == "" {
					fmt.Print("CVV is required, enter CVV: ")
					_, err := fmt.Scanln(&bcCVV)
					if err != nil {
						fmt.Printf("error: could not read CVV: %s\n", err.Error())
					}
				}

				for bcExpiry == "" {
					fmt.Print("card expiry date is required, enter expiry date: ")
					_, err := fmt.Scanln(&bcExpiry)
					if err != nil {
						fmt.Printf("error: could not read expiry date: %s\n", err.Error())
					}
				}
				secretData := models.SecretTypeBankCardData{
					Number:     bcNumber,
					Expiry:     bcExpiry,
					HolderName: bcHolderName,
					CVV:        bcCVV,
					Metadata:   metadata,
				}
				secretDataJSON, _ = json.Marshal(secretData)
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
	cmd.Flags().StringVar(&text, "text", "", "text to save")
	cmd.Flags().StringVar(&lpLogin, "login", "", "login to save")
	cmd.Flags().StringVar(&lpPassword, "password", "", "password to save")
	cmd.Flags().StringVar(&bcNumber, "number", "", "bank card number")
	cmd.Flags().StringVar(&bcCVV, "cvv", "", "bank card CVV")
	cmd.Flags().StringVar(&bcExpiry, "expiry", "", "bank card expiry")
	cmd.Flags().StringVar(&bcHolderName, "holder", "", "bank card holder name")
	cmd.Flags().StringVar(&metadata, "metadata", "", "metadata to save")

	return cmd
}
