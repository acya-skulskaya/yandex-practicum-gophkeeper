package config

type BuildInfo struct {
	Version    string `env:"GOPHKEEPER_BUILD_INFO_VERSION" env-default:"N/A"`
	Date       string `env:"GOPHKEEPER_BUILD_INFO_DATE" env-default:"N/A"`
	CommitHash string `env:"GOPHKEEPER_BUILD_INFO_COMMIT_HASH" env-default:"N/A"`
}

func GetBuildInfo(defaultValue string, version string, date string, commitHash string, cfg BuildInfo) (cfgVersion string, cfgDate string, cfgCommitHash string) {
	if version != defaultValue {
		cfg.Version = version
	}
	if date != defaultValue {
		cfg.Date = date
	}
	if commitHash != defaultValue {
		cfg.CommitHash = commitHash
	}

	return cfg.Version, cfg.Date, cfg.CommitHash
}

func GetBuildInfoStr(info BuildInfo) string {
	return "BUILD INFO: \n" +
		"VERSION: " + info.Version + "\n" +
		"DATE: " + info.Date + "\n" +
		"COMMIT HASH: " + info.CommitHash
}
