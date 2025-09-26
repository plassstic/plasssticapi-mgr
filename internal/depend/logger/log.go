package logger

import (
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once sync.Once
	_log *zap.SugaredLogger
)

func buildSugaredLogger() *zap.SugaredLogger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		NameKey:       "logger",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		zap.InfoLevel,
	)

	logger := zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))

	sigint := make(chan os.Signal)
	signal.Notify(sigint, os.Interrupt)

	go func() {
		<-sigint
		_ = logger.Sync()
	}()

	return logger.Sugar()
}

func initSingleton() {
	once.Do(func() {
		_log = buildSugaredLogger()
	})
}

func GetLogger(name string) *zap.SugaredLogger {
	initSingleton()
	logName := "[ " + strings.ReplaceAll(name, ".", " -> ") + " ]"
	return _log.Named(logName)
}

func GetRawLogger(name string) *zap.Logger {
	initSingleton()
	logName := "[ " + strings.ReplaceAll(name, ".", " -> ") + " ]"
	return _log.Desugar().Named(logName)
}
