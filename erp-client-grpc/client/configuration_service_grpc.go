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
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	configurationServiceGrpcCommandName = "configuration.service.grpc"
)

type IConfigurationServiceGrpc interface {
	GetGenerateCode(ctx context.Context, req *pb.GetGenerateCodeRequest) (res *pb.GetGenerateCodeResponse, err error)
	GetGenerateCustomerCode(ctx context.Context, req *pb.GetGenerateCodeRequest) (res *pb.GetGenerateCodeResponse, err error)
	GetGenerateReferralCode(ctx context.Context, req *pb.GetGenerateCodeRequest) (res *pb.GetGenerateCodeResponse, err error)
	GetGlossaryList(ctx context.Context, req *pb.GetGlossaryListRequest) (res *pb.GetGlossaryListResponse, err error)
	GetGlossaryDetail(ctx context.Context, req *pb.GetGlossaryDetailRequest) (res *pb.GetGlossaryDetailResponse, err error)
	GetConfigAppList(ctx context.Context, req *pb.GetConfigAppListRequest) (res *pb.GetConfigAppListResponse, err error)
	GetConfigAppDetail(ctx context.Context, req *pb.GetConfigAppDetailRequest) (res *pb.GetConfigAppDetailResponse, err error)
	GetWrtDetail(ctx context.Context, req *pb.GetWrtDetailRequest) (res *pb.GetWrtDetailResponse, err error)
	GetWrtIdGP(ctx context.Context, req *pb.GetWrtDetailRequest) (res *pb.GetWrtDetailResponse, err error)
	GetWrtList(ctx context.Context, req *pb.GetWrtListRequest) (res *pb.GetWrtListResponse, err error)
	GetRegionPolicyDetail(ctx context.Context, req *pb.GetRegionPolicyDetailRequest) (res *pb.GetRegionPolicyDetailResponse, err error)
	GetRegionPolicyList(ctx context.Context, req *pb.GetRegionPolicyListRequest) (res *pb.GetRegionPolicyListResponse, err error)
	GetDayOffDetail(ctx context.Context, req *pb.GetDayOffDetailRequest) (res *pb.GetDayOffDetailResponse, err error)
	GetDayOffList(ctx context.Context, req *pb.GetDayOffListRequest) (res *pb.GetDayOffListResponse, err error)
}

type ConfigurationServiceGrpcOption struct {
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

type configurationServiceGrpc struct {
	Option        ConfigurationServiceGrpcOption
	GrpcClient    pb.ConfigurationServiceClient
	HystrixClient *cirbreax.Client
}

func NewConfigurationServiceGrpc(opt ConfigurationServiceGrpcOption) (iConfigurationService IConfigurationServiceGrpc, err error) {
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
		cirbreax.WithCommandName(configurationServiceGrpcCommandName),
		cirbreax.WithHystrixTimeout(opt.Timeout),
		cirbreax.WithMaxConcurrentRequests(opt.MaxConcurrentRequests),
		cirbreax.WithErrorPercentThreshold(opt.ErrorPercentThreshold),
		cirbreax.WithRetryCount(serviceGrpcHTTPRetryCount),
		cirbreax.WithRetrier(retrier),
	)

	gRPCClient := pb.NewConfigurationServiceClient(conn)

	iConfigurationService = configurationServiceGrpc{
		Option:        opt,
		GrpcClient:    gRPCClient,
		HystrixClient: client,
	}
	return
}

func (o configurationServiceGrpc) GetGenerateCode(ctx context.Context, req *pb.GetGenerateCodeRequest) (res *pb.GetGenerateCodeResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetGenerateCode(context.TODO(), req)
		return
	})
	return
}

func (o configurationServiceGrpc) GetGenerateCustomerCode(ctx context.Context, req *pb.GetGenerateCodeRequest) (res *pb.GetGenerateCodeResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetGenerateCustomerCode(context.TODO(), req)
		return
	})
	return
}

func (o configurationServiceGrpc) GetGenerateReferralCode(ctx context.Context, req *pb.GetGenerateCodeRequest) (res *pb.GetGenerateCodeResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetGenerateReferralCode(context.TODO(), req)
		return
	})
	return
}

func (o configurationServiceGrpc) GetGlossaryList(ctx context.Context, req *pb.GetGlossaryListRequest) (res *pb.GetGlossaryListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetGlossaryList(context.TODO(), req)
		return
	})
	return
}

func (o configurationServiceGrpc) GetGlossaryDetail(ctx context.Context, req *pb.GetGlossaryDetailRequest) (res *pb.GetGlossaryDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetGlossaryDetail(context.TODO(), req)
		return
	})
	return
}

func (o configurationServiceGrpc) GetConfigAppList(ctx context.Context, req *pb.GetConfigAppListRequest) (res *pb.GetConfigAppListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetConfigAppList(context.TODO(), req)
		return
	})
	return
}

func (o configurationServiceGrpc) GetConfigAppDetail(ctx context.Context, req *pb.GetConfigAppDetailRequest) (res *pb.GetConfigAppDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetConfigAppDetail(context.TODO(), req)
		return
	})
	return
}

func (o configurationServiceGrpc) GetWrtDetail(ctx context.Context, req *pb.GetWrtDetailRequest) (res *pb.GetWrtDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetWrtDetail(context.TODO(), req)
		return
	})
	return
}

func (o configurationServiceGrpc) GetWrtIdGP(ctx context.Context, req *pb.GetWrtDetailRequest) (res *pb.GetWrtDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetWrtIdGP(context.TODO(), req)
		return
	})
	return
}

func (o configurationServiceGrpc) GetWrtList(ctx context.Context, req *pb.GetWrtListRequest) (res *pb.GetWrtListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetWrtList(context.TODO(), req)
		return
	})
	return
}

func (o configurationServiceGrpc) GetRegionPolicyList(ctx context.Context, req *pb.GetRegionPolicyListRequest) (res *pb.GetRegionPolicyListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetRegionPolicyList(context.TODO(), req)
		return
	})
	return
}
func (o configurationServiceGrpc) GetRegionPolicyDetail(ctx context.Context, req *pb.GetRegionPolicyDetailRequest) (res *pb.GetRegionPolicyDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetRegionPolicyDetail(context.TODO(), req)
		return
	})
	return
}

func (o configurationServiceGrpc) GetDayOffList(ctx context.Context, req *pb.GetDayOffListRequest) (res *pb.GetDayOffListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDayOffList(context.TODO(), req)
		return
	})
	return
}
func (o configurationServiceGrpc) GetDayOffDetail(ctx context.Context, req *pb.GetDayOffDetailRequest) (res *pb.GetDayOffDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDayOffDetail(context.TODO(), req)
		return
	})
	return
}
