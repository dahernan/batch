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

// type Batcher struct {
// 	ctx  context.Context
// 	doFn func(items []Item)
//
// 	buffer    []Item
// 	batchSize int
// 	sendCh    chan Item
// }

func Batcher(ctx context.Context, batchSize int, doFn func(items []Item)) chan Item {
	buffer := make([]Item, 0, batchSize)
	sendCh := make(chan Item)

	go func() {
		for {
			select {
			case item := <-sendCh:
				buffer = append(buffer, item)
				if len(buffer) >= batchSize {
					doFn(buffer)
					buffer = make([]Item, 0, batchSize)
				}
			case <-ctx.Done():
				doFn(buffer)
				return
			}
		}
	}()
	return sendCh
}

func TestBatch(t *testing.T) {
	is := is.New(t)

	b := Batcher(context.Background(), 3, func(items []Item) {
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
