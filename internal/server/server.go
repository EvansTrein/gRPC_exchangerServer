package server

import (
	"context"
	"log/slog"

	"github.com/EvansTrein/exchanger_gRPC/internal/storages"
	pb "github.com/EvansTrein/proto-exchange/exchange"
	"google.golang.org/grpc"
)

type ServerGrpc struct {
	pb.UnimplementedExchangeServiceServer
	db  storages.Database
	log *slog.Logger
}

func RegisterServ(gRPC *grpc.Server, db storages.Database, log *slog.Logger) {
	pb.RegisterExchangeServiceServer(gRPC, &ServerGrpc{db: db, log: log})
}

func (s *ServerGrpc) GetExchangeRates(ctx context.Context, req *pb.Empty) (*pb.ExchangeRatesResponse, error) {
	resp := &pb.ExchangeRatesResponse{}

	return resp, nil
}

func (s *ServerGrpc) GetExchangeRateForCurrency(ctx context.Context, req *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {

	return &pb.ExchangeRateResponse{}, nil
}
