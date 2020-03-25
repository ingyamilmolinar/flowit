package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger
var loggerInitialized = false

// InitLogger initializes a new instance of logrus.Logger for later consumption
func InitLogger() {
	logger = logrus.New()
	logger.SetFormatter(new(logrus.TextFormatter))
	logger.Formatter.(*logrus.TextFormatter).FullTimestamp = true
	logger.SetReportCaller(true)
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(os.Stdout)
	loggerInitialized = true
}

// GetLogger retrieves the initialized instance (calls InitLogger if not already initialized)
func GetLogger() *logrus.Logger {
	if !loggerInitialized {
		InitLogger()
	}
	return logger
}
