package logger

import (
	"github.com/serhiq/skye-trading-bot/internal/config"
	gelf "github.com/snovichkov/zap-gelf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var SugaredLogger *zap.SugaredLogger

func InitLogger(cfg config.Config) (err error) {
	loggingLevel := zap.DebugLevel

	consoleCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stderr,
		zap.NewAtomicLevelAt(loggingLevel),
	)

	gelfCore, err := gelf.NewCore(
		gelf.Addr(cfg.Telemetry.GraylogPath),
		gelf.Level(loggingLevel),
	)

	notSugaredLogger := zap.New(zapcore.NewTee(consoleCore, gelfCore))

	SugaredLogger = notSugaredLogger.Sugar().With(
		"service", cfg.Project.ServiceName,
	)
	//
	//if zapConfig == nil {
	//	log, err = zap.NewDevelopment(zap.AddCallerSkip(1))
	//} else {
	//	if zapConfig.Encoding == "json" {
	//		zapConfig.EncoderConfig = zap.NewProductionEncoderConfig()
	//		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//	} else {
	//		zapConfig.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	//	}
	//	log, err = zapConfig.Build(zap.AddCallerSkip(1))
	//}
	//if err != nil {
	//	return
	//}
	//s.logger = log.Sugar()
	return nil
}

func Sync() {
	SugaredLogger.Info("syncing zap logger")
	_ = SugaredLogger.Sync()
}
