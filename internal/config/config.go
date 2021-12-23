package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP
	DataBase
	Nats
}

type HTTP struct {
	Port string `env:"PORT_HTTP" env-default:"8081"`
}

type DataBase struct {
	URL string `env:"DB_URL" env-default:"user=postgres password=postgres dbname=orders host=localhost sslmode=disable"`
}
type Nats struct {
	ClusterId   string `env:"NATS_CLUSTER" env-default:"test-cluster"`
	ClientId    string `env:"NATS_CLIENT" env-default:"test-user"`
	URL         string `env:"NATS_URL" env-default:"nats://127.0.0.1:4223"`
	Subject     string `env:"NATS_SUBJECT" env-default:"order"`
	DurableName string `env:"NATS_DURABLE_NAME" env-default:"client1"`
	AckWait     string `env:"NATS_ACK_WAIT" env-default:"10"`
}

func New() (Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		return cfg, err
	}
	return cfg, err
}
