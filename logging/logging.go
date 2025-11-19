package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"

	"noname001/filesystem"
)

// TODO: compartmentalized logger, also try using only std logger
// TODO: support for standalone module logger with custom prepend

type WrappedLogger struct {
	*logrus.Logger
}

var (
	Logger *WrappedLogger

	loggingTzStr  string
	loggingTz     *time.Location

	loggingConfig *LoggingConfig
	loggingCron   *cron.Cron
)

func NewCompactLogger() (*WrappedLogger) {
	logger := logrus.New()
	logger.SetFormatter(new(CompactFormatter))
	logger.SetLevel(logrus.InfoLevel)

	return &WrappedLogger{logger}
}

func ConfigureLogging(cfg *LoggingConfig) {
	Logger = NewCompactLogger()

	if cfg == nil {
		Logger.Out = os.Stdout
		loggingTzStr = "UTC"
		loggingTz, _ = time.LoadLocation(loggingTzStr)
		loggingCron = nil

		return
	}


	loggingConfig = cfg

	switch cfg.LogLevel {
	case "error": Logger.SetLevel(logrus.ErrorLevel)
	case "warn": Logger.SetLevel(logrus.WarnLevel)
	case "info": Logger.SetLevel(logrus.InfoLevel)
	case "debug": Logger.SetLevel(logrus.DebugLevel)
	default: Logger.SetLevel(logrus.InfoLevel)
	}

	if cfg.LogTzStr == "" {
		fmt.Printf("logging: empty log_timezone, defaulting to UTC.\n")
		cfg.LogTzStr = "UTC"
	}

	newTz, tzErr := time.LoadLocation(cfg.LogTzStr)
	if tzErr != nil {
		fmt.Printf("logging: invalid log_timezone '%s', defaulting to UTC.\n", cfg.LogTzStr)
		loggingTzStr = "UTC"
		loggingTz, _ = time.LoadLocation(loggingTzStr)
	} else {
		fmt.Printf("logging: logging with timezone '%s'.\n", cfg.LogTzStr)
		loggingTzStr = cfg.LogTzStr
		loggingTz = newTz
	}

	switch cfg.LogTo {
	case "file":
		fmt.Printf("logging: to file, rotator and cleanup enabled.\n")
		
		doLogToFile()
		
		loggingCron := cron.New(
			// cron.WithLocation(loggingTz),
			cron.WithSeconds(),
		)
		loggingCron.AddFunc(fmt.Sprintf("CRON_TZ=%s 1 0 0 * * *", loggingTzStr),
			func() {
				Logger.Infof("logging: rotating log file...")
				doLogToFile()
			},
		)
		loggingCron.Start()

	case "stdout": fallthrough
	default:
		fmt.Printf("logging: to stdout, rotator and cleanup disabled.\n")
	}
}

func doLogToFile() {
	todate := time.Now().In(loggingTz).Format("2006-01-02")
	logFilename := "logg__" + strings.ReplaceAll(todate, "-", "_") + ".log"

	err1 := os.MkdirAll(filesystem.LogDir, filesystem.DEFAULT_DIRECTORY_PERMISSION)
	if err1 != nil {
		fmt.Printf("logging: Unable to create log directory. Will default to previous output config\n")
	} else {
		logFile, err2 := os.OpenFile(filepath.Join(filesystem.LogDir, logFilename), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err2 != nil {
			fmt.Printf("logging: Unable to create log file. Will default to previous output config\n")
		} else {
			Logger.Out = logFile
		}
	}
}


// func (logger) PrefixedError()