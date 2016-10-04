package batch

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cheekybits/is"
)

type Item struct {
	field string
}

func Batcher(ctx context.Context, batchSize int, flushInterval time.Duration, doFn func(ctx context.Context, items []interface{})) chan interface{} {
	var (
		tickCh    <-chan time.Time
		ticker    *time.Ticker
		lastBatch time.Time

		buffer []interface{}
		sendCh chan interface{}
	)
	buffer = make([]interface{}, 0, batchSize)
	sendCh = make(chan interface{})

	if flushInterval != 0 {
		ticker = time.NewTicker(flushInterval)
		tickCh = ticker.C
	}

	go func() {
		for {
			select {
			case item := <-sendCh:
				buffer = append(buffer, item)
				if len(buffer) >= batchSize {
					doFn(ctx, buffer)
					lastBatch = time.Now()
					buffer = make([]interface{}, 0, batchSize)
				}
			case <-tickCh:
				nextBatch := lastBatch.Add(flushInterval)
				now := time.Now()
				if len(buffer) > 0 && now.After(nextBatch) {
					doFn(ctx, buffer)
					lastBatch = now
					buffer = make([]interface{}, 0, batchSize)
				}
			case <-ctx.Done():
				if len(buffer) > 0 {
					doFn(ctx, buffer)
					buffer = nil
				}
				return
			}
		}
	}()
	return sendCh
}

func TestBatch(t *testing.T) {
	is := is.New(t)

	ctx := context.Background()
	b := Batcher(ctx, 3, 0, func(ctx context.Context, items []interface{}) {
		fmt.Printf("Some batch %+v\n", items)
	})

	b <- Item{"one"}
	b <- Item{"two"}
	b <- Item{"three"}

	b <- Item{"four"}
	b <- Item{"five"}
	b <- Item{"six"}

	is.True(false)
}

func TestBatchFlush(t *testing.T) {
	is := is.New(t)

	ctx := context.Background()
	b := Batcher(ctx, 3, 2*time.Second, func(ctx context.Context, items []interface{}) {
		fmt.Printf("Some batch flush %+v\n", items)
	})

	b <- Item{"one"}
	b <- Item{"two"}

	time.Sleep(3 * time.Second)
	b <- Item{"three"}

	b <- Item{"four"}
	b <- Item{"five"}
	b <- Item{"six"}

	is.True(false)
}
