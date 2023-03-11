package writer

import (
	"context"
	"github.com/phoobynet/market-ripper/types"
	"github.com/samber/lo"
	"log"
	"sync"
	"time"
)

type BarWriter struct {
	inputBuffer      []types.Bar
	writeTicker      *time.Ticker
	writeLock        sync.RWMutex
	writeChan        chan []types.Bar
	writtenCount     int64
	writtenCountLock sync.RWMutex
	logTicker        *time.Ticker
	tableName        string
}

func NewBarWriter() *BarWriter {
	writeTicker := time.NewTicker(5 * time.Second)
	writeChan := make(chan []types.Bar, 10_000)

	logTicker := time.NewTicker(time.Second * 5)

	barWriter := &BarWriter{
		writeTicker: writeTicker,
		writeChan:   writeChan,
		logTicker:   logTicker,
	}

	go func() {
		for {
			select {
			case <-writeTicker.C:
				barWriter.copyBuffer()
			case bars := <-writeChan:
				barWriter.flush(bars)
			case <-logTicker.C:
				barWriter.writtenCountLock.RLock()
				log.Printf("bars: %d", barWriter.writtenCount)
				barWriter.writtenCountLock.RUnlock()
			}
		}
	}()

	return &BarWriter{
		writeTicker: writeTicker,
		writeChan:   writeChan,
		logTicker:   logTicker,
	}
}

func (b *BarWriter) Close() {
	b.writeTicker.Stop()
	b.logTicker.Stop()
}

func (b *BarWriter) copyBuffer() {
	b.writeLock.Lock()
	defer b.writeLock.Unlock()

	tempBuffer := make([]types.Bar, len(b.inputBuffer))
	copy(tempBuffer, b.inputBuffer)

	// Clear the input buffer
	b.inputBuffer = make([]types.Bar, 0)

	// Send the buffer to the write channel
	b.writeChan <- tempBuffer
}

func (b *BarWriter) flush(bars []types.Bar) {
	var err error

	if b.tableName == "" {
		if len(bars) > 0 {
			if bars[0].Class == "c" {
				b.tableName = "crypto_bars"
			} else {
				b.tableName = "equity_bars"
			}
		}
	}

	chunks := lo.Chunk(bars, 1_000)

	var c int64

	ctx := context.TODO()

	for _, chunkOfBars := range chunks {
		for _, bar := range chunkOfBars {
			err = lineSender.Table(b.tableName).
				Symbol("ticker", bar.Symbol).
				Float64Column("o", bar.Open).
				Float64Column("h", bar.High).
				Float64Column("l", bar.Low).
				Float64Column("c", bar.Close).
				Float64Column("v", bar.Volume).
				Float64Column("vw", bar.VWAP).
				Float64Column("n", bar.TradeCount).
				TimestampColumn("t", bar.Timestamp.UnixMicro()).
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