package batch

import (
	"context"
	"time"
)

func Batcher(ctx context.Context, batchSize int, flushInterval time.Duration, fn func(ctx context.Context, items []interface{})) chan<- interface{} {
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
					fn(ctx, buffer)
					lastBatch = time.Now()
					buffer = make([]interface{}, 0, batchSize)
				}
			case <-tickCh:
				nextBatch := lastBatch.Add(flushInterval)
				now := time.Now()
				if len(buffer) > 0 && now.After(nextBatch) {
					fn(ctx, buffer)
					lastBatch = now
					buffer = make([]interface{}, 0, batchSize)
				}
			case <-ctx.Done():
				if len(buffer) > 0 {
					fn(ctx, buffer)
					buffer = nil
				}
				return
			}
		}
	}()
	return sendCh
}
