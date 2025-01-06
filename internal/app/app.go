package app

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/EvansTrein/exchanger_gRPC/internal/server"
	"github.com/EvansTrein/exchanger_gRPC/internal/storages"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
	db         storages.Database
}

func New(log *slog.Logger, port int, db storages.Database) *App {
	gRPC := grpc.NewServer()

	server.RegisterServ(gRPC, db, log)

	return &App{
		log:        log,
		gRPCServer: gRPC,
		port:       port,
		db: 		db,
	}
}

func (a *App) MustStart() {

	portListen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		// return fmt.Errorf("%s", err)
		panic(err.Error())
	}

	a.log.Info("grpc server started", slog.String("port", portListen.Addr().String()))
	if err := a.gRPCServer.Serve(portListen); err != nil {
		// return fmt.Errorf("%s", err)
		panic(err.Error())
	}

	// return nil
}

func (a *App) Stop() {
	a.log.Info("application shutdown")
	a.gRPCServer.GracefulStop()

	if err := a.db.Close(); err != nil {
        a.log.Error("failed to close database connection", slog.String("error", err.Error()))
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

	if err = a.db.RatesDownloadFromExternalAPI(storages.TableNameForCurrencyRates); err != nil {
		a.log.Warn("failed to download currency rates from third-party API", slog.String("error", err.Error()))
		
		if err = a.db.LoadDefaultRates(storages.TableNameForCurrencyRates); err != nil {
			panic(err.Error())
		}
		a.log.Warn("attention! currency exchange rates were loaded by default ")
		return
	}

	a.log.Info("currency rates were loaded from a third-party api")
}