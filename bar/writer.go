package bar

import (
	"context"
	"fmt"
	"github.com/phoobynet/market-ripper/bar/models"
	"github.com/phoobynet/market-ripper/config"
	"github.com/questdb/go-questdb-client"
	"github.com/samber/lo"
	"log"
	"sync"
	"time"
)

type Writer struct {
	inputBuffer      []models.Bar
	writeTicker      *time.Ticker
	writeLock        sync.RWMutex
	writeChan        chan []models.Bar
	writtenCount     int64
	writtenCountLock sync.RWMutex
	logTicker        *time.Ticker
	tableName        string
	lineSender       *questdb.LineSender
}

func NewWriter(configuration *config.Config) *Writer {
	sender, err := questdb.NewLineSender(context.TODO(), configuration.IngressAddress())

	if err != nil {
		log.Fatal(err)
	}

	writeTicker := time.NewTicker(5 * time.Second)
	writeChan := make(chan []models.Bar, 10_000)

	logTicker := time.NewTicker(time.Second * 5)

	barWriter := &Writer{
		writeTicker: writeTicker,
		writeChan:   writeChan,
		logTicker:   logTicker,
		tableName:   fmt.Sprintf("%s_bars", configuration.Class),
		lineSender:  sender,
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

	return barWriter
}

func (b *Writer) Write(bar models.Bar) {
	b.writeLock.Lock()
	defer b.writeLock.Unlock()

	b.inputBuffer = append(b.inputBuffer, bar)
}

func (b *Writer) Close() {
	b.writeTicker.Stop()
	b.logTicker.Stop()
	_ = b.lineSender.Close()
}

func (b *Writer) copyBuffer() {
	b.writeLock.Lock()
	defer b.writeLock.Unlock()

	tempBuffer := make([]models.Bar, len(b.inputBuffer))
	copy(tempBuffer, b.inputBuffer)

	// Clear the input buffer
	b.inputBuffer = make([]models.Bar, 0)

	// Send the buffer to the write channel
	b.writeChan <- tempBuffer
}

func (b *Writer) flush(bars []models.Bar) {
	var err error

	chunks := lo.Chunk(bars, 1_000)

	var c int64

	ctx := context.TODO()

	for _, chunkOfBars := range chunks {
		for _, theBar := range chunkOfBars {
			err = b.lineSender.Table(b.tableName).
				Symbol("ticker", theBar.Symbol).
				Float64Column("o", theBar.Open).
				Float64Column("h", theBar.High).
				Float64Column("l", theBar.Low).
				Float64Column("c", theBar.Close).
				Float64Column("v", theBar.Volume).
				Float64Column("vw", theBar.VWAP).
				Float64Column("n", theBar.TradeCount).
				TimestampColumn("t", theBar.Timestamp.UnixMicro()).
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

	b.writtenCountLock.Lock()
	b.writtenCount += c
	b.writtenCountLock.Unlock()
}
