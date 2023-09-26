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
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/storage_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	storageServiceGrpcCommandName = "storage.service.grpc"
)

type IStorageServiceGrpc interface {
	// UploadImageGRPCStream(stream pb.StorageService_UploadImageGRPCStreamServer) (err error)
	UploadImageGRPC(ctx context.Context, req *pb.UploadImageGRPCStreamRequest) (res *pb.UploadImageGRPCStreamResponse, err error)
}

type StorageServiceGrpcOption struct {
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

type storageServiceGrpc struct {
	Option        StorageServiceGrpcOption
	GrpcClient    pb.StorageServiceClient
	HystrixClient *cirbreax.Client
}

func NewStorageServiceGrpc(opt StorageServiceGrpcOption) (iStorageService IStorageServiceGrpc, err error) {
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
		cirbreax.WithCommandName(storageServiceGrpcCommandName),
		cirbreax.WithHystrixTimeout(opt.Timeout),
		cirbreax.WithMaxConcurrentRequests(opt.MaxConcurrentRequests),
		cirbreax.WithErrorPercentThreshold(opt.ErrorPercentThreshold),
		cirbreax.WithRetryCount(serviceGrpcHTTPRetryCount),
		cirbreax.WithRetrier(retrier),
	)

	gRPCClient := pb.NewStorageServiceClient(conn)

	iStorageService = storageServiceGrpc{
		Option:        opt,
		GrpcClient:    gRPCClient,
		HystrixClient: client,
	}
	return
}

// func (o storageServiceGrpc) UploadImageGRPCStream(stream pb.StorageService_UploadImageGRPCStreamServer) (err error) {
// 	err = o.HystrixClient.Execute(func() (err error) {
// 		_, err = o.GrpcClient.UploadImageGRPCStream(context.TODO(), stream)
// 		return
// 	})
// 	return
// }

func (o storageServiceGrpc) UploadImageGRPC(ctx context.Context, req *pb.UploadImageGRPCStreamRequest) (res *pb.UploadImageGRPCStreamResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UploadImageGRPC(context.TODO(), req)
		return
	})
	return
}
