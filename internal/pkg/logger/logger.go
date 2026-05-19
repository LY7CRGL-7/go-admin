package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

// Init 初始化日志
func Init(logFile, level string, maxSize, maxBackups, maxAge int) error {
	// 创建日志目录
	dir := filepath.Dir(logFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	// 配置日志级别
	logLevel := zapcore.InfoLevel
	switch level {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "info":
		logLevel = zapcore.InfoLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	}
	
	// 文件日志
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    maxSize,    // MB
		MaxBackups: maxBackups,
		MaxAge:     maxAge,     // days
		Compress:   true,
	})
	
	// 控制台日志
	consoleWriter := zapcore.AddSync(os.Stdout)
	
	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	
	// 创建 Core
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		fileWriter,
		logLevel,
	)
	
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		consoleWriter,
		logLevel,
	)
	
	// 合并 Core
	core := zapcore.NewTee(fileCore, consoleCore)
	
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	
	return nil
}

// GetLogger 获取 logger 实例
func GetLogger() *zap.Logger {
	if logger == nil {
		// 默认配置
		Init("/tmp/admin.log", "info", 100, 10, 30)
	}
	return logger
}

// Sugar 获取 SugaredLogger
func Sugar() *zap.SugaredLogger {
	return GetLogger().Sugar()
}

// Debug 级别日志
func Debug(msg string, fields ...interface{}) {
	GetLogger().Sugar().Debugw(msg, fields...)
}

// Info 级别日志
func Info(msg string, fields ...interface{}) {
	GetLogger().Sugar().Infow(msg, fields...)
}

// Warn 级别日志
func Warn(msg string, fields ...interface{}) {
	GetLogger().Sugar().Warnw(msg, fields...)
}

// Error 级别日志
func Error(msg string, fields ...interface{}) {
	GetLogger().Sugar().Errorw(msg, fields...)
}

// Fatal 级别日志
func Fatal(msg string, fields ...interface{}) {
	GetLogger().Sugar().Fatalw(msg, fields...)
}
