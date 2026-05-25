package runner

import (
	"context"
	"log/slog"
	"sync"
)

type Runnable func(ctx context.Context) error

type RunMap map[string]Runnable

// Start runs every Runnable in its own goroutine. If any task returns an
// error, the shared context is cancelled so the remaining tasks can shut
// down cleanly. Start blocks until all tasks have returned and reports the
// first error observed (if any).
func (r *RunMap) Start(ctx context.Context, logger *slog.Logger) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		wg       sync.WaitGroup
		errOnce  sync.Once
		firstErr error
	)

	for name, fn := range *r {
		wg.Add(1)
		go func(name string, fn Runnable) {
			defer wg.Done()
			logger.InfoContext(ctx, "starting runner", slog.String("task", name))
			err := fn(ctx)
			if err != nil {
				firstErr = err
			}
			if err := fn(ctx); err != nil {
				errOnce.Do(func() {
					firstErr = err
				})
				cancel()
			}
		}(name, fn)
	}

	wg.Wait()
	return firstErr
}
