package bootstrap

import (
	"ble2mqtt/internal/config"
	"log/slog"
	"os"
)

var levels map[string]slog.Level

func init() {
	levels = map[string]slog.Level{
		"debug": slog.LevelDebug, // Raw data arriving, regardless if it's used or not.
		"info":  slog.LevelInfo,  // Logs data being sent to the MQTT server.
		"warn":  slog.LevelWarn,  // MQTT server disconnected, sensors are sending bad packets.
		"error": slog.LevelError, // MQTT server cannot be reconnected after a disconnect
	}
}

func GetLogger(config *config.Config) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: levels[config.Logging.Level]}))
	return logger
}
