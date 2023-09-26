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
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/account_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	accountServiceGrpcCommandName = "account.service.grpc"
)

type IAccountServiceGrpc interface {
	GetUserList(ctx context.Context, req *pb.GetUserListRequest) (res *pb.GetUserListResponse, err error)
	GetUserDetail(ctx context.Context, req *pb.GetUserDetailRequest) (res *pb.GetUserDetailResponse, err error)
	GetUserEmailAuth(ctx context.Context, req *pb.GetUserEmailAuthRequest) (res *pb.GetUserEmailAuthResponse, err error)
	GetUserRolesByUserId(ctx context.Context, req *pb.GetUserRoleByUserIdRequest) (res *pb.GetUserRoleByUserIdResponse, err error)
	UpdateUserSalesAppToken(ctx context.Context, req *pb.UpdateUserSalesAppTokenRequest) (res *pb.GetUserDetailResponse, err error)
	GetUserBySalesAppLoginToken(ctx context.Context, req *pb.GetUserBySalesAppLoginTokenRequest) (res *pb.GetUserDetailResponse, err error)
	GetRoleDetail(ctx context.Context, req *pb.GetRoleDetailRequest) (res *pb.GetRoleDetailResponse, err error)
	UpdateUserEdnAppToken(ctx context.Context, req *pb.UpdateUserEdnAppTokenRequest) (res *pb.GetUserDetailResponse, err error)
	UpdateUserPurchaserAppToken(ctx context.Context, req *pb.UpdateUserPurchaserAppTokenRequest) (res *pb.GetUserDetailResponse, err error)
	GetUserByEdnAppLoginToken(ctx context.Context, req *pb.GetUserByEdnAppLoginTokenRequest) (res *pb.GetUserDetailResponse, err error)
	GetDivisionDetail(ctx context.Context, req *pb.GetDivisionDetailRequest) (res *pb.GetDivisionDetailResponse, err error)
	GetDivisionDefaultByCustomerType(ctx context.Context, req *pb.GetDivisionDefaultByCustomerTypeRequest) (res *pb.GetDivisionDefaultByCustomerTypeResponse, err error)
}

type AccountServiceGrpcOption struct {
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

type accountServiceGrpc struct {
	Option        AccountServiceGrpcOption
	GrpcClient    pb.AccountServiceClient
	HystrixClient *cirbreax.Client
}

func NewAccountServiceGrpc(opt AccountServiceGrpcOption) (iAccountService IAccountServiceGrpc, err error) {
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
		cirbreax.WithCommandName(accountServiceGrpcCommandName),
		cirbreax.WithHystrixTimeout(opt.Timeout),
		cirbreax.WithMaxConcurrentRequests(opt.MaxConcurrentRequests),
		cirbreax.WithErrorPercentThreshold(opt.ErrorPercentThreshold),
		cirbreax.WithRetryCount(serviceGrpcHTTPRetryCount),
		cirbreax.WithRetrier(retrier),
	)

	gRPCClient := pb.NewAccountServiceClient(conn)

	iAccountService = accountServiceGrpc{
		Option:        opt,
		GrpcClient:    gRPCClient,
		HystrixClient: client,
	}
	return
}

func (o accountServiceGrpc) GetUserList(ctx context.Context, req *pb.GetUserListRequest) (res *pb.GetUserListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetUserList(context.TODO(), req)
		return
	})
	return
}

func (o accountServiceGrpc) GetUserDetail(ctx context.Context, req *pb.GetUserDetailRequest) (res *pb.GetUserDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetUserDetail(context.TODO(), req)
		return
	})
	return
}

func (o accountServiceGrpc) GetUserEmailAuth(ctx context.Context, req *pb.GetUserEmailAuthRequest) (res *pb.GetUserEmailAuthResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetUserEmailAuth(context.TODO(), req)
		return
	})
	return
}

func (o accountServiceGrpc) GetUserRolesByUserId(ctx context.Context, req *pb.GetUserRoleByUserIdRequest) (res *pb.GetUserRoleByUserIdResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetUserRoleByUserId(context.TODO(), req)
		return
	})
	return
}

func (o accountServiceGrpc) UpdateUserSalesAppToken(ctx context.Context, req *pb.UpdateUserSalesAppTokenRequest) (res *pb.GetUserDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateUserSalesAppToken(context.TODO(), req)
		return
	})
	return
}

func (o accountServiceGrpc) GetUserBySalesAppLoginToken(ctx context.Context, req *pb.GetUserBySalesAppLoginTokenRequest) (res *pb.GetUserDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetUserBySalesAppLoginToken(context.TODO(), req)
		return
	})
	return
}

func (o accountServiceGrpc) GetRoleDetail(ctx context.Context, req *pb.GetRoleDetailRequest) (res *pb.GetRoleDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetRoleDetail(context.TODO(), req)
		return
	})
	return
}

func (o accountServiceGrpc) UpdateUserEdnAppToken(ctx context.Context, req *pb.UpdateUserEdnAppTokenRequest) (res *pb.GetUserDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateUserEdnAppToken(context.TODO(), req)
		return
	})
	return
}

func (o accountServiceGrpc) UpdateUserPurchaserAppToken(ctx context.Context, req *pb.UpdateUserPurchaserAppTokenRequest) (res *pb.GetUserDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateUserPurchaserAppToken(context.TODO(), req)
		return
	})
	return
}

func (o accountServiceGrpc) GetUserByEdnAppLoginToken(ctx context.Context, req *pb.GetUserByEdnAppLoginTokenRequest) (res *pb.GetUserDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetUserByEdnAppLoginToken(context.TODO(), req)
		return
	})
	return
}

func (o accountServiceGrpc) GetDivisionDetail(ctx context.Context, req *pb.GetDivisionDetailRequest) (res *pb.GetDivisionDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDivisionDetail(context.TODO(), req)
		return
	})
	return
}

func (o accountServiceGrpc) GetDivisionDefaultByCustomerType(ctx context.Context, req *pb.GetDivisionDefaultByCustomerTypeRequest) (res *pb.GetDivisionDefaultByCustomerTypeResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDivisionDefaultByCustomerType(context.TODO(), req)
		return
	})
	return
}
