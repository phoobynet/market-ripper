package config

import (
	"errors"
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/pelletier/go-toml/v2"
	"github.com/phoobynet/market-ripper/file"
	"github.com/phoobynet/market-ripper/utils"
	"github.com/questdb/go-questdb-client"
	"log"
	"os"
)

type Config struct {
	Title                       string
	Symbols                     []string
	Class                       alpaca.AssetClass
	DBHost                      string `toml:"db_host"`
	DBILPPort                   string `toml:"db_ilp_port"`
	DBPGPort                    string `toml:"db_pg_port"`
	SnapshotRefreshIntervalMins int    `toml:"snapshot_refresh_interval_mins"`
}

func Load(configPath string) *Config {
	file.MustExist(configPath)

	var config *Config

	data, err := os.ReadFile(configPath)

	if err != nil {
		log.Fatal(err)
	}

	err = toml.Unmarshal(
		data,
		&config,
	)

	if err != nil {
		log.Fatal(err)
	}

	config.clean()

	if err := config.validate(); err != nil {
		log.Fatal(err)
	}

	return config
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"title: %s, class: %s, symbols: %d, db_host: %s, db_ilp_port: %s, db_pg_port: %s, snapshot_refresh_interval_mins: %d",
		c.Title,
		c.Class,
		len(c.Symbols),
		c.DBHost,
		c.DBILPPort,
		c.DBPGPort,
		c.SnapshotRefreshIntervalMins,
	)
}

func (c *Config) GetIngressAddress() questdb.LineSenderOption {
	return questdb.WithAddress(
		fmt.Sprintf(
			"%s:%s",
			c.DBHost,
			c.DBILPPort,
		),
	)
}

func (c *Config) GetPGAddress() string {
	return fmt.Sprintf(
		"postgresql://admin:quest@%s:%s/qdb",
		c.DBHost,
		c.DBPGPort,
	)
}

// clean removes any invalid characters from the ticker symbols, trims whitespace and converts to uppercase.
func (c *Config) clean() {
	var cleanedSymbols []string
	var cleanedSymbol string

	for _, symbol := range c.Symbols {
		cleanedSymbol = utils.CleanTicker(symbol)
		if cleanedSymbol != "" {
			cleanedSymbols = append(
				cleanedSymbols,
				cleanedSymbol,
			)
		}
	}

	c.Symbols = cleanedSymbols
}

func (c *Config) validate() error {
	if c.Title == "" {
		return errors.New("title is required")
	}

	if c.Class == "" {
		return errors.New("class is required, either 'us_equity' or 'crypto'")
	} else if (c.Class != "us_equity") && (c.Class != "crypto") {
		return errors.New("class must be either 'us_equity' or 'crypto'")
	}

	if len(c.Symbols) == 0 {
		return errors.New("symbols is required")
	}

	if c.DBHost == "" {
		return errors.New("db_host is required")
	}

	if c.DBILPPort == "" {
		return errors.New("db_ilp_port is required")
	}

	if c.DBPGPort == "" {
		return errors.New("db_pg_port is required")
	}

	return nil
}
