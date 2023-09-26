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

	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/notification_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	notificationServiceGrpcCommandName = "notification.service.grpc"
)

type INotificationServiceGrpc interface {
	SendNotificationTransaction(ctx context.Context, req *pb.SendNotificationTransactionRequest) (res *pb.SendNotificationTransactionResponse, err error)
	GetNotificationTransactionList(ctx context.Context, req *pb.GetNotificationTransactionListRequest) (res *pb.GetNotificationTransactionListResponse, err error)
	UpdateReadNotificationTransaction(ctx context.Context, req *pb.UpdateReadNotificationTransactionRequest) (res *pb.UpdateReadNotificationTransactionResponse, err error)
	CountUnreadNotificationTransaction(ctx context.Context, req *pb.CountUnreadNotificationTransactionRequest) (res *pb.CountUnreadNotificationTransactionResponse, err error)
	SendNotificationCampaign(ctx context.Context, req *pb.SendNotificationCampaignRequest) (res *pb.SendNotificationCampaignResponse, err error)
	GetNotificationCampaignList(ctx context.Context, req *pb.GetNotificationCampaignListRequest) (res *pb.GetNotificationCampaignListResponse, err error)
	UpdateReadNotificationCampaign(ctx context.Context, req *pb.UpdateReadNotificationCampaignRequest) (res *pb.UpdateReadNotificationCampaignResponse, err error)
	CountUnreadNotificationCampaign(ctx context.Context, req *pb.CountUnreadNotificationCampaignRequest) (res *pb.CountUnreadNotificationCampaignResponse, err error)
	SendNotificationHelper(ctx context.Context, req *pb.SendNotificationHelperRequest) (res *pb.SuccessResponse, err error)
	SendNotificationPurchaser(ctx context.Context, req *pb.SendNotificationPurchaserRequest) (res *pb.SendNotificationPurchaserResponse, err error)
	SendNotificationCancelSalesOrder(ctx context.Context, req *pb.SendNotificationCancelSalesOrderRequest) (res *pb.SendNotificationCancelSalesOrderResponse, err error)
	SendNotificationPaymentReminder(ctx context.Context, req *pb.SendNotificationPaymentReminderRequest) (res *pb.SendNotificationPaymentReminderResponse, err error)
}

type NotificationServiceGrpcOption struct {
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

type NotificationServiceGrpc struct {
	Option        NotificationServiceGrpcOption
	GrpcClient    pb.NotificationServiceClient
	HystrixClient *cirbreax.Client
}

func NewNotificationServiceGrpc(opt NotificationServiceGrpcOption) (iNotificationService INotificationServiceGrpc, err error) {
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
		cirbreax.WithCommandName(notificationServiceGrpcCommandName),
		cirbreax.WithHystrixTimeout(opt.Timeout),
		cirbreax.WithMaxConcurrentRequests(opt.MaxConcurrentRequests),
		cirbreax.WithErrorPercentThreshold(opt.ErrorPercentThreshold),
		cirbreax.WithRetryCount(serviceGrpcHTTPRetryCount),
		cirbreax.WithRetrier(retrier),
	)

	gRPCClient := pb.NewNotificationServiceClient(conn)

	iNotificationService = NotificationServiceGrpc{
		Option:        opt,
		GrpcClient:    gRPCClient,
		HystrixClient: client,
	}
	return
}

func (o NotificationServiceGrpc) SendNotificationTransaction(ctx context.Context, req *pb.SendNotificationTransactionRequest) (res *pb.SendNotificationTransactionResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SendNotificationTransaction(context.TODO(), req)
		return
	})
	return
}

func (o NotificationServiceGrpc) GetNotificationTransactionList(ctx context.Context, req *pb.GetNotificationTransactionListRequest) (res *pb.GetNotificationTransactionListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetNotificationTransactionList(context.TODO(), req)
		return
	})
	return
}

func (o NotificationServiceGrpc) UpdateReadNotificationTransaction(ctx context.Context, req *pb.UpdateReadNotificationTransactionRequest) (res *pb.UpdateReadNotificationTransactionResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateReadNotificationTransaction(context.TODO(), req)
		return
	})
	return
}

func (o NotificationServiceGrpc) CountUnreadNotificationTransaction(ctx context.Context, req *pb.CountUnreadNotificationTransactionRequest) (res *pb.CountUnreadNotificationTransactionResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CountUnreadNotificationTransaction(context.TODO(), req)
		return
	})
	return
}

func (o NotificationServiceGrpc) SendNotificationCampaign(ctx context.Context, req *pb.SendNotificationCampaignRequest) (res *pb.SendNotificationCampaignResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SendNotificationCampaign(context.TODO(), req)
		return
	})
	return
}

func (o NotificationServiceGrpc) GetNotificationCampaignList(ctx context.Context, req *pb.GetNotificationCampaignListRequest) (res *pb.GetNotificationCampaignListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetNotificationCampaignList(context.TODO(), req)
		return
	})
	return
}

func (o NotificationServiceGrpc) UpdateReadNotificationCampaign(ctx context.Context, req *pb.UpdateReadNotificationCampaignRequest) (res *pb.UpdateReadNotificationCampaignResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateReadNotificationCampaign(context.TODO(), req)
		return
	})
	return
}

func (o NotificationServiceGrpc) CountUnreadNotificationCampaign(ctx context.Context, req *pb.CountUnreadNotificationCampaignRequest) (res *pb.CountUnreadNotificationCampaignResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CountUnreadNotificationCampaign(context.TODO(), req)
		return
	})
	return
}

func (o NotificationServiceGrpc) SendNotificationHelper(ctx context.Context, req *pb.SendNotificationHelperRequest) (res *pb.SuccessResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SendNotificationHelper(context.TODO(), req)
		return
	})
	return
}

func (o NotificationServiceGrpc) SendNotificationPurchaser(ctx context.Context, req *pb.SendNotificationPurchaserRequest) (res *pb.SendNotificationPurchaserResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SendNotificationPurchaser(context.TODO(), req)
		return
	})
	return
}

func (o NotificationServiceGrpc) SendNotificationCancelSalesOrder(ctx context.Context, req *pb.SendNotificationCancelSalesOrderRequest) (res *pb.SendNotificationCancelSalesOrderResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SendNotificationCancelSalesOrder(context.TODO(), req)
		return
	})
	return
}

func (o NotificationServiceGrpc) SendNotificationPaymentReminder(ctx context.Context, req *pb.SendNotificationPaymentReminderRequest) (res *pb.SendNotificationPaymentReminderResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SendNotificationPaymentReminder(context.TODO(), req)
		return
	})
	return
}
