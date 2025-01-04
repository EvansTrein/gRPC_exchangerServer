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

	server.RegisterServ(gRPC, db)

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
	a.gRPCServer.GracefulStop()
}
