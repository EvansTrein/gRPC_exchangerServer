package tests

import (
	"context"
	"testing"
	"time"

	pb "github.com/EvansTrein/proto-exchange/exchange"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AppTest struct {
	*testing.T                              
	ExchangeClient pb.ExchangeServiceClient 
}

const (
	grpcHost = "localhost"
	grpcPort = "44000"
	requestTimeout = time.Second * 10
)

func NewAppTest(t *testing.T) (context.Context, *AppTest) {
	t.Helper()

	ctx, cancelCtx := context.WithTimeout(context.Background(), requestTimeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	target := grpcHost + ":" + grpcPort
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc test server connection failed: %v", err)
	}

	return ctx, &AppTest{
		T:              t,
		ExchangeClient: pb.NewExchangeServiceClient(conn),
	}
}


