## Batch library with batch size, and flush interval

Usage:

Please do copy & paste the code or [drop it](https://github.com/matryer/drop), don't create more dependencies :)

```
ctx, cancel := context.WithCancel(context.Background())

// batch size of 3, and flush the batch every 4 Seconds even if the batch is not full
b, doneCh := Batcher(ctx, 3, 4 * time.Second, func(ctx context.Context, items []interface{}) {
  fmt.Printf("Some batch %+v\n", items)
})

b <- "one"
b <- "two"
b <- "three"

b <- "four"
b <- "five"
b <- "six"

cancel()
<-doneCh
```
