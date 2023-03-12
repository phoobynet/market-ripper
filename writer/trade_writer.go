package writer

import (
	"context"
	"github.com/phoobynet/market-ripper/types"
	"github.com/samber/lo"
	"log"
	"sync"
	"time"
)

type TradeWriter struct {
	inputBuffer      []types.Trade
	writeTicker      *time.Ticker
	writeLock        sync.RWMutex
	writeChan        chan []types.Trade
	writtenCount     int64
	writtenCountLock sync.RWMutex
	logTicker        *time.Ticker
	tableName        string
}

func NewTradeWriter() *TradeWriter {
	writeTicker := time.NewTicker(5 * time.Second)
	writeChan := make(chan []types.Trade, 10_000)

	logTicker := time.NewTicker(time.Second * 5)

	tradeWriter := &TradeWriter{
		writeTicker: writeTicker,
		writeChan:   writeChan,
		logTicker:   logTicker,
	}

	go func() {
		for {
			select {
			case <-writeTicker.C:
				tradeWriter.copyBuffer()
			case trades := <-writeChan:
				tradeWriter.flush(trades)
			case <-logTicker.C:
				tradeWriter.writtenCountLock.RLock()
				log.Printf("Trades: %d", tradeWriter.writtenCount)
				tradeWriter.writtenCountLock.RUnlock()
			}
		}
	}()

	return &TradeWriter{
		writeTicker: writeTicker,
		writeChan:   writeChan,
		logTicker:   logTicker,
	}
}

func (b *TradeWriter) Write(trade types.Trade) {
	b.writeLock.Lock()
	defer b.writeLock.Unlock()

	b.inputBuffer = append(b.inputBuffer, trade)
}

func (b *TradeWriter) Close() {
	b.writeTicker.Stop()
	b.logTicker.Stop()
}

func (b *TradeWriter) copyBuffer() {
	b.writeLock.Lock()
	defer b.writeLock.Unlock()

	tempBuffer := make([]types.Trade, len(b.inputBuffer))
	copy(tempBuffer, b.inputBuffer)

	// Clear the input buffer
	b.inputBuffer = make([]types.Trade, 0)

	// Send the buffer to the write channel
	b.writeChan <- tempBuffer
}

func (b *TradeWriter) flush(Trades []types.Trade) {
	var err error

	if b.tableName == "" {
		if len(Trades) > 0 {
			if Trades[0].Class == "c" {
				b.tableName = "crypto_trades"
			} else {
				b.tableName = "equity_trades"
			}
		}
	}

	chunks := lo.Chunk(Trades, 1_000)

	var c int64

	ctx := context.TODO()

	for _, chunkOfTrades := range chunks {
		for _, trade := range chunkOfTrades {
			err = lineSender.Table(b.tableName).
				Symbol("ticker", trade.Symbol).
				Float64Column("price", trade.Price).
				Float64Column("size", trade.Size).
				TimestampColumn("trade_timestamp", trade.Timestamp.UnixMicro()).
				AtNow(ctx)

			if err != nil {
				log.Fatal(err)
			}

			c++
		}

		err = lineSender.Flush(ctx)

		if err != nil {
			log.Fatal(err)
		}
	}

	b.writtenCountLock.Lock()
	b.writtenCount += c
	defer b.writtenCountLock.Unlock()
}
