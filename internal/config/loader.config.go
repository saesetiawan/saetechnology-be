package config

func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		Name:     getEnv("DB_NAME", ""),
		Username: getEnv("DB_USERNAME", ""),
		Password: getEnv("DB_PASSWORD", ""),
		Secure:   getEnv("DB_SECURE", "false") == "true",
		CAFile:   getEnv("DB_CA", ""),
	}
}
