# Market Ripper

Captures trade's and bars from [Alpaca's Market Stream](https://alpaca.markets/data), writing the data to a [QuestDB](https://questdb.io/docs/) database

In addition, snapshots and assets are downloaded on each restart.

This does not provide query services, just capture.  You will need to build another application on top of QuestDB to get anything useful out of it.

> This is an alternative to another project I worked on called [trade-ripper](https://github.com/phoobynet/trade-ripper).  [trade-ripper](https://github.com/phoobynet/trade-ripper) did not use Alpaca's Go library, and was, quite frankly, overly complicated.

## Requirements

- [Alpaca Market Data](https://alpaca.markets/data) SIP access (currently $99)
- [QuestDB](https://questdb.io/docs/) - A time series database that is very fast.

## Installation

Set the following environment variables

- `APCA_API_KEY_ID` - Your Alpaca Key
- `APCA_API_SECRET_KEY` - Your Alpaca Secret

Install the packages

```bash
go install github.com/phoobynet/market-ripper@latest
```

Create `.toml` file, and decide what symbols you would like to include.  

- `title` - Whatever you want it to be.
- `class` - Must be either `us_equity` or `crypto`
- `symbols` - e.g. `"AAPL"`, or `"BTC/USDT"` etc. Note that `*` means everything.  **The symbols must all be within a single class of financial instrument.**
- `db_host` - The host address of QuestDB
- `db_ilp_port` - ILP ingestion port; the default is `9009`
- `db_pg_port` - Postgres(ish) port; the default is `8812`

**Example `config.toml`**

```toml
# What ever you want it to be
title = "Everything"

# class should be either crypto or us_equity
class = "crypto"

# List of symbols to observe.  "*" means everything in that class.
symbols = [
    "*",
]

# The QuestDB host
db_host = "localhost"

# ILP (in-line protocol) is used for fast inserts
db_ilp_port = "9009"

# Port of queries on some DDL execution using a Postgres adapter
db_pg_port = "8812"

# Setting to 0 will cause the snapshots to update every 24 hours;
snapshot_refresh_interval_mins = 5
```

Assuming you have a QuestDB server up and running, start `market-ripper`

```bash
market-ripper --config config.toml
```

## Tables

- `assets` - A list of assets.  If assets are older than 24 hours, the table is refreshed.
- `crypto_snapshots` - Snapshots are created on start up and refreshed every hour
- `us_equity_snapshots` - Snapshots are created on start up and refreshed every hour
- `crypto_trades` - Crypto trades
- `us_equity_trades` - Equity trades
- `crypto_bars` - Crypto bars
- `us_equity_bars` - Equity bars
