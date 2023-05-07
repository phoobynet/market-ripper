package trade

import (
	"context"
	"fmt"
	"github.com/phoobynet/market-ripper/config"
	"github.com/questdb/go-questdb-client"
	"github.com/samber/lo"
	"log"
	"sync"
	"time"
)

type Writer struct {
	inputBuffer      []Trade
	writeTicker      *time.Ticker
	writeLock        sync.RWMutex
	writeChan        chan []Trade
	writtenCount     int64
	writtenCountLock sync.RWMutex
	logTicker        *time.Ticker
	tableName        string
	lineSender       *questdb.LineSender
}

func NewWriter(configuration *config.Config) *Writer {
	sender, err := questdb.NewLineSender(context.TODO(), configuration.IngressAddress())

	if err != nil {
		log.Fatalf("Error initializing lineSender: %v", err)
	}

	writeTicker := time.NewTicker(time.Second)
	writeChan := make(chan []Trade, 10_000)

	logTicker := time.NewTicker(time.Second * 5)

	tradeWriter := &Writer{
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

func (b *Writer) Write(theTrade Trade) {
	b.writeLock.Lock()
	defer b.writeLock.Unlock()

	b.inputBuffer = append(b.inputBuffer, theTrade)
}

func (b *Writer) Close() {
	b.writeTicker.Stop()
	b.logTicker.Stop()
	_ = b.lineSender.Close()
}

func (b *Writer) copyBuffer() {
	b.writeLock.Lock()
	defer b.writeLock.Unlock()

	tempBuffer := make([]Trade, len(b.inputBuffer))
	copy(tempBuffer, b.inputBuffer)

	b.inputBuffer = make([]Trade, 0)

	b.writeChan <- tempBuffer
}

func (b *Writer) flush(trades []Trade) {
	b.writtenCountLock.Lock()
	defer b.writtenCountLock.Unlock()
	var err error

	chunks := lo.Chunk(trades, 1_000)

	var c int64

	ctx := context.TODO()

	for _, chunkOfTrades := range chunks {
		for _, theTrade := range chunkOfTrades {
			err = b.lineSender.Table(b.tableName).
				Symbol("ticker", theTrade.Symbol).
				Float64Column("price", theTrade.Price).
				Float64Column("size", theTrade.Size).
				StringColumn("taker_side", theTrade.TakerSide).
				TimestampColumn("trade_timestamp", theTrade.Timestamp.UnixMicro()).
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
