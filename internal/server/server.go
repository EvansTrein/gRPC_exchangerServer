package server

import (
	"context"
	"log"

	"github.com/EvansTrein/exchanger_gRPC/internal/storages"
	pb "github.com/EvansTrein/proto-exchange/exchange"
	"google.golang.org/grpc"
)

type ServerGrpc struct {
	pb.UnimplementedExchangeServiceServer
	db storages.Database
}

func RegisterServ(gRPC *grpc.Server, db storages.Database) {
	pb.RegisterExchangeServiceServer(gRPC, &ServerGrpc{db: db})
}

func (s *ServerGrpc) GetExchangeRates(ctx context.Context, req *pb.Empty) (*pb.ExchangeRatesResponse, error) {
	log.Println("GetExchangeRates")
	resp := &pb.ExchangeRatesResponse{}

	result, err := s.db.Rate(ctx, "CNY")
	if err != nil {
		return nil, err
	}

	log.Println(result)

	answer := make(map[string]float32)
	answer["USD"] = 100
	answer["EUR"] = 200
	answer["RUB"] = 50
	answer[result.Currency] = result.Value

	resp.Rates = answer

	return resp, nil
}

func (s *ServerGrpc) GetExchangeRateForCurrency(ctx context.Context, req *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {

	return &pb.ExchangeRateResponse{}, nil
}
