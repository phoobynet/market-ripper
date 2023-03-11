package query

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/phoobynet/market-ripper/config"
	"log"
)

var connection *pgx.Conn

func Connect(configuration *config.Config) {
	pgConn, err := pgx.Connect(context.TODO(), fmt.Sprintf("postgresql://admin:quest@%s:%s/qdb", configuration.DBHost, configuration.DBPGPort))

	if err != nil {
		log.Fatal(err)
	}

	connection = pgConn
}

func Disconnect() {
	_ = connection.Close(context.TODO())
}
