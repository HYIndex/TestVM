package logging

import (
	"github.com/Sirupsen/logrus"
	"os"
	"testvm/conf"
)

var logger = logrus.New()

func init() {
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Level = logrus.DebugLevel
	logger.Out = os.Stdout
}

func GetLogger() *logrus.Logger {
	level, err := logrus.ParseLevel(conf.GlobalConfig().LOG_Level)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"package" : "logging",
		}).Errorf("ParseLevel fail, Error: %v\n", err)
		return nil
	}
	logger.Level = level
	return logger
}
