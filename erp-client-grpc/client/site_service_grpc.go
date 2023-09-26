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

	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/site_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	siteServiceGrpcCommandName = "site.service.grpc"
)

type ISiteServiceGrpc interface {
	GetKoliList(ctx context.Context, req *pb.GetKoliListRequest) (res *pb.GetKoliListResponse, err error)
	GetKoliDetail(ctx context.Context, req *pb.GetKoliDetailRequest) (res *pb.GetKoliDetailResponse, err error)
	GetSalesOrderDeliveryKoli(ctx context.Context, req *pb.GetSalesOrderDeliveryKoliRequest) (res *pb.GetSalesOrderDeliveryKoliResponse, err error)

	LoginHelper(ctx context.Context, req *pb.LoginHelperRequest) (res *pb.LoginHelperResponse, err error)
	// Picker
	GetPickingOrderHeader(ctx context.Context, req *pb.GetPickingOrderHeaderRequest) (res *pb.GetPickingOrderHeaderResponse, err error)
	GetPickingOrderDetail(ctx context.Context, req *pb.GetPickingOrderDetailRequest) (res *pb.GetPickingOrderDetailResponse, err error)
	GetAggregatedProductSalesOrder(ctx context.Context, req *pb.GetAggregatedProductSalesOrderRequest) (res *pb.GetAggregatedProductSalesOrderResponse, err error)
	StartPickingOrder(ctx context.Context, req *pb.StartPickingOrderRequest) (res *pb.SuccessResponse, err error)
	SubmitPicking(ctx context.Context, req *pb.SubmitPickingRequest) (res *pb.SuccessResponse, err error)
	GetSalesOrderPicking(ctx context.Context, req *pb.GetSalesOrderPickingRequest) (res *pb.GetSalesOrderPickingResponse, err error)
	GetSalesOrderPickingDetail(ctx context.Context, req *pb.GetSalesOrderPickingDetailRequest) (res *pb.GetSalesOrderPickingDetailResponse, err error)
	SubmitSalesOrder(ctx context.Context, req *pb.SubmitSalesOrderRequest) (res *pb.SuccessResponse, err error)
	History(ctx context.Context, req *pb.HistoryRequest) (res *pb.HistoryResponse, err error)
	HistoryDetail(ctx context.Context, req *pb.HistoryDetailRequest) (res *pb.HistoryDetailResponse, err error)
	PickerWidget(ctx context.Context, req *pb.PickerWidgetRequest) (res *pb.PickerWidgetResponse, err error)
	// SPV & Checker
	GetSalesOrderToCheck(ctx context.Context, req *pb.GetSalesOrderToCheckRequest) (res *pb.GetSalesOrderToCheckResponse, err error)
	// SPV
	SPVGetSalesOrderToCheckDetail(ctx context.Context, req *pb.GetSalesOrderToCheckDetailRequest) (res *pb.GetSalesOrderToCheckDetailResponse, err error)
	SPVRejectSalesOrder(ctx context.Context, req *pb.SPVRejectSalesOrderRequest) (res *pb.SuccessResponse, err error)
	SPVAcceptSalesOrder(ctx context.Context, req *pb.SPVAcceptSalesOrderRequest) (res *pb.SuccessResponse, err error)
	SPVWidget(ctx context.Context, req *pb.SPVWidgetRequest) (res *pb.SPVWidgetResponse, err error)
	SPVWrtMonitoring(ctx context.Context, req *pb.GetWrtMonitoringListRequest) (res *pb.GetWrtMonitoringListResponse, err error)
	SPVWrtMonitoringDetail(ctx context.Context, req *pb.GetWrtMonitoringDetailRequest) (res *pb.GetWrtMonitoringDetailResponse, err error)
	// Checker
	CheckerGetSalesOrderToCheckDetail(ctx context.Context, req *pb.GetSalesOrderToCheckDetailRequest) (res *pb.GetSalesOrderToCheckDetailResponse, err error)
	CheckerStartChecking(ctx context.Context, req *pb.CheckerStartCheckingRequest) (res *pb.SuccessResponse, err error)
	CheckerSubmitChecking(ctx context.Context, req *pb.CheckerSubmitCheckingRequest) (res *pb.SuccessResponse, err error)
	CheckerRejectSalesOrder(ctx context.Context, req *pb.CheckerRejectSalesOrderRequest) (res *pb.SuccessResponse, err error)
	CheckerGetDeliveryKoli(ctx context.Context, req *pb.CheckerGetDeliveryKoliRequest) (res *pb.CheckerGetDeliveryKoliResponse, err error)
	CheckerAcceptSalesOrder(ctx context.Context, req *pb.CheckerAcceptSalesOrderRequest) (res *pb.CheckerAcceptSalesOrderResponse, err error)
	CheckerHistory(ctx context.Context, req *pb.CheckerHistoryRequest) (res *pb.CheckerHistoryResponse, err error)
	CheckerHistoryDetail(ctx context.Context, req *pb.CheckerHistoryDetailRequest) (res *pb.CheckerHistoryDetailResponse, err error)
	CheckerWidget(ctx context.Context, req *pb.CheckerWidgetRequest) (res *pb.CheckerWidgetResponse, err error)
}

