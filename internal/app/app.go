package app

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/EvansTrein/gRPC_exchangerServer/internal/server"
	"github.com/EvansTrein/gRPC_exchangerServer/internal/storages"
	"google.golang.org/grpc"
)

type App struct {
	log               *slog.Logger
	gRPCServer        *grpc.Server
	port              int
	db                storages.Database
	connectionTimeout time.Duration
}

func New(log *slog.Logger, port int, db storages.Database, connectionTimeout time.Duration) *App {
	gRPC := grpc.NewServer(grpc.ConnectionTimeout(connectionTimeout))

	server.RegisterServ(gRPC, db, log)

	return &App{
		log:               log,
		gRPCServer:        gRPC,
		port:              port,
		db:                db,
		connectionTimeout: connectionTimeout,
	}
}

func (a *App) MustStart() {

	portListen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		panic(err.Error())
	}

	a.log.Info("grpc server started", slog.String("port", portListen.Addr().String()))
	if err := a.gRPCServer.Serve(portListen); err != nil {
		panic(err.Error())
	}
}

func (a *App) Stop() {
	a.log.Info("application shutdown")
	a.gRPCServer.GracefulStop()

	if err := a.db.Close(); err != nil {
		a.log.Error("failed to close database connection", "error", err)
	} else {
		a.log.Info("database connection closed successfully")
	}
}

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
