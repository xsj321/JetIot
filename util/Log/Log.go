package Log

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var logger = logrus.New()
var Entry *logrus.Entry

func init() {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	logger.Formatter = &logrus.TextFormatter{}
	logger.SetReportCaller(true)
	Entry = logger.WithFields(logrus.Fields{
		"Time": time.Now(),
	})
}

func D() func(args ...interface{}) {
	return Entry.Debug
}

func W() func(args ...interface{}) {
	return Entry.Warn
}

func E() func(args ...interface{}) {
	return Entry.Error
}

func I() func(args ...interface{}) {
	return Entry.Info
}
