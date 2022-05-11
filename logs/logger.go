package logs

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
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

func Init(level string, isDev bool) {

	var logLevel logrus.Level
	var err error
	if logLevel, err = logrus.ParseLevel(level); err != nil {
		panic(fmt.Errorf("invalid log level: %w", err))
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:   options.MessageKeyName,
			logrus.FieldKeyLevel: options.LevelKeyName,
			logrus.FieldKeyTime:  options.TimeKeyName,
			logrus.FieldKeyFunc:  options.FuncKeyName,
			logrus.FieldKeyFile:  options.FileKeyName,
		},
	})

	if isDev {
		logrus.SetFormatter(&prefixed.TextFormatter{
			ForceFormatting: true,
		})
	}

	logrus.SetLevel(logLevel)
}
