package models

import (
	"time"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/api/gophkeeper"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/logger"
	"go.uber.org/zap"
)

const (
	SecretTable        = "secrets"
	SecretVersionTable = "secret_versions"

	SecretTypeLoginPassword = "login_password"
	SecretTypeText          = "text"
	SecretTypeBinary        = "binary"
	SecretTypeBankCard      = "bank_card"
)

type Secret struct {
	UpdatedAt *time.Time      `db:"updated_at"`
	CreatedAt *time.Time      `db:"created_at"`
	Name      string          `db:"name"`
	Type      string          `db:"type"`
	Versions  []SecretVersion `db:"-"`
	UserID    uint            `db:"user_id"`
	ID        uint            `db:"id"`
}

type SecretVersion struct {
	CreatedAt *time.Time `db:"created_at"`
	Data      string     `db:"data"`
	SecretID  uint       `db:"secret_id"`
	ID        uint       `db:"id"`
}

type SecretTypeLoginPasswordData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Metadata string `json:"metadata"`
}

type SecretTypeTextData struct {
	Text     string `json:"text"`
	Metadata string `json:"metadata"`
}

type SecretTypeBankCardData struct {
	Number     string `json:"number"`
	Expiry     string `json:"expiry"`
	HolderName string `json:"holder_name"`
	CVV        string `json:"cvv"`
	Metadata   string `json:"metadata"`
}

type SecretTypeBinaryData struct {
	OriginalFilePath string `json:"original_file_path"`
	FileName         string `json:"file_name"`
	Metadata         string `json:"metadata"`
	FileSize         int64  `json:"file_size"`
}

func MatchPBEnumWIthSecretType(enumString string) string {
	switch enumString {
	case gophkeeper.SecretType_SECRET_TYPE_LOGIN_PASSWORD.String():
		return SecretTypeLoginPassword
	case gophkeeper.SecretType_SECRET_TYPE_TEXT.String():
		return SecretTypeText
	case gophkeeper.SecretType_SECRET_TYPE_BANK_CARD.String():
		return SecretTypeBankCard
	case gophkeeper.SecretType_SECRET_TYPE_BINARY.String():
		return SecretTypeBinary
	}
	logger.Log.Warn("could not match enum type", zap.String("enumString", enumString))
	return ""
}

func MatchSecretTypeWithPBEnum(secretType string) gophkeeper.SecretType {
	switch secretType {
	case SecretTypeLoginPassword:
		return gophkeeper.SecretType_SECRET_TYPE_LOGIN_PASSWORD
	case SecretTypeText:
		return gophkeeper.SecretType_SECRET_TYPE_TEXT
	case SecretTypeBankCard:
		return gophkeeper.SecretType_SECRET_TYPE_BANK_CARD
	case SecretTypeBinary:
		return gophkeeper.SecretType_SECRET_TYPE_BINARY
	}
	return gophkeeper.SecretType_SECRET_TYPE_UNSPECIFIED
}
