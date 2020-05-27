package io

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is the pointer to the already configured logrus logger instance
var Logger *logrus.Logger

// InitLogger initializes a new instance of logrus.Logger for later consumption
func init() {
	Logger = logrus.New()
	Logger.SetFormatter(new(logrus.TextFormatter))
	Logger.Formatter.(*logrus.TextFormatter).FullTimestamp = true
	Logger.SetReportCaller(true)
	Logger.SetLevel(logrus.DebugLevel)
	Logger.SetOutput(os.Stdout)
}
