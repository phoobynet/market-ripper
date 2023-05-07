package config

import (
	"errors"
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/pelletier/go-toml/v2"
	"github.com/phoobynet/market-ripper/file"
	"github.com/phoobynet/market-ripper/ticker"
	"github.com/questdb/go-questdb-client"
	"log"
	"os"
)

type Config struct {
	Title                       string
	Symbols                     []string
	Class                       alpaca.AssetClass
	DBQuestHost                 string `toml:"db_quest_host"`
	DBQuestILPPort              string `toml:"db_quest_ilp_port"`
	DBQuestPGPort               string `toml:"db_quest_pg_port"`
	SnapshotRefreshIntervalMins int    `toml:"snapshot_refresh_interval_mins"`
}

// Load loads a configuration file.
func Load(configPath string) (*Config, error) {
	file.MustExist(configPath)

	var config *Config

	data, err := os.ReadFile(configPath)

	if err != nil {
		return nil, err
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
		return nil, err
	}

	return config, nil
}

// String returns a string representation of the config.
func (c *Config) String() string {
	return fmt.Sprintf("%+v", *c)
}

// IngressAddress returns a questdb.LineSenderOption with the address set to the ingress address.
func (c *Config) IngressAddress() questdb.LineSenderOption {
	return questdb.WithAddress(
		fmt.Sprintf(
			"%s:%s",
			c.DBQuestHost,
			c.DBQuestILPPort,
		),
	)
}

// PGAddress returns a postgresql connection string - used for SQL (ish) queries.
// Deprecated: Use DSN() instead.
func (c *Config) PGAddress() string {
	return fmt.Sprintf(
		"postgresql://admin:quest@%s:%s/qdb",
		c.DBQuestHost,
		c.DBQuestPGPort,
	)
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=admin password=quest dbname=qdb sslmode=disable", c.DBQuestHost, c.DBQuestPGPort)
}

func (c *Config) clean() {
	var cleanedSymbols []string
	var cleanedSymbol string

	for _, symbol := range c.Symbols {
		cleanedSymbol = ticker.Clean(symbol)
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

	if c.DBQuestHost == "" {
		return errors.New("db_quest_host is required")
	}

	if c.DBQuestILPPort == "" {
		return errors.New("db_quest_ilp_port is required")
	}

	if c.DBQuestPGPort == "" {
		return errors.New("db_quest_pg_port is required")
	}

	return nil
}
