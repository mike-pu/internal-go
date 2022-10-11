package logging

import (
	"log"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log - Global logger
var Log = zap.NewNop()
var SugaredLog = Log.Sugar()

const logFilePath = "/data/logs/logs.log"

var initOnce sync.Once

func Init(level zapcore.Level) {
	initOnce.Do(func() {
		nonErrorLevelEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool { return l >= level && l < zapcore.ErrorLevel })

		// Configure console output.
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

		logFileSyncer, err := createLogFileSyncer(logFilePath)
		if err != nil {
			log.Fatalln("Failed to start zap logger")
		}

		// Join the outputs, encoders, and level-handling functions into zapcore.
		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stderr), zapcore.ErrorLevel),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), nonErrorLevelEnabler),
			zapcore.NewCore(fileEncoder, logFileSyncer, zapcore.DebugLevel),
		)

		Log = zap.New(core)
		SugaredLog = Log.Sugar()

		// Redirecting standard logging to zap
		zap.RedirectStdLog(Log)
	})
}

func createLogFileSyncer(logFilePath string) (zapcore.WriteSyncer, error) {
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil, err
	}

	return zapcore.AddSync(logFile), nil
}