package Log

import (
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strconv"
)

var logger = logrus.New()
var Entry *logrus.Entry

func init() {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	logger.Formatter = &logrus.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		FullTimestamp:    true,
		CallerPrettyfier: callerFormatter,
	}
	logger.SetReportCaller(true)
	Entry = logger.WithFields(logrus.Fields{})
}

func callerFormatter(frame *runtime.Frame) (function string, file string) {
	line := frame.Line
	funcName := "[" + frame.Function + "]" + "[" + strconv.Itoa(line) + "]"
	return funcName, ""
}

func D() func(args ...interface{}) {
	logger.SetReportCaller(true)
	return Entry.Debug
}

func W() func(args ...interface{}) {
	logger.SetReportCaller(true)
	return Entry.Warn
}

func E() func(args ...interface{}) {
	logger.SetReportCaller(true)
	return Entry.Error
}

func I() func(args ...interface{}) {
	logger.SetReportCaller(false)
	return Entry.Info
}
