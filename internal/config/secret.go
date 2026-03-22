package config

type Secret struct {
	//nolint:govet // ignore
	SecretKey string `env:"GOPHKEEPER_SECRET_SECRET_KEY" env-required`
	FilesDir  string `env:"GOPHKEEPER_SECRET_FILES_DIR" env-default:"./data/"`
}
