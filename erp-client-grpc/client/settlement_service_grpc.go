package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/env"
	"git.edenfarm.id/project-version3/erp-pkg/erp-client-grpc/cirbreax"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/settlement_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	settlementServiceGrpcCommandName = "settlement.service.grpc"
)

type ISettlementServiceGrpc interface {
	GetSalesInvoiceExternalXendit(ctx context.Context, req *pb.GetSalesInvoiceExternalRequest) (res *pb.GetSalesInvoiceExternalResponse, err error)
	CreateSalesInvoiceExternal(ctx context.Context, req *pb.CreateSalesInvoiceExternalRequest) (res *pb.CreateSalesInvoiceExternalResponse, err error)
	GenerateFixedVaXendit(ctx context.Context, req *pb.GenerateFixedVaXenditRequest) (res *pb.GenerateFixedVaXenditResponse, err error)
}

type SettlementServiceGrpcOption struct {
	Host                  string
	Port                  int
	Timeout               time.Duration
	MaxConcurrentRequests int
	ErrorPercentThreshold int
	Tls                   bool
	PemPath               string
	Secret                string
	Realtime              bool
}

type settlementServiceGrpc struct {
	Option        SettlementServiceGrpcOption
	GrpcClient    pb.SettlementServiceClient
	HystrixClient *cirbreax.Client
}

func NewSettlementServiceGrpc(opt SettlementServiceGrpcOption) (iSettlementService ISettlementServiceGrpc, err error) {
	var opts []grpc.DialOption
	env, e := env.Env("env")
	if e != nil {
		return
	}
	serviceGrpcHTTPBackoffInterval := time.Duration(env.GetInt("client.serviceGrpcHTTPBackoffInterval")) * time.Millisecond
	serviceGrpcHTTPMaxJitterInterval := time.Duration(env.GetInt("client.serviceGrpcHTTPMaxJitterInterval")) * time.Millisecond
	serviceGrpcHTTPTimeout := time.Duration(env.GetInt("client.serviceGrpcHTTPTimeout")) * time.Millisecond
	serviceGrpcHTTPRetryCount := env.GetInt("client.serviceGrpcHTTPRetryCount")

	if opt.Tls {
		var pemServerCA []byte
		pemServerCA, err = ioutil.ReadFile(opt.PemPath)
		if err != nil {
			return
		}

		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(pemServerCA) {
			err = errors.New("failed to add server ca's certificate")
			return
		}

		// Create the credentials and return it
		config := &tls.Config{
			RootCAs: certPool,
		}

		tlsCredentials := credentials.NewTLS(config)

		opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	opts = append(opts, grpc.WithReturnConnectionError())
	opts = append(opts, grpc.FailOnNonTempDialError(true))

	var conn *grpc.ClientConn
	conn, err = grpc.Dial(fmt.Sprintf("%s:%d", opt.Host, opt.Port),
		opts...,
	)
	if err != nil {
		return
	}
	env.GetString("client.serviceGrpcHTTPTimeout")
	backoff := cirbreax.NewConstantBackoff(serviceGrpcHTTPBackoffInterval, serviceGrpcHTTPMaxJitterInterval)
	retrier := cirbreax.NewRetrier(backoff)

	client := cirbreax.NewHttpClient(
		cirbreax.WithHTTPTimeout(serviceGrpcHTTPTimeout),
		cirbreax.WithCommandName(settlementServiceGrpcCommandName),
		cirbreax.WithHystrixTimeout(opt.Timeout),
		cirbreax.WithMaxConcurrentRequests(opt.MaxConcurrentRequests),
		cirbreax.WithErrorPercentThreshold(opt.ErrorPercentThreshold),
		cirbreax.WithRetryCount(serviceGrpcHTTPRetryCount),
		cirbreax.WithRetrier(retrier),
	)

	gRPCClient := pb.NewSettlementServiceClient(conn)

	iSettlementService = settlementServiceGrpc{
		Option:        opt,
		GrpcClient:    gRPCClient,
		HystrixClient: client,
	}
	return
}

func (o settlementServiceGrpc) CreateSalesInvoiceExternal(ctx context.Context, req *pb.CreateSalesInvoiceExternalRequest) (res *pb.CreateSalesInvoiceExternalResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateSalesInvoiceExternal(context.TODO(), req)
		return
	})
	return
}
func (o settlementServiceGrpc) GetSalesInvoiceExternalXendit(ctx context.Context, req *pb.GetSalesInvoiceExternalRequest) (res *pb.GetSalesInvoiceExternalResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesInvoiceExternalXendit(context.TODO(), req)
		return
	})
	return
}

func (o settlementServiceGrpc) GenerateFixedVaXendit(ctx context.Context, req *pb.GenerateFixedVaXenditRequest) (res *pb.GenerateFixedVaXenditResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GenerateFixedVaXendit(context.TODO(), req)
		return
	})
	return
}
