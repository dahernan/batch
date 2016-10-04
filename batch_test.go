package batch

import (
	"context"
	"fmt"
	"testing"

	"github.com/cheekybits/is"
)

type Item struct {
	field string
}

func Batcher(ctx context.Context, batchSize int, doFn func(ctx context.Context, items []interface{})) chan interface{} {
	buffer := make([]interface{}, 0, batchSize)
	sendCh := make(chan interface{})

	go func() {
		for {
			select {
			case item := <-sendCh:
				buffer = append(buffer, item)
				if len(buffer) >= batchSize {
					doFn(ctx, buffer)
					buffer = make([]interface{}, 0, batchSize)
				}
			case <-ctx.Done():
				doFn(ctx, buffer)
				buffer = nil
				return
			}
		}
	}()
	return sendCh
}

func TestBatch(t *testing.T) {
	is := is.New(t)

	ctx := context.Background()
	b := Batcher(ctx, 3, func(ctx context.Context, items []interface{}) {
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
