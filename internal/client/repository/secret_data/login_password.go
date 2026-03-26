package secret_data

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type SecretTypeLoginPasswordData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Metadata string `json:"metadata"`
}

func (b *SecretTypeLoginPasswordData) PrintData() error {
	fmt.Printf("LOGIN: %s\n", b.Login)
	fmt.Printf("PASSWORD: %s\n", b.Password)
	fmt.Printf("METADATA: %s\n", b.Metadata)
	return nil
}

func (b *SecretTypeLoginPasswordData) Fill(requireNotEmpty bool) (string, error) {
	var metadata, lpLogin, lpPassword string

	if metadata == "" {
		fmt.Print("enter new metadata or leave empty: ")
		//nolint:errcheck,gosec // can be left empty
		fmt.Scanln(&metadata)
	}

	if requireNotEmpty {
		for lpLogin == "" {
			fmt.Print("login is required, enter login: ")
			_, err := fmt.Scanln(&lpLogin)
			if err != nil {
				fmt.Printf("error: could not read login: %s\n", err.Error())
			}
		}
		for lpPassword == "" {
			fmt.Print("password is required, enter password: ")
			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				fmt.Printf("error: could not read password: %s\n", err.Error())
			}
			lpPassword = string(bytePassword)
			fmt.Println()
		}
	} else {
		if lpLogin == "" {
			fmt.Print("enter new login if you want to update it or leave empty: ")
			//nolint:errcheck,gosec // can be left empty
			fmt.Scanln(&lpLogin)
		}
		if lpPassword == "" {
			fmt.Print("enter new password if you want to update it or leave empty: ")
			bytePassword, _ := term.ReadPassword(int(os.Stdin.Fd()))
			lpPassword = string(bytePassword)
			fmt.Println()
		}
	}

	if lpLogin != "" {
		b.Login = lpLogin
	}
	if lpPassword != "" {
		b.Password = lpPassword
	}
	if metadata != "" {
		b.Metadata = metadata
	}

	return "", nil
}
