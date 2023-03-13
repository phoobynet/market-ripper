package query

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/phoobynet/market-ripper/config"
	"log"
)

var connection *pgx.Conn

func Connect(configuration *config.Config) {
	pgConn, err := pgx.Connect(context.TODO(), configuration.GetPGAddress())

	if err != nil {
		log.Fatal(err)
	}

	err = pgConn.Ping(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	connection = pgConn
}

func Disconnect() {
	_ = connection.Close(context.TODO())
}
