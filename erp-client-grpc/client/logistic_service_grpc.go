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

	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	logisticServiceGrpcCommandName = "logistic.service.grpc"
)

type ILogisticServiceGrpc interface {
	// Delivery Run Sheet
	GetDeliveryRunSheetList(ctx context.Context, req *pb.GetDeliveryRunSheetListRequest) (res *pb.GetDeliveryRunSheetListResponse, err error)
	GetDeliveryRunSheetDetail(ctx context.Context, req *pb.GetDeliveryRunSheetDetailRequest) (res *pb.GetDeliveryRunSheetDetailResponse, err error)
	CreateDeliveryRunSheet(ctx context.Context, req *pb.CreateDeliveryRunSheetRequest) (res *pb.CreateDeliveryRunSheetResponse, err error)
	FinishDeliveryRunSheet(ctx context.Context, req *pb.FinishDeliveryRunSheetRequest) (res *pb.FinishDeliveryRunSheetResponse, err error)
	// Delivery Run Sheet Item
	GetDeliveryRunSheetItemList(ctx context.Context, req *pb.GetDeliveryRunSheetItemListRequest) (res *pb.GetDeliveryRunSheetItemListResponse, err error)
	GetDeliveryRunSheetItemDetail(ctx context.Context, req *pb.GetDeliveryRunSheetItemDetailRequest) (res *pb.GetDeliveryRunSheetItemDetailResponse, err error)
	CreateDeliveryRunSheetItemPickup(ctx context.Context, req *pb.CreateDeliveryRunSheetItemRequest) (res *pb.CreateDeliveryRunSheetItemResponse, err error)
	CreateDeliveryRunSheetItemDelivery(ctx context.Context, req *pb.CreateDeliveryRunSheetItemRequest) (res *pb.CreateDeliveryRunSheetItemResponse, err error)
	StartDeliveryRunSheetItem(ctx context.Context, req *pb.StartDeliveryRunSheetItemRequest) (res *pb.StartDeliveryRunSheetItemResponse, err error)
	PostponeDeliveryRunSheetItem(ctx context.Context, req *pb.PostponeDeliveryRunSheetItemRequest) (res *pb.PostponeDeliveryRunSheetItemResponse, err error)
	FailPickupDeliveryRunSheetItem(ctx context.Context, req *pb.FailPickupDeliveryRunSheetItemRequest) (res *pb.FailPickupDeliveryRunSheetItemResponse, err error)
	FailDeliveryDeliveryRunSheetItem(ctx context.Context, req *pb.FailDeliveryDeliveryRunSheetItemRequest) (res *pb.FailDeliveryDeliveryRunSheetItemResponse, err error)
	SuccessDeliveryRunSheetItem(ctx context.Context, req *pb.SuccessDeliveryRunSheetItemRequest) (res *pb.SuccessDeliveryRunSheetItemResponse, err error)
	ArrivedDeliveryRunSheetItem(ctx context.Context, req *pb.ArrivedDeliveryRunSheetItemRequest) (res *pb.ArrivedDeliveryRunSheetItemResponse, err error)
	// Delivery Run Return
	GetDeliveryRunReturnList(ctx context.Context, req *pb.GetDeliveryRunReturnListRequest) (res *pb.GetDeliveryRunReturnListResponse, err error)
	GetDeliveryRunReturnDetail(ctx context.Context, req *pb.GetDeliveryRunReturnDetailRequest) (res *pb.GetDeliveryRunReturnDetailResponse, err error)
	CreateDeliveryRunReturn(ctx context.Context, req *pb.CreateDeliveryRunReturnRequest) (res *pb.CreateDeliveryRunReturnResponse, err error)
	UpdateDeliveryRunReturn(ctx context.Context, req *pb.UpdateDeliveryRunReturnRequest) (res *pb.UpdateDeliveryRunReturnResponse, err error)
	DeleteDeliveryRunReturn(ctx context.Context, req *pb.DeleteDeliveryRunReturnRequest) (res *pb.DeleteDeliveryRunReturnResponse, err error)
	// Delivery Run Return Item
	GetDeliveryRunReturnItemList(ctx context.Context, req *pb.GetDeliveryRunReturnItemListRequest) (res *pb.GetDeliveryRunReturnItemListResponse, err error)
	GetDeliveryRunReturnItemDetail(ctx context.Context, req *pb.GetDeliveryRunReturnItemDetailRequest) (res *pb.GetDeliveryRunReturnItemDetailResponse, err error)
	CreateDeliveryRunReturnItem(ctx context.Context, req *pb.CreateDeliveryRunReturnItemRequest) (res *pb.CreateDeliveryRunReturnItemResponse, err error)
	UpdateDeliveryRunReturnItem(ctx context.Context, req *pb.UpdateDeliveryRunReturnItemRequest) (res *pb.UpdateDeliveryRunReturnItemResponse, err error)
	DeleteDeliveryRunReturnItem(ctx context.Context, req *pb.DeleteDeliveryRunReturnItemRequest) (res *pb.DeleteDeliveryRunReturnItemResponse, err error)
	// Address Coordinate Log
	GetAddressCoordinateLogList(ctx context.Context, req *pb.GetAddressCoordinateLogListRequest) (res *pb.GetAddressCoordinateLogListResponse, err error)
	GetAddressCoordinateLogDetail(ctx context.Context, req *pb.GetAddressCoordinateLogDetailRequest) (res *pb.GetAddressCoordinateLogDetailResponse, err error)
	CreateAddressCoordinateLog(ctx context.Context, req *pb.CreateAddressCoordinateLogRequest) (res *pb.CreateAddressCoordinateLogResponse, err error)
	GetMostTrustedAddressCoordinateLog(ctx context.Context, req *pb.GetMostTrustedAddressCoordinateLogRequest) (res *pb.GetMostTrustedAddressCoordinateLogResponse, err error)
	// Courier Log
	CreateCourierLog(ctx context.Context, req *pb.CreateCourierLogRequest) (res *pb.CreateCourierLogResponse, err error)
	GetLastCourierLog(ctx context.Context, req *pb.GetLastCourierLogRequest) (res *pb.GetLastCourierLogResponse, err error)
	// Merchant Delivery Log
	CreateMerchantDeliveryLog(ctx context.Context, req *pb.CreateMerchantDeliveryLogRequest) (res *pb.CreateMerchantDeliveryLogResponse, err error)
	GetFirstMerchantDeliveryLog(ctx context.Context, req *pb.GetFirstMerchantDeliveryLogRequest) (res *pb.GetFirstMerchantDeliveryLogResponse, err error)
	// Postpone Delivery Log
	CreatePostponeDeliveryLog(ctx context.Context, req *pb.CreatePostponeDeliveryLogRequest) (res *pb.CreatePostponeDeliveryLogResponse, err error)
	// Geocode
	Geocode(ctx context.Context, req *pb.GeocodeAddressRequest) (res *pb.GeocodeAddressResponse, err error)
}

