package broker

import (
	rabbitmq "github.com/wagslane/go-rabbitmq"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/config"
)

func NewRabbitConnection(config *config.Config) *rabbitmq.Conn {
	conn, err := rabbitmq.NewConn(config.RabbitMQURL)
	if err != nil {
		panic(err)
	}
	return conn
}
