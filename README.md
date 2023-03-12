# Market Ripper

Captures trade's and bars from [Alpaca's Market Stream](https://alpaca.markets/data), writing the data to a [QuestDB](https://questdb.io/docs/) database

In addition, snapshots and assets are downloaded on each restart.

This does not provide query services, just capture.  You will need to build another application on top of QuestDB to get anything useful out of it.

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
- `symbols` - e.g. `"AAPL"`, etc. Note that `*` means everything.  The symbols must all be within a single class of financial instrument.
- `db_host` - The host address of QuestDB
- `db_ilp_port` - ILP ingestion port; the default is `9009`
- `db_pg_port` - Postgres(ish) port; the default is `8812`

**Example `config.toml`**

```toml
title = "Everything"

class = "crypto"

symbols = [
    "*",
]

db_host = "localhost"
db_lip_port = "9009"
db_pg_port = "8812"
```

Assuming you have a QuestDB server up and running, start `market-ripper`

```bash
market-ripper --config config.toml
```


