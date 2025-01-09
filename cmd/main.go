package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/EvansTrein/gRPC_exchangerServer/internal/app"
	"github.com/EvansTrein/gRPC_exchangerServer/internal/config"
	"github.com/EvansTrein/gRPC_exchangerServer/internal/storages/sqlite"
	"github.com/EvansTrein/gRPC_exchangerServer/pkg/logs"
)

func main() {
	var cfg *config.Config
	var appLog *slog.Logger

	cfg = config.MustLoadConf()
	appLog = logs.InitLog(cfg.Env)

	db, err := sqlite.New(cfg.StoragePath, appLog)
	if err != nil {
		appLog.Error("failed to initialize database", "error", err)
		return
	}

	application := app.New(appLog, cfg.GrpcServ.Port, db, cfg.GrpcServ.ConnectionTimeout)

	application.MustRatesInit()

	go func() {
		application.MustStart()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	application.Stop()
}
