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

func Handler(level string, text string, err []string) {
	if err != nil {
		switch level {
		case "info":
			zap.S().Infow(text, "Message", err)
		case "warning":
			zap.S().Warnw(text, "Message", err)
		case "error":
			zap.S().Errorw(text, "Message", err)
		case "panic":
			zap.S().Panicw(text, "Message", err)
		case "fatal":
			zap.S().Fatalw(text, "Message", err)
		}
	}
}