package io

import (
	"io/ioutil"
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
	if os.Getenv("DEBUG") == "true" {
		Logger.SetOutput(os.Stdout)
	} else {
		Logger.SetOutput(ioutil.Discard)
	}
}
