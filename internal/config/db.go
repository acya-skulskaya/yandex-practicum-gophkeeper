package config

type DB struct {
	//nolint:govet // ignore
	DSN            string `env:"GOPHKEEPER_DB_DSN" env-default:"postgres://user:pass@localhost:5432/gophkeeper?sslmode=disable" env-required`
	MigrationsPath string `env:"GOPHKEEPER_DB_MIGRATIONS_PATH" env-default:"file://./migrations"`
}
