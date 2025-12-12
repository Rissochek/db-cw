package postgres

import (
	"fmt"
	"time"

	"github.com/Rissochek/db-cw/internal/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Postgres struct {
	conn *sqlx.DB
}

func NewPostgres(conn *sqlx.DB) *Postgres {
	return &Postgres{
		conn: conn,
	}
}

func CreateConnection(dsn string) *sqlx.DB {
	time.Sleep(3 * time.Second)
	zap.S().Infof("waiting for db connection")

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		zap.S().Fatalf("failed connect to postgres: %v", err)
	}

	return db
}

func CreateDsnFromEnv() string {
	host := utils.GetKeyFromEnv("POSTGRES_HOST")
	user := utils.GetKeyFromEnv("POSTGRES_USER")
	password := utils.GetKeyFromEnv("POSTGRES_PASSWORD")
	dbName := utils.GetKeyFromEnv("POSTGRES_DB")
	port := utils.GetKeyFromEnv("POSTGRES_PORT")
	timezone := utils.GetKeyFromEnv("POSTGRES_TZ")
	sslMode := utils.GetKeyFromEnv("POSTGRES_SSL")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v", host, user, password, dbName, port, sslMode, timezone)
	return dsn
}
