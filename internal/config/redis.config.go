package config

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func LoadRedisConfig(config *Config) *RedisConfig {
	return &RedisConfig{
		Addr:     config.RedisAddr,
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       1,
	}
}
