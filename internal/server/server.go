package server

import (
	"context"
	"errors"
	"log/slog"

	"github.com/EvansTrein/gRPC_exchangerServer/internal/storages"
	"github.com/EvansTrein/gRPC_exchangerServer/pkg/utils"
	pb "github.com/EvansTrein/proto-exchange/exchange"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrinValidData = errors.New("invalid data in the request")

type ServerGrpc struct {
	pb.UnimplementedExchangeServiceServer
	db  storages.Database
	log *slog.Logger
}

// our server registration
func RegisterServ(gRPC *grpc.Server, db storages.Database, log *slog.Logger) {
	pb.RegisterExchangeServiceServer(gRPC, &ServerGrpc{db: db, log: log})
}

// gRPC method to get all exchange rates
func (s *ServerGrpc) GetExchangeRates(ctx context.Context, req *pb.Empty) (*pb.ExchangeRatesResponse, error) {
	const op = "func GetExchangeRates"
	log := s.log.With(
		slog.String("operation", op),
		slog.Any("query context", ctx),
	)
	log.Debug("call of gRPC method GetExchangeRates")

	var resp pb.ExchangeRatesResponse

	result, err := s.db.AllRates(ctx)
	if err != nil {
		log.Error("failed to retrieve data from the database", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to retrieve data from the database: %v", err)
	}

	resp.Rates = result

	s.log.Info("data for all courses has been successfully submitted")
	return &resp, nil
}

// gRPC method for obtaining exchange rates
func (s *ServerGrpc) GetExchangeRateForCurrency(ctx context.Context, req *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {
	const op = "func GetExchangeRateForCurrency"
	log := s.log.With(
		slog.String("operation", op),
		slog.Any("query context", ctx),
		slog.String("FromCurrency request parameter", req.GetFromCurrency()),
		slog.String("ToCurrency request parameter", req.GetToCurrency()),
	)

	log.Debug("call of gRPC method GetExchangeRateForCurrency")

	if isValid := utils.ValidateCurrencyRequest(req.GetFromCurrency(), req.GetToCurrency()); !isValid {
		log.Error("the data in the request did not pass the validity check", "error", ErrinValidData)
		return nil, status.Error(codes.InvalidArgument, ErrinValidData.Error())
	} else {
		log.Debug("the data in the request successfully passed the validity check")
	}

	var resp pb.ExchangeRateResponse

	result, err := s.db.Rate(ctx, req.GetFromCurrency(), req.GetToCurrency())
	if err != nil {
		if err == storages.ErrExchangeRateNotFound {
			log.Error("exchange rate not found in the database", "error", err)
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Error("failed to retrieve data from the database", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to retrieve data from the database: %v", err)
	}

	// add data to the response
	resp.FromCurrency = result.BaseCurrency
	resp.ToCurrency = result.ToCurrency
	resp.Rate = result.Rate

	s.log.Info("exchange rate data successfully sent")
	return &resp, nil
}
