package fastweb

import (
	"go.uber.org/zap/zapcore"
	"io"
	"time"

	"go.uber.org/zap"
)

type Logger interface {
	Info(...interface{})
	Infof(string, ...interface{})

	Debug(...interface{})
	Debugf(string, ...interface{})

	Warn(...interface{})
	Warnf(string, ...interface{})

	Error(...interface{})
	Errorf(string, ...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
}

type zapLogger struct {
	log *zap.SugaredLogger
}

func DefaultLogger(w io.Writer, development bool) *zapLogger {
	var log *zap.Logger
	var err error

	if development {
		log, err = zap.NewDevelopment()
	} else {
		log, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}

	log.WithOptions(
		zap.Fields(
			zap.String("ts", time.Now().Format(time.RFC3339)),
		),
		zap.ErrorOutput(zapcore.NewMultiWriteSyncer(zapcore.AddSync(w))),
	)

	return &zapLogger{
		log: log.Sugar(),
	}
}

func (l *zapLogger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *zapLogger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *zapLogger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *zapLogger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *zapLogger) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l *zapLogger) Panic(args ...interface{}) {
	l.log.Panic(args...)
}

func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.log.Panicf(format, args...)
}
