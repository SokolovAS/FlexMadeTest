package configuration

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Configuration struct {
	ApplicationPort string `env:"APPLICATION_PORT"`
	DBConfig        DBConfig
}

type DBConfig struct {
	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
}

func New() (*Configuration, error) {
	cfg := Configuration{}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c DBConfig) GetPostgresDsn() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
	)
}
