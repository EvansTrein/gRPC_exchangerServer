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

var cfg *config.Config
var appLog *slog.Logger

func init() {
	cfg = config.MustLoadConf()

	appLog = logs.InitLog(cfg.Env)
}

func main() {
	db, err := sqlite.New(cfg.StoragePath, appLog)
	if err != nil {
		appLog.Error("failed to initialize database", slog.String("error", err.Error()))
		return
	}

	application := app.New(appLog, cfg.GrpcServ.Port, db)

	application.MustRatesInit()

	go func() {
		application.MustStart()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	application.Stop()
}
