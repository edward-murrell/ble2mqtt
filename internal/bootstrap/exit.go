package bootstrap

import (
	"log/slog"
	"os"
)

func ExitOnErr(err error) {
	if err != nil {
		slog.Error("Start error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
