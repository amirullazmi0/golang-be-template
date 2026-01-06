package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

type LoggerConfig struct {
	Debug      bool
	LogToFile  bool
	LogFilePath string
	MaxSize    int // megabytes
	MaxBackups int
	MaxAge     int // days
	Compress   bool
}

// InitLogger initializes the global logger with Grafana Loki compatible format
func InitLogger(cfg LoggerConfig) error {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Set log level
	var level zapcore.Level
	if cfg.Debug {
		level = zapcore.DebugLevel
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		level = zapcore.InfoLevel
	}

	// Create cores for different outputs
	var cores []zapcore.Core

	// Console output
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	if !cfg.Debug {
		// Use JSON encoder for production (Loki compatible)
		consoleEncoder = zapcore.NewJSONEncoder(encoderConfig)
	}
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		level,
	)
	cores = append(cores, consoleCore)

	// File output (always JSON for Loki/Promtail)
	if cfg.LogToFile {
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		
		// Create logs directory if not exists
		if err := os.MkdirAll("logs", 0755); err != nil {
			return err
		}

		logFile, err := os.OpenFile(cfg.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		fileCore := zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(logFile),
			level,
		)
		cores = append(cores, fileCore)
	}

	// Combine all cores
	core := zapcore.NewTee(cores...)

	// Create logger with caller info
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	Log = logger
	return nil
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Info(msg, fields...)
	}
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Error(msg, fields...)
	}
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Debug(msg, fields...)
	}
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Warn(msg, fields...)
	}
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Fatal(msg, fields...)
	} else {
		os.Exit(1)
	}
}

// Sync flushes any buffered log entries
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
