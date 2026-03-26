package config

const (
	MetadataAuthKeyName = "authorization"
)

type Auth struct {
	//nolint:govet // ignore
	SecretKey string `env:"GOPHKEEPER_AUTH_SECRET_KEY" env-required`
}
