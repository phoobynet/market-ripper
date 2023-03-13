package writer

import (
	"context"
	"fmt"
	"github.com/phoobynet/market-ripper/config"
	"github.com/phoobynet/market-ripper/types"
	"github.com/questdb/go-questdb-client"
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
	lineSender       *questdb.LineSender
}

func NewTradeWriter(configuration *config.Config) *TradeWriter {
	sender, err := questdb.NewLineSender(context.TODO(), configuration.GetIngressAddress())

	if err != nil {
		log.Fatalf("Error initializing lineSender: %v", err)
	}

	writeTicker := time.NewTicker(time.Second)
	writeChan := make(chan []types.Trade, 10_000)

	logTicker := time.NewTicker(time.Second * 5)

	tradeWriter := &TradeWriter{
		writeTicker: writeTicker,
		writeChan:   writeChan,
		logTicker:   logTicker,
		tableName:   fmt.Sprintf("%s_trades", configuration.Class),
		lineSender:  sender,
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

	return tradeWriter
}

func (b *TradeWriter) Write(trade types.Trade) {
	b.writeLock.Lock()
	defer b.writeLock.Unlock()

	b.inputBuffer = append(b.inputBuffer, trade)
}

func (b *TradeWriter) Close() {
	b.writeTicker.Stop()
	b.logTicker.Stop()
	_ = b.lineSender.Close()
}

func (b *TradeWriter) copyBuffer() {
	b.writeLock.Lock()
	defer b.writeLock.Unlock()

	tempBuffer := make([]types.Trade, len(b.inputBuffer))
	copy(tempBuffer, b.inputBuffer)

	b.inputBuffer = make([]types.Trade, 0)

	b.writeChan <- tempBuffer
}

func (b *TradeWriter) flush(trades []types.Trade) {
	b.writtenCountLock.Lock()
	defer b.writtenCountLock.Unlock()
	var err error

	chunks := lo.Chunk(trades, 1_000)

	var c int64

	ctx := context.TODO()

	for _, chunkOfTrades := range chunks {
		for _, trade := range chunkOfTrades {
			err = b.lineSender.Table(b.tableName).
				Symbol("ticker", trade.Symbol).
				Float64Column("price", trade.Price).
				Float64Column("size", trade.Size).
				StringColumn("taker_side", trade.TakerSide).
				TimestampColumn("trade_timestamp", trade.Timestamp.UnixMicro()).
				AtNow(ctx)

			if err != nil {
				log.Fatal(err)
			}

			c++
		}

		err = b.lineSender.Flush(ctx)

		if err != nil {
			log.Fatal(err)
		}
	}

	b.writtenCount += c
}
