package database

import (
	"github.com/phoobynet/market-ripper/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(configuration *config.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(configuration.DSN()), &gorm.Config{})
}
