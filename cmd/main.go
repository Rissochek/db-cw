package main

import (
	"errors"
	"net/http"

	"github.com/Rissochek/db-cw/internal/app"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// @title DB CW
// @version 1.0
// @host localhost:8080

func main() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	logger, _ := config.Build()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	if err := app.InitApp().SetupRoutes().Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		zap.S().Errorf("failed to start server: %v", err)
	}
}
