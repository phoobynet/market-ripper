package writer

import (
	"context"
	"fmt"
	"github.com/phoobynet/market-ripper/config"
	"github.com/questdb/go-questdb-client"
	"log"
)

var lineSender *questdb.LineSender

func StartLineSender(configuration *config.Config) {
	sender, err := questdb.NewLineSender(context.TODO(), questdb.WithAddress(fmt.Sprintf("%s:%s", configuration.DBHost, configuration.DBILPPort)))

	if err != nil {
		log.Fatalf("Error initializing lineSender: %v", err)
	}

	lineSender = sender
}

func CloseLineSender() {
	_ = lineSender.Close()
}
