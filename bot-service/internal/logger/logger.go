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
	return nil
}

func Sync() {
	SugaredLogger.Info("syncing zap logger")
	_ = SugaredLogger.Sync()
}
