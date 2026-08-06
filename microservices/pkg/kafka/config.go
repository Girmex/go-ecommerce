package kafka

import "github.com/caarlos0/env/v11"

type Config struct {
	Brokers []string `env:"KAFKA_BROKERS" envSeparator:","`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}