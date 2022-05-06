package logger

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	options *LogrusOptions = &LogrusOptions{
		TimestampFormat: time.RFC3339,
		MessageKeyName:  logrus.FieldKeyMsg,
		LevelKeyName:    logrus.FieldKeyLevel,
		TimeKeyName:     logrus.FieldKeyTime,
		FuncKeyName:     logrus.FieldKeyFunc,
		FileKeyName:     logrus.FieldKeyFile,
		LogLevel:        logrus.DebugLevel.String(),
	}
)

func Init(opts *LogrusOptions) {
	if opts == nil {
		opts = options
	}

	var logLevel logrus.Level
	var err error
	if logLevel, err = logrus.ParseLevel(opts.LogLevel); err != nil {
		panic(fmt.Errorf("invalid log level: %w", err))
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:   opts.MessageKeyName,
			logrus.FieldKeyLevel: opts.LevelKeyName,
			logrus.FieldKeyTime:  opts.TimeKeyName,
			logrus.FieldKeyFunc:  opts.FuncKeyName,
			logrus.FieldKeyFile:  opts.FileKeyName,
		},
	})

	logrus.SetLevel(logLevel)
}