type LogisticServiceGrpcOption struct {
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

type logisticServiceGrpc struct {
	Option        LogisticServiceGrpcOption
	GrpcClient    pb.LogisticServiceClient
	HystrixClient *cirbreax.Client
}

func NewLogisticServiceGrpc(opt LogisticServiceGrpcOption) (iLogisticService ILogisticServiceGrpc, err error) {
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
		cirbreax.WithCommandName(logisticServiceGrpcCommandName),
		cirbreax.WithHystrixTimeout(opt.Timeout),
		cirbreax.WithMaxConcurrentRequests(opt.MaxConcurrentRequests),
		cirbreax.WithErrorPercentThreshold(opt.ErrorPercentThreshold),
		cirbreax.WithRetryCount(serviceGrpcHTTPRetryCount),
		cirbreax.WithRetrier(retrier),
	)

	gRPCClient := pb.NewLogisticServiceClient(conn)

	iLogisticService = logisticServiceGrpc{
		Option:        opt,
		GrpcClient:    gRPCClient,
		HystrixClient: client,
	}
	return
}

func (o logisticServiceGrpc) GetDeliveryRunSheetList(ctx context.Context, req *pb.GetDeliveryRunSheetListRequest) (res *pb.GetDeliveryRunSheetListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDeliveryRunSheetList(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) GetDeliveryRunSheetDetail(ctx context.Context, req *pb.GetDeliveryRunSheetDetailRequest) (res *pb.GetDeliveryRunSheetDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDeliveryRunSheetDetail(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) CreateDeliveryRunSheet(ctx context.Context, req *pb.CreateDeliveryRunSheetRequest) (res *pb.CreateDeliveryRunSheetResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateDeliveryRunSheet(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) FinishDeliveryRunSheet(ctx context.Context, req *pb.FinishDeliveryRunSheetRequest) (res *pb.FinishDeliveryRunSheetResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.FinishDeliveryRunSheet(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) GetDeliveryRunSheetItemList(ctx context.Context, req *pb.GetDeliveryRunSheetItemListRequest) (res *pb.GetDeliveryRunSheetItemListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDeliveryRunSheetItemList(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) GetDeliveryRunSheetItemDetail(ctx context.Context, req *pb.GetDeliveryRunSheetItemDetailRequest) (res *pb.GetDeliveryRunSheetItemDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDeliveryRunSheetItemDetail(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) CreateDeliveryRunSheetItemPickup(ctx context.Context, req *pb.CreateDeliveryRunSheetItemRequest) (res *pb.CreateDeliveryRunSheetItemResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateDeliveryRunSheetItemPickup(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) CreateDeliveryRunSheetItemDelivery(ctx context.Context, req *pb.CreateDeliveryRunSheetItemRequest) (res *pb.CreateDeliveryRunSheetItemResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateDeliveryRunSheetItemDelivery(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) StartDeliveryRunSheetItem(ctx context.Context, req *pb.StartDeliveryRunSheetItemRequest) (res *pb.StartDeliveryRunSheetItemResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.StartDeliveryRunSheetItem(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) PostponeDeliveryRunSheetItem(ctx context.Context, req *pb.PostponeDeliveryRunSheetItemRequest) (res *pb.PostponeDeliveryRunSheetItemResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.PostponeDeliveryRunSheetItem(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) FailPickupDeliveryRunSheetItem(ctx context.Context, req *pb.FailPickupDeliveryRunSheetItemRequest) (res *pb.FailPickupDeliveryRunSheetItemResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.FailPickupDeliveryRunSheetItem(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) FailDeliveryDeliveryRunSheetItem(ctx context.Context, req *pb.FailDeliveryDeliveryRunSheetItemRequest) (res *pb.FailDeliveryDeliveryRunSheetItemResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.FailDeliveryDeliveryRunSheetItem(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) SuccessDeliveryRunSheetItem(ctx context.Context, req *pb.SuccessDeliveryRunSheetItemRequest) (res *pb.SuccessDeliveryRunSheetItemResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SuccessDeliveryRunSheetItem(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) ArrivedDeliveryRunSheetItem(ctx context.Context, req *pb.ArrivedDeliveryRunSheetItemRequest) (res *pb.ArrivedDeliveryRunSheetItemResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.ArrivedDeliveryRunSheetItem(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) GetDeliveryRunReturnList(ctx context.Context, req *pb.GetDeliveryRunReturnListRequest) (res *pb.GetDeliveryRunReturnListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDeliveryRunReturnList(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) GetDeliveryRunReturnDetail(ctx context.Context, req *pb.GetDeliveryRunReturnDetailRequest) (res *pb.GetDeliveryRunReturnDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDeliveryRunReturnDetail(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) CreateDeliveryRunReturn(ctx context.Context, req *pb.CreateDeliveryRunReturnRequest) (res *pb.CreateDeliveryRunReturnResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateDeliveryRunReturn(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) UpdateDeliveryRunReturn(ctx context.Context, req *pb.UpdateDeliveryRunReturnRequest) (res *pb.UpdateDeliveryRunReturnResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateDeliveryRunReturn(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) DeleteDeliveryRunReturn(ctx context.Context, req *pb.DeleteDeliveryRunReturnRequest) (res *pb.DeleteDeliveryRunReturnResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.DeleteDeliveryRunReturn(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) GetDeliveryRunReturnItemList(ctx context.Context, req *pb.GetDeliveryRunReturnItemListRequest) (res *pb.GetDeliveryRunReturnItemListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDeliveryRunReturnItemList(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) GetDeliveryRunReturnItemDetail(ctx context.Context, req *pb.GetDeliveryRunReturnItemDetailRequest) (res *pb.GetDeliveryRunReturnItemDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDeliveryRunReturnItemDetail(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) CreateDeliveryRunReturnItem(ctx context.Context, req *pb.CreateDeliveryRunReturnItemRequest) (res *pb.CreateDeliveryRunReturnItemResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateDeliveryRunReturnItem(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) UpdateDeliveryRunReturnItem(ctx context.Context, req *pb.UpdateDeliveryRunReturnItemRequest) (res *pb.UpdateDeliveryRunReturnItemResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateDeliveryRunReturnItem(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) DeleteDeliveryRunReturnItem(ctx context.Context, req *pb.DeleteDeliveryRunReturnItemRequest) (res *pb.DeleteDeliveryRunReturnItemResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.DeleteDeliveryRunReturnItem(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) GetAddressCoordinateLogList(ctx context.Context, req *pb.GetAddressCoordinateLogListRequest) (res *pb.GetAddressCoordinateLogListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetAddressCoordinateLogList(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) GetAddressCoordinateLogDetail(ctx context.Context, req *pb.GetAddressCoordinateLogDetailRequest) (res *pb.GetAddressCoordinateLogDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetAddressCoordinateLogDetail(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) CreateAddressCoordinateLog(ctx context.Context, req *pb.CreateAddressCoordinateLogRequest) (res *pb.CreateAddressCoordinateLogResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateAddressCoordinateLog(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) GetMostTrustedAddressCoordinateLog(ctx context.Context, req *pb.GetMostTrustedAddressCoordinateLogRequest) (res *pb.GetMostTrustedAddressCoordinateLogResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetMostTrustedAddressCoordinateLog(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) CreateCourierLog(ctx context.Context, req *pb.CreateCourierLogRequest) (res *pb.CreateCourierLogResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateCourierLog(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) CreateMerchantDeliveryLog(ctx context.Context, req *pb.CreateMerchantDeliveryLogRequest) (res *pb.CreateMerchantDeliveryLogResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateMerchantDeliveryLog(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) GetFirstMerchantDeliveryLog(ctx context.Context, req *pb.GetFirstMerchantDeliveryLogRequest) (res *pb.GetFirstMerchantDeliveryLogResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetFirstMerchantDeliveryLog(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) CreatePostponeDeliveryLog(ctx context.Context, req *pb.CreatePostponeDeliveryLogRequest) (res *pb.CreatePostponeDeliveryLogResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreatePostponeDeliveryLog(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) Geocode(ctx context.Context, req *pb.GeocodeAddressRequest) (res *pb.GeocodeAddressResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.Geocode(context.TODO(), req)
		return
	})
	return
}

func (o logisticServiceGrpc) GetLastCourierLog(ctx context.Context, req *pb.GetLastCourierLogRequest) (res *pb.GetLastCourierLogResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetLastCourierLog(context.TODO(), req)
		return
	})
	return
}
