package database

import (
	"github.com/phoobynet/market-ripper/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect connects to the QuestDB database using the Postgres protocol
func Connect(configuration *config.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(configuration.DSN()), &gorm.Config{})
}
