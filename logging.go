package main

import (
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

var levels map[string]log.Level

func init() {
	levels = map[string]log.Level{
		"debug": log.DebugLevel, // Raw data arriving, regardless if it's used or not.
		"info":  log.InfoLevel,  // Logs data being sent to the MQTT server.
		"warn":  log.WarnLevel,  // MQTT server disconnected, sensors are sending bad packets.
		"error": log.WarnLevel,  // MQTT server cannot be reconnected after a disconnect
		"fatal": log.FatalLevel, // Something else is using the Bluetooth controller
		"panic": log.PanicLevel, // Someone has unplugged the Bluetooth controller
	}
}

func getLogger(config *Config) *log.Logger {
	var output io.Writer
	if config.Logging.File == "stdout" {
		output = os.Stdout
	} else if config.Logging.File == "stderr" {
		output = os.Stderr
	} else {
		var ferr error
		output, ferr = os.OpenFile(config.Logging.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		if ferr != nil {
			panicCheck("opening file for logging", ferr)
		}
	}

	level, lok := levels[config.Logging.Level]
	if !lok {
		level = log.DebugLevel
	}

	return &log.Logger{
		Out: output,
		Formatter: &log.TextFormatter{
			DisableTimestamp: !config.Logging.Timestamp,
			FullTimestamp:    config.Logging.Timestamp,
		},
		Hooks:        make(log.LevelHooks),
		Level:        level,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
}
