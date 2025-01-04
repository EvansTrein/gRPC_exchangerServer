package main

import (
	"log/slog"

	"github.com/EvansTrein/exchanger_gRPC/internal/app"
	"github.com/EvansTrein/exchanger_gRPC/internal/config"
	"github.com/EvansTrein/exchanger_gRPC/internal/storages/sqlite"
	"github.com/EvansTrein/exchanger_gRPC/pkg/logs"
)

var cfg *config.Config
var log *slog.Logger

func init() {
	cfg = config.MustLoadConf()

	log = logs.InitLog(cfg.Env)
}

func main() {
	db, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize database", "error", err)
		return
	}

	application := app.New(log, cfg.GrpcServ.Port, db)

	application.MustStart()
}




