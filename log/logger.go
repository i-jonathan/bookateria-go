package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLog() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := config.Build()
	return logger
}

func Handler(level string, text string, err error) {
	if err != nil {
		switch level {
		case "info":
			zap.S().Info(text, err)
		case "warning":
			zap.S().Warn(text, err)
		case "error":
			zap.S().Error(text, err)
		case "panic":
			zap.S().Panic(text, err)
		case "fatal":
			zap.S().Fatal(text, err)
		}
	}
}
