package common

import (
	"os"

	"github.com/Sirupsen/logrus"
)

func init() {
	logrus.SetOutput(os.Stderr)
	switch ConfigString("LOG_LEVEL") {
	default:
		logrus.SetLevel(logrus.ErrorLevel)
	case "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	}
}

type GormLogger struct{}

func (*GormLogger) Print(v ...interface{}) {
	if v[0] == "sql" {
		logrus.WithFields(logrus.Fields{"module": "gorm", "type": "sql"}).Debug(v[3])
	}
	if v[0] == "log" {
		logrus.WithFields(logrus.Fields{"module": "gorm", "type": "log"}).Info(v[2])
	}
}
