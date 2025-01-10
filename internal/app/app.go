package app

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/EvansTrein/gRPC_exchangerServer/internal/config"
	"github.com/EvansTrein/gRPC_exchangerServer/internal/server"
	"github.com/EvansTrein/gRPC_exchangerServer/internal/storages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	db         storages.Database
	conf       *config.GrpcServer
}

// new application creation
func New(log *slog.Logger, db storages.Database, conf *config.GrpcServer) *App {

	ka := keepalive.ServerParameters{
		MaxConnectionIdle:     conf.MaxConnectionIdle,     // Максимальное время бездействия соединения
		MaxConnectionAge:      conf.MaxConnectionAge,      // Максимальное время жизни соединения
		MaxConnectionAgeGrace: conf.MaxConnectionAgeGrace, // Время для завершения активных запросов
		Time:                  conf.Time,                  // Время между keepalive ping
		Timeout:               conf.Timeout,               // Таймаут на ответ от клиента
	}

	gRPC := grpc.NewServer(grpc.KeepaliveParams(ka))

	server.RegisterServ(gRPC, db, log)

	return &App{
		log:        log,
		gRPCServer: gRPC,
		db:         db,
		conf:       conf,
	}
}

// application start
func (a *App) MustStart() {
	log := a.log.With(
		slog.Int("Port", a.conf.Port),
		slog.String("MaxConnectionIdle", a.conf.MaxConnectionIdle.String()),
		slog.String("MaxConnectionAge", a.conf.MaxConnectionAge.String()),
		slog.String("MaxConnectionAgeGrace", a.conf.MaxConnectionAgeGrace.String()),
		slog.String("Time", a.conf.Time.String()),
		slog.String("Timeout", a.conf.Timeout.String()),
	)
	log.Debug("started gRPC server")

	portListen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.conf.Port))
	if err != nil {
		panic(err.Error())
	}

	a.log.Info(" gRPC server successfully started", slog.String("port", portListen.Addr().String()))
	if err := a.gRPCServer.Serve(portListen); err != nil {
		panic(err.Error())
	}
}

// stopping the application and closing the database connection
func (a *App) Stop() {
	a.log.Info("application shutdown")
	a.gRPCServer.GracefulStop()

	if err := a.db.Close(); err != nil {
		a.log.Error("failed to close database connection", "error", err)
	} else {
		a.log.Info("database connection closed successfully")
	}
}

// if there is no data in the table, then load them from API and if it fails to load them,
// then load default ones, if it fails here too, then panic
func (a *App) MustRatesInit() {
	exsist, err := a.db.IsTableEmpty(storages.TableNameForCurrencyRates)
	if err != nil {
		panic(err.Error())
	}

	if !exsist {
		a.log.Info("loading of exchange rates is not required, the data already exists in the database")
		return
	}

	a.log.Info("there is no exchange rate data, let's start downloading them")

	if err = a.db.RatesDownloadFromExternalAPI(); err != nil {
		a.log.Warn("failed to download currency rates from third-party API", "error", err)

		if err = a.db.LoadDefaultRates(); err != nil {
			a.log.Error("failed to load default exchange rates", "error", err)
			panic(err.Error())
		}
		a.log.Warn("attention! currency exchange rates were loaded by DEFAULT ")
		return
	}

	a.log.Info("attention! currency rates were loaded from a third-party api")
}
