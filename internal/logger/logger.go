package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapService interface {
	Info(fields ...interface{})
}

type ZapStorage struct {
	*zap.Logger
}

var Zap ZapService = &ZapStorage{zap.NewNop()}

//var Log *zap.Logger = zap.NewNop()

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return fmt.Errorf("parse atomic level error: %w", err)
	}

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Level = lvl
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	logger, err := config.Build()
	if err != nil {
		return fmt.Errorf("logger build error: %w", err)
	}

	Zap = &ZapStorage{logger}

	return nil
}

func (z *ZapStorage) Info(fields ...interface{}) {

	z.Logger.Sugar().Infoln(fields...)

	//Log.Info(message, fields...)
}

/*func Debug(message string, fields ...zap.Field) {
	Log.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	Log.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	Log.Fatal(message, fields...)
}*/
