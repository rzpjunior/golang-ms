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
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	catalogServiceGrpcCommandName = "catalog.service.grpc"
)

type ICatalogServiceGrpc interface {
	GetItemImageList(ctx context.Context, req *pb.GetItemImageListRequest) (res *pb.GetItemImageListResponse, err error)
	GetItemImageDetail(ctx context.Context, req *pb.GetItemImageDetailRequest) (res *pb.GetItemImageDetailResponse, err error)
	GetItemCategoryList(ctx context.Context, req *pb.GetItemCategoryListRequest) (res *pb.GetItemCategoryListResponse, err error)
	GetItemCategoryDetail(ctx context.Context, req *pb.GetItemCategoryDetailRequest) (res *pb.GetItemCategoryDetailResponse, err error)
	GetItemList(ctx context.Context, req *pb.GetItemListRequest) (res *pb.GetItemListResponse, err error)
	GetItemDetail(ctx context.Context, req *pb.GetItemDetailRequest) (res *pb.GetItemDetailResponse, err error)
	GetItemDetailByInternalId(ctx context.Context, req *pb.GetItemDetailByInternalIdRequest) (res *pb.GetItemDetailByInternalIdResponse, err error)
	GetItemDetailMasterComplexByInternalID(ctx context.Context, req *pb.GetItemDetailByInternalIdRequest) (res *pb.GetItemDetailByInternalIdResponse, err error)
	GetItemListInternal(ctx context.Context, req *pb.GetItemListRequest) (res *pb.GetItemListResponse, err error)
	GetItemDetailInternal(ctx context.Context, req *pb.GetItemDetailRequest) (res *pb.GetItemDetailResponse, err error)
	SyncMongo(ctx context.Context, req *pb.SyncMongoRequest) (res *pb.SyncMongoResponse, err error)
	GetItemListMongo(ctx context.Context, req *pb.GetItemListRequest) (res *pb.GetItemListResponse, err error)
}

type CatalogServiceGrpcOption struct {
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

type catalogServiceGrpc struct {
	Option        CatalogServiceGrpcOption
	GrpcClient    pb.CatalogServiceClient
	HystrixClient *cirbreax.Client
}

func NewCatalogServiceGrpc(opt CatalogServiceGrpcOption) (iCatalogService ICatalogServiceGrpc, err error) {
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
		cirbreax.WithCommandName(catalogServiceGrpcCommandName),
		cirbreax.WithHystrixTimeout(opt.Timeout),
		cirbreax.WithMaxConcurrentRequests(opt.MaxConcurrentRequests),
		cirbreax.WithErrorPercentThreshold(opt.ErrorPercentThreshold),
		cirbreax.WithRetryCount(serviceGrpcHTTPRetryCount),
		cirbreax.WithRetrier(retrier),
	)

	gRPCClient := pb.NewCatalogServiceClient(conn)

	iCatalogService = catalogServiceGrpc{
		Option:        opt,
		GrpcClient:    gRPCClient,
		HystrixClient: client,
	}
	return
}

func (o catalogServiceGrpc) GetItemImageList(ctx context.Context, req *pb.GetItemImageListRequest) (res *pb.GetItemImageListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemImageList(context.TODO(), req)
		return
	})
	return
}

func (o catalogServiceGrpc) GetItemImageDetail(ctx context.Context, req *pb.GetItemImageDetailRequest) (res *pb.GetItemImageDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemImageDetail(context.TODO(), req)
		return
	})
	return
}

func (o catalogServiceGrpc) GetItemCategoryList(ctx context.Context, req *pb.GetItemCategoryListRequest) (res *pb.GetItemCategoryListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemCategoryList(context.TODO(), req)
		return
	})
	return
}

func (o catalogServiceGrpc) GetItemCategoryDetail(ctx context.Context, req *pb.GetItemCategoryDetailRequest) (res *pb.GetItemCategoryDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemCategoryDetail(context.TODO(), req)
		return
	})
	return
}

func (o catalogServiceGrpc) GetItemList(ctx context.Context, req *pb.GetItemListRequest) (res *pb.GetItemListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemList(context.TODO(), req)
		return
	})
	return
}
func (o catalogServiceGrpc) GetItemListMongo(ctx context.Context, req *pb.GetItemListRequest) (res *pb.GetItemListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemListMongo(context.TODO(), req)
		return
	})
	return
}

func (o catalogServiceGrpc) GetItemDetail(ctx context.Context, req *pb.GetItemDetailRequest) (res *pb.GetItemDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemDetail(context.TODO(), req)
		return
	})
	return
}

func (o catalogServiceGrpc) GetItemDetailByInternalId(ctx context.Context, req *pb.GetItemDetailByInternalIdRequest) (res *pb.GetItemDetailByInternalIdResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemDetailByInternalId(context.TODO(), req)
		return
	})
	return
}

func (o catalogServiceGrpc) GetItemDetailMasterComplexByInternalID(ctx context.Context, req *pb.GetItemDetailByInternalIdRequest) (res *pb.GetItemDetailByInternalIdResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemDetailMasterComplexByInternalID(context.TODO(), req)
		return
	})
	return
}

func (o catalogServiceGrpc) GetItemListInternal(ctx context.Context, req *pb.GetItemListRequest) (res *pb.GetItemListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemListInternal(context.TODO(), req)
		return
	})
	return
}

func (o catalogServiceGrpc) GetItemDetailInternal(ctx context.Context, req *pb.GetItemDetailRequest) (res *pb.GetItemDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemDetailInternal(context.TODO(), req)
		return
	})
	return
}

func (o catalogServiceGrpc) SyncMongo(ctx context.Context, req *pb.SyncMongoRequest) (res *pb.SyncMongoResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SyncMongo(context.TODO(), req)
		return
	})
	return
}
