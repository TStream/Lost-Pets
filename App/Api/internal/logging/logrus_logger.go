package logging

import (
	"io/ioutil"
	"lostpets"
	"strings"

	"github.com/orandin/lumberjackrus"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// LogrusWrapper wraps logrus in the domain Logger interface
type LogrusWrapper struct {
	logger   *logrus.Logger
	maxDepth int
}

// EntryWrapper wraps logrus entry in the domain Logger interface
type EntryWrapper struct {
	entry    *logrus.Entry
	maxDepth int
}

type LogrusConfig struct {
	Depth  int    `json:"depth"`  // Depth is the number of levels to traverse when unwrapping an error stack.
	Level  string `json:"level"`  // Level of log statements to write: DEBUG | INFO | ERROR
	File   string `json:"file"`   // File to log to
	Stdout bool   `json:"stdOut"` // The logger will log to stdOut if true
}

//NewLogrusWrapper sets up a LogrusWrapper logger.
// Debug level logging only appears in StdOut.
// Defaults:
//  Depth: 0
// 	Level: Info
//	File: ""
//  StdOut: false
func NewLogrusWrapper(config LogrusConfig) (*LogrusWrapper, error) {
	logger := logrus.New()

	//TTY Setup
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	var logLevel logrus.Level
	switch strings.ToLower(config.Level) {
	case "debug":
		logger.SetLevel(logrus.DebugLevel) //set logger directly, only the stdout will use debug.
		logLevel = logrus.InfoLevel
		break
	case "error":
		logLevel = logrus.ErrorLevel //To keep logging simple we only support 3 levels for now
		break
	case "info":
		logLevel = logrus.InfoLevel
		break
	default:
		logLevel = logrus.InfoLevel
		break
	}

	//Log File Setup
	hook, err := lumberjackrus.NewHook(
		&lumberjackrus.LogFile{
			Filename:   config.File,
			MaxSize:    100,
			MaxBackups: 10,
			MaxAge:     365,
			Compress:   false,
			LocalTime:  false,
		},
		logLevel,
		&logrus.TextFormatter{},
		nil, //opts to send diff levels to diffrent files
	)

	if err != nil {
		return nil, err
	}

	logger.AddHook(hook)

	if !config.Stdout {
		logger.Out = ioutil.Discard
	}

	return &LogrusWrapper{logger, config.Depth}, nil
}

func (l *LogrusWrapper) Debug(message string, args ...interface{}) {
	l.logger.Debugf(message, args...)
}

func (l *LogrusWrapper) Info(message string, args ...interface{}) {
	l.logger.Infof(message, args...)
}

func (l *LogrusWrapper) Error(message string, args ...interface{}) {
	l.logger.Errorf(message, args...)
}

// UnwrapError prints an err message and it's stack trace if one is available.
// No more then maxDepth levels of the trace will be displayed, after which the
// number of remaining errors in the stack will be displayed.
func (l *LogrusWrapper) UnwrapError(err error) {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	l.logger.Error(err)
	if err, ok := errors.Cause(err).(stackTracer); ok {
		st := err.StackTrace()
		extra := len(st) - l.maxDepth - 1
		if l.maxDepth > 0 && extra > 0 {
			l.logger.Printf("\t%+v", st[0:l.maxDepth])
			l.logger.Printf(" .....and %d more\n", extra)
			return
		}

		l.logger.Printf("\t%+v", st[0:])
	}
}

// WithFields creates a new LogEntryWrapper containing the fields to log.
// The fields will not be logged until Debug | Info or Error are called
func (l *LogrusWrapper) WithFields(fields map[string]interface{}) lostpets.Logger {
	return newEntryWrapper(l.logger.WithFields(map[string]interface{}(fields)), l.maxDepth)
}

func newEntryWrapper(entry *logrus.Entry, maxDepth int) *EntryWrapper {
	return &EntryWrapper{entry, maxDepth}
}

func (e *EntryWrapper) Debug(message string, args ...interface{}) {
	e.entry.Debugf(message, args...)
}

func (e *EntryWrapper) Info(message string, args ...interface{}) {
	e.entry.Infof(message, args...)
}

func (e *EntryWrapper) Error(message string, args ...interface{}) {
	e.entry.Errorf(message, args...)
}

func (e *EntryWrapper) UnwrapError(err error) {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	e.entry.Error(err.Error())
	if err, ok := errors.Cause(err).(stackTracer); ok {
		st := err.StackTrace()
		extra := len(st) - e.maxDepth - 1
		if e.maxDepth > 0 && extra > 0 {
			e.entry.Printf("%+v", st[0:e.maxDepth])
			e.entry.Printf(" .....and %d more\n", extra)
			return
		}

		e.entry.Printf("%+v", st[0:])
	}
}
