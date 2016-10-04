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

func TestBatch(t *testing.T) {
	is := is.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	batches := 0
	b, doneCh := Batcher(ctx, 3, 0, func(ctx context.Context, items []interface{}) {
		fmt.Printf("Some batch %+v\n", items)
		batches++
	})

	b <- Item{"one"}
	b <- Item{"two"}
	b <- Item{"three"}

	b <- Item{"four"}
	b <- Item{"five"}
	b <- Item{"six"}

	cancel()
	<-doneCh
	is.Equal(batches, 2)
}

func TestBatchFlush(t *testing.T) {
	is := is.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	batches := 0
	b, doneCh := Batcher(ctx, 3, 2*time.Second, func(ctx context.Context, items []interface{}) {
		fmt.Printf("Some batch flush %+v\n", items)
		batches++
	})

	b <- Item{"one"}
	b <- Item{"two"}

	time.Sleep(3 * time.Second)
	b <- Item{"three"}

	b <- Item{"four"}
	b <- Item{"five"}
	b <- Item{"six"}

	cancel()
	<-doneCh
	is.Equal(batches, 3)

}
