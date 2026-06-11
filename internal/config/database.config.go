package config

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
	Secure   bool
	CAFile   string
}
