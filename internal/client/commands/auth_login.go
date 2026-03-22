package commands

import (
	"fmt"
	"syscall"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/auth"
	grpcHandler "github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/handlers/grpc/auth"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func GetLoginCmd(authHandler *grpcHandler.Handler) *cobra.Command {
	var login, password string

	cmd := &cobra.Command{
		GroupID: CommandGroupAuth,
		Use:     "login --login=YOUR_LOGIN --password=YOUR_PASSWORD",
		Short:   "log in to your account",
		RunE: func(cmd *cobra.Command, args []string) error {
			for login == "" {
				fmt.Print("login is required, enter login: ")
				_, err := fmt.Scanln(&login)
				if err != nil {
					fmt.Printf("error: could not read login: %s\n", err.Error())
				}
			}

			for password == "" {
				fmt.Print("password is required, enter password: ")
				bytePassword, err := term.ReadPassword(syscall.Stdin)
				if err != nil {
					fmt.Printf("error: could not read password: %s\n", err.Error())
				}
				password = string(bytePassword)
				fmt.Println()
			}

			fmt.Printf("trying to log in as %s\n", login)

			token, err := authHandler.Login(cmd.Context(), login, password)
			if err != nil {
				return fmt.Errorf("could not log in: %w", err)
			}

			if err = auth.SaveSession(token); err != nil {
				return fmt.Errorf("could not save session: %w", err)
			}

			fmt.Printf("successfully logged in as %s\n", login)

			return nil
		},
	}

	cmd.Flags().StringVar(&login, "login", "", "login")
	cmd.Flags().StringVar(&password, "password", "", "password")

	return cmd
}
