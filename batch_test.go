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

type Batcher struct {
	ctx  context.Context
	doFn func(items []Item)

	buffer    []Item
	batchSize int
	sendCh    chan Item
}

func NewBatcher(ctx context.Context, batchSize int, doFn func(items []Item)) *Batcher {
	b := Batcher{
		ctx:       ctx,
		doFn:      doFn,
		batchSize: batchSize,
		buffer:    make([]Item, 0, batchSize),
		sendCh:    make(chan Item),
	}

	go b.run()
	return &b
}

func (b *Batcher) Send(item Item) {
	b.sendCh <- item
}

func (b *Batcher) run() {
	for {
		select {
		case item := <-b.sendCh:
			b.buffer = append(b.buffer, item)
			if len(b.buffer) >= b.batchSize {
				b.do()
			}
		case <-b.ctx.Done():
			b.do()
			return
		}
	}
}

func (b *Batcher) do() {
	b.doFn(b.buffer)
	b.buffer = make([]Item, 0, b.batchSize)
}

func TestBatch(t *testing.T) {
	is := is.New(t)

	b := NewBatcher(context.Background(), 3, func(items []Item) {
		fmt.Printf("Some batch %+v\n", items)
	})

	b.Send(Item{"one"})
	b.Send(Item{"two"})
	b.Send(Item{"three"})

	b.Send(Item{"four"})
	b.Send(Item{"five"})
	b.Send(Item{"six"})

	is.True(false)
}
