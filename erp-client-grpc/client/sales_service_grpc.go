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
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	salesServiceGrpcCommandName = "sales.service.grpc"
)

type ISalesServiceGrpc interface {
	GetSalesOrderList(ctx context.Context, req *pb.GetSalesOrderListRequest) (res *pb.GetSalesOrderListResponse, err error)
	GetSalesOrderDetail(ctx context.Context, req *pb.GetSalesOrderDetailRequest) (res *pb.GetSalesOrderDetailResponse, err error)
	GetSalesOrderItemList(ctx context.Context, req *pb.GetSalesOrderItemListRequest) (res *pb.GetSalesOrderItemListResponse, err error)
	GetSalesOrderItemDetail(ctx context.Context, req *pb.GetSalesOrderItemDetailRequest) (res *pb.GetSalesOrderItemDetailResponse, err error)
	CreateSalesOrder(ctx context.Context, req *pb.CreateSalesOrderRequest) (res *pb.CreateSalesOrderResponse, err error)
	UpdateSalesOrder(ctx context.Context, req *pb.UpdateSalesOrderRequest) (res *pb.UpdateSalesOrderResponse, err error)
	GetSalesOrderListMobile(ctx context.Context, req *pb.GetSalesOrderListRequest) (res *pb.GetSalesOrderListResponse, err error)
	GetSalesOrderFeedbackList(ctx context.Context, req *pb.GetSalesOrderFeedbackListRequest) (res *pb.GetSalesOrderFeedbackListResponse, err error)
	// GetSalesOrderFeedbackDetail(ctx context.Context, req *pb.GetSalesOrderFeedbackDetailRequest) (res *pb.GetSalesOrderFeedbackDetailResponse, err error)
	CreateSalesOrderFeedback(ctx context.Context, req *pb.CreateSalesOrderFeedbackRequest) (res *pb.CreateSalesOrderFeedbackResponse, err error)
	GetPaymentChannelList(ctx context.Context, req *pb.GetPaymentChannelListRequest) (res *pb.GetPaymentChannelListResponse, err error)
	GetPaymentMethodList(ctx context.Context, req *pb.GetPaymentMethodListRequest) (res *pb.GetPaymentMethodListResponse, err error)
	GetPaymentGroupCombList(ctx context.Context, req *pb.GetPaymentGroupCombListRequest) (res *pb.GetPaymentGroupCombListResponse, err error)
	GetSalesInvoiceGPMobileList(ctx context.Context, req *pb.GetSalesInvoiceGPMobileListRequest) (res *pb.GetSalesInvoiceGPMobileListResponse, err error)
	GetSalesOrderListCronJob(ctx context.Context, req *pb.GetSalesOrderListCronjobRequest) (res *pb.GetSalesOrderListCronjobResponse, err error)
	UpdateSalesOrderRemindPayment(ctx context.Context, req *pb.UpdateSalesOrderRemindPaymentRequest) (res *pb.UpdateSalesOrderRemindPaymentResponse, err error)
	ExpiredSalesOrder(ctx context.Context, req *pb.ExpiredSalesOrderRequest) (res *pb.ExpiredSalesOrderResponse, err error)
	CreateSalesOrderPaid(ctx context.Context, req *pb.CreateSalesOrderPaidRequest) (res *pb.CreateSalesOrderPaidResponse, err error)
	GetDeltaPrintSiEdnDetail(ctx context.Context, req *pb.GetDeltaPrintSiEdnDetailRequest) (res *pb.GetDeltaPrintSiEdnDetailResponse, err error)
	GetDeltaPrintSpEdnDetail(ctx context.Context, req *pb.GetDeltaPrintSpEdnDetailRequest) (res *pb.GetDeltaPrintSpEdnDetailResponse, err error)
	CreateDeltaPrintSiEdn(ctx context.Context, req *pb.CreateDeltaPrintSiEdnRequest) (res *pb.CreateDeltaPrintSiEdnResponse, err error)
	UpdateDeltaPrintSiEdn(ctx context.Context, req *pb.UpdateDeltaPrintSiEdnRequest) (res *pb.UpdateDeltaPrintSiEdnResponse, err error)
	CreateDeltaPrintSpEdn(ctx context.Context, req *pb.CreateDeltaPrintSpEdnRequest) (res *pb.CreateDeltaPrintSpEdnResponse, err error)
	UpdateDeltaPrintSpEdn(ctx context.Context, req *pb.UpdateDeltaPrintSpEdnRequest) (res *pb.UpdateDeltaPrintSpEdnResponse, err error)
}

type SalesServiceGrpcOption struct {
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

type salesServiceGrpc struct {
	Option        SalesServiceGrpcOption
	GrpcClient    pb.SalesServiceClient
	HystrixClient *cirbreax.Client
}

func NewSalesServiceGrpc(opt SalesServiceGrpcOption) (iSalesService ISalesServiceGrpc, err error) {
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
		cirbreax.WithCommandName(salesServiceGrpcCommandName),
		cirbreax.WithHystrixTimeout(opt.Timeout),
		cirbreax.WithMaxConcurrentRequests(opt.MaxConcurrentRequests),
		cirbreax.WithErrorPercentThreshold(opt.ErrorPercentThreshold),
		cirbreax.WithRetryCount(serviceGrpcHTTPRetryCount),
		cirbreax.WithRetrier(retrier),
	)

	gRPCClient := pb.NewSalesServiceClient(conn)

	iSalesService = salesServiceGrpc{
		Option:        opt,
		GrpcClient:    gRPCClient,
		HystrixClient: client,
	}
	return
}

