package bootstrap

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// MakeSignal returns a context that is cancelled when SIGINT or SIGHUP is
// received. After the signal is received, the process is forcefully exited
// after 10 seconds to guarantee shutdown even if cleanup hangs.
func MakeSignal(logger *slog.Logger, grace time.Duration) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGHUP)

	go func() {
		sig := <-ch
		logger.Info("received signal shutting down", slog.String("signal", sig.String()))
		cancel()

		time.Sleep(grace)
		logger.Warn("shutdown grace period exceeded, exiting")
		os.Exit(1)
	}()

	return ctx
}
