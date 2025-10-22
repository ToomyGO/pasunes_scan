package pipeline

import (
	"context"
	"sync"
)

type Stage[T any] func(<-chan T) <-chan T

// PipeLine مراحل مختلف پردازش رو به هم وصل می‌کنه.
func Pipeline[T any](in <-chan T, stages ...Stage[T]) <-chan T {
	out := in

	for _, stage := range stages {
		out = stage(out)
	}
	return out
}

// FanOut چند worker ایجاد می‌کنه تا هم‌زمان روی یه ورودی کار کنن.
func FanOut[T any](in <-chan T, num int, stage Stage[T]) []<-chan T {
	outs := make([]<-chan T, num)

	for i := 0; i < num; i++ {
		outs[i] = stage(in)
	}
	return outs
}

// FanIn خروجی چند goroutine مختلف رو با هم ترکیب می‌کنه.
func FanIn[T any](ctx context.Context, chans ...<-chan T) <-chan T {
	out := make(chan T)

	if len(chans) == 0 {
		close(out)
		return out
	}

	var wg sync.WaitGroup
	wg.Add(len(chans))

	multiplex := func(c <-chan T) {
		defer wg.Done()

		for v := range c {
			select {
			case out <- v:
			case <-ctx.Done():
				return
			}
		}
	}

	for _, c := range chans {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func Merge[T any](ctx context.Context, channels []<-chan T) <-chan T {
	return FanIn(ctx, channels...)
}