func (o salesServiceGrpc) GetSalesOrderList(ctx context.Context, req *pb.GetSalesOrderListRequest) (res *pb.GetSalesOrderListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesOrderList(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) GetSalesOrderDetail(ctx context.Context, req *pb.GetSalesOrderDetailRequest) (res *pb.GetSalesOrderDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesOrderDetail(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) GetSalesOrderItemList(ctx context.Context, req *pb.GetSalesOrderItemListRequest) (res *pb.GetSalesOrderItemListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesOrderItemList(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) GetSalesOrderItemDetail(ctx context.Context, req *pb.GetSalesOrderItemDetailRequest) (res *pb.GetSalesOrderItemDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesOrderItemDetail(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) CreateSalesOrder(ctx context.Context, req *pb.CreateSalesOrderRequest) (res *pb.CreateSalesOrderResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateSalesOrder(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) UpdateSalesOrder(ctx context.Context, req *pb.UpdateSalesOrderRequest) (res *pb.UpdateSalesOrderResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateSalesOrder(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) GetSalesOrderListMobile(ctx context.Context, req *pb.GetSalesOrderListRequest) (res *pb.GetSalesOrderListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesOrderListMobile(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) GetSalesOrderFeedbackList(ctx context.Context, req *pb.GetSalesOrderFeedbackListRequest) (res *pb.GetSalesOrderFeedbackListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesOrderFeedbackList(context.TODO(), req)
		return
	})
	return
}

// func (o salesServiceGrpc) GetSalesOrderFeedbackDetail(ctx context.Context, req *pb.GetSalesOrderFeedbackDetailRequest) (res *pb.GetSalesOrderFeedbackDetailResponse, err error) {
// 	err = o.HystrixClient.Execute(func() (err error) {
// 		res, err = o.GrpcClient.GetSalesOrderFeedbackDetail(context.TODO(), req)
// 		return
// 	})
// 	return
// }

func (o salesServiceGrpc) CreateSalesOrderFeedback(ctx context.Context, req *pb.CreateSalesOrderFeedbackRequest) (res *pb.CreateSalesOrderFeedbackResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateSalesOrderFeedback(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) GetPaymentChannelList(ctx context.Context, req *pb.GetPaymentChannelListRequest) (res *pb.GetPaymentChannelListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPaymentChannelList(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) GetPaymentMethodList(ctx context.Context, req *pb.GetPaymentMethodListRequest) (res *pb.GetPaymentMethodListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPaymentMethodList(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) GetPaymentGroupCombList(ctx context.Context, req *pb.GetPaymentGroupCombListRequest) (res *pb.GetPaymentGroupCombListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPaymentGroupCombList(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) GetSalesInvoiceGPMobileList(ctx context.Context, req *pb.GetSalesInvoiceGPMobileListRequest) (res *pb.GetSalesInvoiceGPMobileListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesInvoiceGPMobileList(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) GetSalesOrderListCronJob(ctx context.Context, req *pb.GetSalesOrderListCronjobRequest) (res *pb.GetSalesOrderListCronjobResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesOrderListCronJob(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) UpdateSalesOrderRemindPayment(ctx context.Context, req *pb.UpdateSalesOrderRemindPaymentRequest) (res *pb.UpdateSalesOrderRemindPaymentResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateSalesOrderRemindPayment(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) ExpiredSalesOrder(ctx context.Context, req *pb.ExpiredSalesOrderRequest) (res *pb.ExpiredSalesOrderResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.ExpiredSalesOrder(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) CreateSalesOrderPaid(ctx context.Context, req *pb.CreateSalesOrderPaidRequest) (res *pb.CreateSalesOrderPaidResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateSalesOrderPaid(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) GetDeltaPrintSiEdnDetail(ctx context.Context, req *pb.GetDeltaPrintSiEdnDetailRequest) (res *pb.GetDeltaPrintSiEdnDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDeltaPrintSiEdnDetail(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) GetDeltaPrintSpEdnDetail(ctx context.Context, req *pb.GetDeltaPrintSpEdnDetailRequest) (res *pb.GetDeltaPrintSpEdnDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDeltaPrintSpEdnDetail(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) CreateDeltaPrintSiEdn(ctx context.Context, req *pb.CreateDeltaPrintSiEdnRequest) (res *pb.CreateDeltaPrintSiEdnResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateDeltaPrintSiEdn(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) UpdateDeltaPrintSiEdn(ctx context.Context, req *pb.UpdateDeltaPrintSiEdnRequest) (res *pb.UpdateDeltaPrintSiEdnResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateDeltaPrintSiEdn(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) CreateDeltaPrintSpEdn(ctx context.Context, req *pb.CreateDeltaPrintSpEdnRequest) (res *pb.CreateDeltaPrintSpEdnResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateDeltaPrintSpEdn(context.TODO(), req)
		return
	})
	return
}

func (o salesServiceGrpc) UpdateDeltaPrintSpEdn(ctx context.Context, req *pb.UpdateDeltaPrintSpEdnRequest) (res *pb.UpdateDeltaPrintSpEdnResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateDeltaPrintSpEdn(context.TODO(), req)
		return
	})
	return
}
