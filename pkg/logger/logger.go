package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
	"os"
)

type Logger struct {
	*logrus.Logger
	file *os.File
}

func (lg *Logger) Close() error {
	err := lg.file.Close()
	if err != nil {
		return err
	}
	return nil
}

func New(fileName string, logLevel string) *Logger {
	lg := logrus.New()
	var file *os.File = nil
	if fileName == "" {
		lg.Out = os.Stdout
	} else {
		f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return nil
		}
		lg.Out = f
		file = f
	}
	if logLevel != "" {
		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			return nil
		}
		lg.SetLevel(level)
	}
	lg.SetReportCaller(true)

	formatter := &prefixed.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	}
	lg.SetFormatter(formatter)

	return &Logger{lg, file}
}
