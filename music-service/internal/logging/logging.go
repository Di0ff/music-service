package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger = logrus.New()

func Init() {
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	level, exists := os.LookupEnv("LOG_LEVEL")
	if exists {
		logLevel, err := logrus.ParseLevel(level)
		if err == nil {
			Logger.SetLevel(logLevel)
		} else {
			Logger.Warn("Invalid value for the logging level")
			Logger.SetLevel(logrus.InfoLevel)
		}
	} else {
		Logger.SetLevel(logrus.InfoLevel)
	}

	Logger.SetOutput(os.Stdout)
}
