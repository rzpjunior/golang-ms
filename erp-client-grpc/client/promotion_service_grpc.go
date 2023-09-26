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
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/promotion_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	promotionServiceGrpcCommandName = "promotion.service.grpc"
)

type IPromotionServiceGrpc interface {
	GetVoucherMobileList(ctx context.Context, req *pb.GetVoucherMobileListRequest) (res *pb.GetVoucherMobileListResponse, err error)
	GetVoucherMobileDetail(ctx context.Context, req *pb.GetVoucherMobileDetailRequest) (res *pb.GetVoucherMobileDetailResponse, err error)
	GetVoucherItemList(ctx context.Context, req *pb.GetVoucherItemListRequest) (res *pb.GetVoucherItemListResponse, err error)
	CreateVoucher(ctx context.Context, req *pb.CreateVoucherRequest) (res *pb.CreateVoucherResponse, err error)
	GetVoucherLogList(ctx context.Context, req *pb.GetVoucherLogListRequest) (res *pb.GetVoucherLogListResponse, err error)
	CreateVoucherLog(ctx context.Context, req *pb.CreateVoucherLogRequest) (res *pb.CreateVoucherLogResponse, err error)
	CancelVoucherLog(ctx context.Context, req *pb.CancelVoucherLogRequest) (res *pb.CancelVoucherLogResponse, err error)
	UpdateVoucher(ctx context.Context, req *pb.UpdateVoucherRequest) (res *pb.UpdateVoucherResponse, err error)
	CreatePriceTieringLog(ctx context.Context, req *pb.CreatePriceTieringLogRequest) (res *pb.CreatePriceTieringLogResponse, err error)
	CancelPriceTieringLog(ctx context.Context, req *pb.CancelPriceTieringLogRequest) (res *pb.CancelPriceTieringLogResponse, err error)
	GetPriceTieringLogList(ctx context.Context, req *pb.GetPriceTieringLogListRequest) (res *pb.GetPriceTieringLogListResponse, err error)
}

type PromotionServiceGrpcOption struct {
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

type promotionServiceGrpc struct {
	Option        PromotionServiceGrpcOption
	GrpcClient    pb.PromotionServiceClient
	HystrixClient *cirbreax.Client
}

func NewPromotionServiceGrpc(opt PromotionServiceGrpcOption) (iPromotionService IPromotionServiceGrpc, err error) {
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

	backoff := cirbreax.NewConstantBackoff(serviceGrpcHTTPBackoffInterval, serviceGrpcHTTPMaxJitterInterval)
	retrier := cirbreax.NewRetrier(backoff)

	client := cirbreax.NewHttpClient(
		cirbreax.WithHTTPTimeout(serviceGrpcHTTPTimeout),
		cirbreax.WithCommandName(promotionServiceGrpcCommandName),
		cirbreax.WithHystrixTimeout(opt.Timeout),
		cirbreax.WithMaxConcurrentRequests(opt.MaxConcurrentRequests),
		cirbreax.WithErrorPercentThreshold(opt.ErrorPercentThreshold),
		cirbreax.WithRetryCount(serviceGrpcHTTPRetryCount),
		cirbreax.WithRetrier(retrier),
	)

	gRPCClient := pb.NewPromotionServiceClient(conn)

	iPromotionService = promotionServiceGrpc{
		Option:        opt,
		GrpcClient:    gRPCClient,
		HystrixClient: client,
	}
	return
}

func (o promotionServiceGrpc) GetVoucherMobileList(ctx context.Context, req *pb.GetVoucherMobileListRequest) (res *pb.GetVoucherMobileListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVoucherMobileList(context.TODO(), req)
		return
	})
	return
}

func (o promotionServiceGrpc) GetVoucherMobileDetail(ctx context.Context, req *pb.GetVoucherMobileDetailRequest) (res *pb.GetVoucherMobileDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVoucherMobileDetail(context.TODO(), req)
		return
	})
	return
}

func (o promotionServiceGrpc) GetVoucherItemList(ctx context.Context, req *pb.GetVoucherItemListRequest) (res *pb.GetVoucherItemListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVoucherItemList(context.TODO(), req)
		return
	})
	return
}

func (o promotionServiceGrpc) CreateVoucher(ctx context.Context, req *pb.CreateVoucherRequest) (res *pb.CreateVoucherResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateVoucher(context.TODO(), req)
		return
	})
	return
}

func (o promotionServiceGrpc) GetVoucherLogList(ctx context.Context, req *pb.GetVoucherLogListRequest) (res *pb.GetVoucherLogListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVoucherLogList(context.TODO(), req)
		return
	})
	return
}

func (o promotionServiceGrpc) CreateVoucherLog(ctx context.Context, req *pb.CreateVoucherLogRequest) (res *pb.CreateVoucherLogResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateVoucherLog(context.TODO(), req)
		return
	})
	return
}

func (o promotionServiceGrpc) CancelVoucherLog(ctx context.Context, req *pb.CancelVoucherLogRequest) (res *pb.CancelVoucherLogResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CancelVoucherLog(context.TODO(), req)
		return
	})
	return
}

func (o promotionServiceGrpc) UpdateVoucher(ctx context.Context, req *pb.UpdateVoucherRequest) (res *pb.UpdateVoucherResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateVoucher(context.TODO(), req)
		return
	})
	return
}

func (o promotionServiceGrpc) CreatePriceTieringLog(ctx context.Context, req *pb.CreatePriceTieringLogRequest) (res *pb.CreatePriceTieringLogResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreatePriceTieringLog(context.TODO(), req)
		return
	})
	return
}

func (o promotionServiceGrpc) CancelPriceTieringLog(ctx context.Context, req *pb.CancelPriceTieringLogRequest) (res *pb.CancelPriceTieringLogResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CancelPriceTieringLog(context.TODO(), req)
		return
	})
	return
}

func (o promotionServiceGrpc) GetPriceTieringLogList(ctx context.Context, req *pb.GetPriceTieringLogListRequest) (res *pb.GetPriceTieringLogListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPriceTieringLogList(context.TODO(), req)
		return
	})
	return
}
