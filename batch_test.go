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
