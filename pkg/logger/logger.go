package logger

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
	"strings"
)

func InitLogger(level string, filePath string) (*log.Logger, error) {
	logger := log.New("")

	lvl, err := parseLevel(level)
	if err != nil {
		return nil, err
	}

	logger.SetLevel(lvl)

	if file, err := os.Open(filePath); err == nil {
		logger.SetOutput(file)
		logger.Infof("Logger output set to %s", filePath)
	} else {
		logger.Info("Logger output set to stdout")
		logger.SetOutput(os.Stdout)
	}

	return logger, nil
}

func parseLevel(level string) (log.Lvl, error) {
	switch strings.ToLower(level) {
	case "off":
		return log.OFF, nil
	case "error", "err":
		return log.ERROR, nil
	case "warn", "warning":
		return log.WARN, nil
	case "info":
		return log.INFO, nil
	case "debug":
		return log.DEBUG, nil
	}

	var l log.Lvl
	return l, fmt.Errorf("not a valid log Level: %q", level)
}
