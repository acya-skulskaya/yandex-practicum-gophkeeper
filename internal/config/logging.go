package config

type Logging struct {
	Level               string `env:"GOPHKEEPER_LOGGING_LEVEL" env-default:"info"`
	OutputFilePath      string `env:"GOPHKEEPER_LOGGING_OUTPUT_FILE_PATH" env-default:"./logs/output.log"`
	ErrorOutputFilePath string `env:"GOPHKEEPER_LOGGING_ERROR_OUTPUT_FILE_PATH" env-default:"./logs/error.log"`
	OutputToFile        bool   `env:"GOPHKEEPER_LOGGING_OUTPUT_TO_FILE" env-default:"false"`
	ErrorOutputToFile   bool   `env:"GOPHKEEPER_LOGGING_ERROR_OUTPUT_TO_FILE" env-default:"false"`
}