type SiteServiceGrpcOption struct {
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

type SiteServiceGrpc struct {
	Option        SiteServiceGrpcOption
	GrpcClient    pb.SiteServiceClient
	HystrixClient *cirbreax.Client
}

func NewSiteServiceGrpc(opt SiteServiceGrpcOption) (iSiteService ISiteServiceGrpc, err error) {
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
		cirbreax.WithCommandName(siteServiceGrpcCommandName),
		cirbreax.WithHystrixTimeout(opt.Timeout),
		cirbreax.WithMaxConcurrentRequests(opt.MaxConcurrentRequests),
		cirbreax.WithErrorPercentThreshold(opt.ErrorPercentThreshold),
		cirbreax.WithRetryCount(serviceGrpcHTTPRetryCount),
		cirbreax.WithRetrier(retrier),
	)

	gRPCClient := pb.NewSiteServiceClient(conn)

	iSiteService = SiteServiceGrpc{
		Option:        opt,
		GrpcClient:    gRPCClient,
		HystrixClient: client,
	}
	return
}

func (o SiteServiceGrpc) GetKoliList(ctx context.Context, req *pb.GetKoliListRequest) (res *pb.GetKoliListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetKoliList(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) GetKoliDetail(ctx context.Context, req *pb.GetKoliDetailRequest) (res *pb.GetKoliDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetKoliDetail(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) GetSalesOrderDeliveryKoli(ctx context.Context, req *pb.GetSalesOrderDeliveryKoliRequest) (res *pb.GetSalesOrderDeliveryKoliResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesOrderDeliveryKoli(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) LoginHelper(ctx context.Context, req *pb.LoginHelperRequest) (res *pb.LoginHelperResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.LoginHelper(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) GetPickingOrderHeader(ctx context.Context, req *pb.GetPickingOrderHeaderRequest) (res *pb.GetPickingOrderHeaderResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPickingOrderHeader(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) GetPickingOrderDetail(ctx context.Context, req *pb.GetPickingOrderDetailRequest) (res *pb.GetPickingOrderDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPickingOrderDetail(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) GetAggregatedProductSalesOrder(ctx context.Context, req *pb.GetAggregatedProductSalesOrderRequest) (res *pb.GetAggregatedProductSalesOrderResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetAggregatedProductSalesOrder(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) StartPickingOrder(ctx context.Context, req *pb.StartPickingOrderRequest) (res *pb.SuccessResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.StartPickingOrder(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) SubmitPicking(ctx context.Context, req *pb.SubmitPickingRequest) (res *pb.SuccessResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SubmitPicking(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) GetSalesOrderPicking(ctx context.Context, req *pb.GetSalesOrderPickingRequest) (res *pb.GetSalesOrderPickingResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesOrderPicking(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) GetSalesOrderPickingDetail(ctx context.Context, req *pb.GetSalesOrderPickingDetailRequest) (res *pb.GetSalesOrderPickingDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesOrderPickingDetail(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) SubmitSalesOrder(ctx context.Context, req *pb.SubmitSalesOrderRequest) (res *pb.SuccessResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SubmitSalesOrder(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) History(ctx context.Context, req *pb.HistoryRequest) (res *pb.HistoryResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.History(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) HistoryDetail(ctx context.Context, req *pb.HistoryDetailRequest) (res *pb.HistoryDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.HistoryDetail(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) GetSalesOrderToCheck(ctx context.Context, req *pb.GetSalesOrderToCheckRequest) (res *pb.GetSalesOrderToCheckResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesOrderToCheck(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) SPVGetSalesOrderToCheckDetail(ctx context.Context, req *pb.GetSalesOrderToCheckDetailRequest) (res *pb.GetSalesOrderToCheckDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SPVGetSalesOrderToCheckDetail(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) SPVRejectSalesOrder(ctx context.Context, req *pb.SPVRejectSalesOrderRequest) (res *pb.SuccessResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SPVRejectSalesOrder(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) SPVAcceptSalesOrder(ctx context.Context, req *pb.SPVAcceptSalesOrderRequest) (res *pb.SuccessResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SPVAcceptSalesOrder(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) CheckerGetSalesOrderToCheckDetail(ctx context.Context, req *pb.GetSalesOrderToCheckDetailRequest) (res *pb.GetSalesOrderToCheckDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CheckerGetSalesOrderToCheckDetail(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) CheckerStartChecking(ctx context.Context, req *pb.CheckerStartCheckingRequest) (res *pb.SuccessResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CheckerStartChecking(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) CheckerSubmitChecking(ctx context.Context, req *pb.CheckerSubmitCheckingRequest) (res *pb.SuccessResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CheckerSubmitChecking(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) CheckerRejectSalesOrder(ctx context.Context, req *pb.CheckerRejectSalesOrderRequest) (res *pb.SuccessResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CheckerRejectSalesOrder(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) CheckerGetDeliveryKoli(ctx context.Context, req *pb.CheckerGetDeliveryKoliRequest) (res *pb.CheckerGetDeliveryKoliResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CheckerGetDeliveryKoli(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) CheckerAcceptSalesOrder(ctx context.Context, req *pb.CheckerAcceptSalesOrderRequest) (res *pb.CheckerAcceptSalesOrderResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CheckerAcceptSalesOrder(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) CheckerHistory(ctx context.Context, req *pb.CheckerHistoryRequest) (res *pb.CheckerHistoryResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CheckerHistory(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) CheckerHistoryDetail(ctx context.Context, req *pb.CheckerHistoryDetailRequest) (res *pb.CheckerHistoryDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CheckerHistoryDetail(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) PickerWidget(ctx context.Context, req *pb.PickerWidgetRequest) (res *pb.PickerWidgetResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.PickerWidget(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) SPVWidget(ctx context.Context, req *pb.SPVWidgetRequest) (res *pb.SPVWidgetResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SPVWidget(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) CheckerWidget(ctx context.Context, req *pb.CheckerWidgetRequest) (res *pb.CheckerWidgetResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CheckerWidget(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) SPVWrtMonitoring(ctx context.Context, req *pb.GetWrtMonitoringListRequest) (res *pb.GetWrtMonitoringListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SPVWrtMonitoring(context.TODO(), req)
		return
	})
	return
}

func (o SiteServiceGrpc) SPVWrtMonitoringDetail(ctx context.Context, req *pb.GetWrtMonitoringDetailRequest) (res *pb.GetWrtMonitoringDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SPVWrtMonitoringDetail(context.TODO(), req)
		return
	})
	return
}
