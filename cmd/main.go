package main

import (
	"context"
	"time"

	"github.com/Rissochek/db-cw/internal/faking"
	"github.com/Rissochek/db-cw/internal/repository/postgres"
	"github.com/Rissochek/db-cw/internal/service"
	"github.com/Rissochek/db-cw/internal/utils"
	"go.uber.org/zap"
)

var (
	seed = int64(42)
)

func main() {
	utils.LoadEnvFile()

	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	faker := faking.NewDataFaker(seed)

	dsn := postgres.CreateDsnFromEnv()
	conn := postgres.CreateConnection(dsn)
	repo := postgres.NewPostgres(conn)

	service := service.NewService(faker, repo)
	service.FillDatabase(ctx, seed)
}
