package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// InitLogger initializes a new instance of logrus.Logger for later consumption
func init() {
	logger = logrus.New()
	logger.SetFormatter(new(logrus.TextFormatter))
	logger.Formatter.(*logrus.TextFormatter).FullTimestamp = true
	logger.SetReportCaller(true)
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(os.Stdout)
}

// GetLogger retrieves the initialized instance (calls InitLogger if not already initialized)
func GetLogger() *logrus.Logger {
	return logger
}
