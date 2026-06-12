package config

import (
	"crypto/rsa"
	"saetechnology-be/internal/delivery/http/security/jwt"
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort             string
	AppEnv              string
	DatabaseConfig      *DatabaseConfig
	RedisAddr           string
	RabbitMQURL         string
	Secret              string
	BrevoApiKey         string
	RegisterSenderName  string
	RegisterSenderEmail string
	MailpitSMTPAddr     string
	QueueRegisterEmail  string
	RegisterLinkUrl     string
	PanelBaseURL        string
	CustomerBaseURL     string
	JwtPrivateKey       *rsa.PrivateKey
	JwtPublicKey        *rsa.PublicKey
	S3AwsRegion         string
	S3BucketName        string
	S3AccessKeyID       string
	S3SecretAccessKey   string
	S3Endpoint          string
	S3PublicBaseURL     string
	TracerEndPoint      string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file found (production mode assumed)")
	}
	jwtPrivateKey, err := jwt.LoadPrivateKey()
	if err != nil {
		log.Fatal(err)
	}
	jwtPublicKey, err := jwt.LoadPublicKey()
	if err != nil {
		log.Fatal(err)
	}
	return &Config{
		AppPort:             getEnv("APP_PORT", "3135"),
		AppEnv:              getEnv("APP_ENV", "development"),
		DatabaseConfig:      LoadDatabaseConfig(),
		RedisAddr:           getEnv("REDIS_ADDR", "localhost:6379"),
		RabbitMQURL:         getEnv("RABBITMQ_URL", ""),
		Secret:              getEnv("SECRET_KEY", ""),
		RegisterSenderName:  getEnv("EMAIL_REGISTER_NAME", ""),
		RegisterSenderEmail: getEnv("EMAIL_REGISTER_ADDRESS", ""),
		MailpitSMTPAddr:     getEnv("MAILPIT_SMTP_ADDR", "localhost:1025"),
		RegisterLinkUrl:     getEnv("EMAIL_REGISTER_LINK", ""),
		PanelBaseURL:        getEnv("PANEL_BASE_URL", "http://localhost:3901"),
		CustomerBaseURL:     getEnv("CUSTOMER_BASE_URL", "http://localhost:3902"),
		QueueRegisterEmail:  getEnv("RABBITMQ_QUEUE_EMAIL", ""),
		JwtPrivateKey:       jwtPrivateKey,
		JwtPublicKey:        jwtPublicKey,
		BrevoApiKey:         getEnv("BREVO_API_KEY", ""),
		S3AwsRegion:         getEnv("S3_AWS_REGION", ""),
		S3Endpoint:          getEnv("S3_ENDPOINT", ""),
		S3BucketName:        getEnv("S3_BUCKET_NAME", ""),
		S3PublicBaseURL:     getEnv("S3_PUBLIC_BASE_URL", ""),
		S3AccessKeyID:       getEnv("S3_ACCESS_KEY_ID", ""),
		S3SecretAccessKey:   getEnv("S3_SECRET_ACCESS_KEY", ""),
		TracerEndPoint:      getEnv("TRACER_END_POINT", ""),
	}
}
