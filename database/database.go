package database

import (
	"FlexMadeTest/configuration"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(cfg *configuration.Configuration) (*gorm.DB, error) {
	connection, err := gorm.Open(postgres.Open(cfg.DBConfig.GetPostgresDsn()), &gorm.Config{})
	if err != nil {
		return connection, err
	}

	return connection, err
}
