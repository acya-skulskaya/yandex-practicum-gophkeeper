package auth

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	KeyringServiceName = "gophkeeper"
	KeyringTokenKey    = "token"
)

func SaveSession(token string) error {
	err := keyring.Set(KeyringServiceName, KeyringTokenKey, token)
	if err != nil {
		return fmt.Errorf("could not set token: %w", err)
	}
	return nil
}

func DeleteSession() error {
	err := keyring.Delete(KeyringServiceName, KeyringTokenKey)
	if err != nil {
		return fmt.Errorf("could not unset token: %w", err)
	}
	return nil
}

func GetSessionToken() (string, error) {
	token, err := keyring.Get(KeyringServiceName, KeyringTokenKey)
	if err != nil {
		return "", fmt.Errorf("could not get token: %w", err)
	}
	return token, nil
}
