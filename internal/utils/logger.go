package utils

import (
	"flag"
	"github.com/armory/armory-cli/cmd"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func GetLogger() *logrus.Logger {
	if logger == nil {
		logger := logrus.New()
		lvl := logrus.FatalLevel
		if cmd.GlobalCommandConfig.VerboseMode {
			lvl = logrus.DebugLevel
		}
		logger.SetLevel(lvl)
		logger.SetFormatter(&logrus.TextFormatter{})
		_ = flag.Set("logtostderr", "true")
	}

	return logger
}