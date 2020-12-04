package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

func InitLog() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := config.Build()
	return logger
}

func Handler(err error) {
	file, issue := os.OpenFile("log/error.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if issue != nil {
		log.Printf("Error opening file: %v", err)
		return
	}

	log.SetOutput(file)
	log.Println(err)
	err = file.Close()
	log.Println(err)
	return
}