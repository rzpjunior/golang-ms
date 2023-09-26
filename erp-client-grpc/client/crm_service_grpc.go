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
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/crm_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	crmServiceGrpcCommandName = "crm.service.grpc"
)

type ICrmServiceGrpc interface {
	GetSalesAssignmentList(ctx context.Context, req *pb.GetSalesAssignmentListRequest) (res *pb.GetSalesAssignmentListResponse, err error)
	GetSalesAssignmentDetail(ctx context.Context, req *pb.GetSalesAssignmentDetailRequest) (res *pb.GetSalesAssignmentDetailResponse, err error)
	GetSalesAssignmentItemList(ctx context.Context, req *pb.GetSalesAssignmentItemListRequest) (res *pb.GetSalesAssignmentItemListResponse, err error)
	GetSalesAssignmentItemDetail(ctx context.Context, req *pb.GetSalesAssignmentItemDetailRequest) (res *pb.GetSalesAssignmentItemDetailResponse, err error)
	GetSalesAssignmentObjectiveList(ctx context.Context, req *pb.GetSalesAssignmentObjectiveListRequest) (res *pb.GetSalesAssignmentObjectiveListResponse, err error)
	GetSalesAssignmentObjectiveDetail(ctx context.Context, req *pb.GetSalesAssignmentObjectiveDetailRequest) (res *pb.GetSalesAssignmentObjectiveDetailResponse, err error)
	CheckTaskCustomerAcquisitionActive(ctx context.Context, req *pb.CheckTaskCustomerAcquisitionRequest) (res *pb.CheckTaskCustomerAcquisitionResponse, err error)
	CheckTaskSalesAssignmentItemActive(ctx context.Context, req *pb.CheckTaskSalesAssignmentItemRequest) (res *pb.CheckTaskSalesAssignmentItemResponse, err error)
	UpdateSubmitTaskVisitFU(ctx context.Context, req *pb.UpdateSubmitTaskVisitFURequest) (res *pb.UpdateSubmitTaskVisitFUResponse, err error)
	CheckoutTaskVisitFU(ctx context.Context, req *pb.CheckoutTaskVisitFURequest) (res *pb.CheckoutTaskVisitFUResponse, err error)
	BulkCheckoutTaskVisitFU(ctx context.Context, req *pb.BulkCheckoutTaskVisitFURequest) (res *pb.BulkCheckoutTaskVisitFUResponse, err error)
	SubmitTaskCustomerAcquisition(ctx context.Context, req *pb.SubmitTaskCustomerAcquisitionRequest) (res *pb.SubmitTaskCustomerAcquisitionResponse, err error)
	SubmitTaskFailed(ctx context.Context, req *pb.SubmitTaskFailedRequest) (res *pb.SubmitTaskFailedResponse, err error)
	CreateSalesAssignmentItem(ctx context.Context, req *pb.CreateSalesAssignmentItemRequest) (res *pb.GetSalesAssignmentItemDetailResponse, err error)
	GetCustomerAcquisitionById(ctx context.Context, req *pb.GetCustomerAcquisitionByIdRequest) (res *pb.GetCustomerAcquisitionDetailResponse, err error)
	GetCustomerAcquisitionList(ctx context.Context, req *pb.GetCustomerAcquisitionListRequest) (res *pb.GetCustomerAcquisitionListResponse, err error)
	GetCustomerAcquisitionListWithExcludedIds(ctx context.Context, req *pb.GetCustomerAcquisitionListWithExcludedIdsRequest) (res *pb.GetCustomerAcquisitionListResponse, err error)
	GetCountCustomerAcquisition(ctx context.Context, req *pb.GetCountCustomerAcquisitionRequest) (res *pb.GetCountCustomerAcquisitionResponse, err error)
	GetSalesSubmissionList(ctx context.Context, req *pb.GetSalesSubmissionListRequest) (res *pb.GetSalesAssignmentItemListResponse, err error)
	GetCustomerDetail(ctx context.Context, req *pb.GetCustomerDetailRequest) (res *pb.GetCustomerDetailResponse, err error)
	UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest) (res *pb.UpdateCustomerResponse, err error)
	CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (res *pb.CreateCustomerResponse, err error)
	GetProspectiveCustomerList(ctx context.Context, req *pb.GetProspectiveCustomerListRequest) (res *pb.GetProspectiveCustomerListResponse, err error)
	GetProspectiveCustomerDetail(ctx context.Context, req *pb.GetProspectiveCustomerDetailRequest) (res *pb.GetProspectiveCustomerDetailResponse, err error)
	DeleteProspectiveCustomer(ctx context.Context, req *pb.DeleteProspectiveCustomerRequest) (res *pb.DeleteProspectiveCustomerResponse, err error)
	CreateProspectiveCustomer(ctx context.Context, req *pb.CreateProspectiveCustomerRequest) (res *pb.CreateProspectiveCustomerResponse, err error)
	GetCustomerID(ctx context.Context, req *pb.GetCustomerIDRequest) (res *pb.GetCustomerIDResponse, err error) 
	UpdateFixedVa(ctx context.Context, req *pb.UpdateFixedVaRequest) (res *pb.UpdateFixedVaResponse, err error) 
}

type CrmServiceGrpcOption struct {
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

type crmServiceGrpc struct {
	Option        CrmServiceGrpcOption
	GrpcClient    pb.CrmServiceClient
	HystrixClient *cirbreax.Client
}

func NewCrmServiceGrpc(opt CrmServiceGrpcOption) (iCrmService ICrmServiceGrpc, err error) {
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
		cirbreax.WithCommandName(crmServiceGrpcCommandName),
		cirbreax.WithHystrixTimeout(opt.Timeout),
		cirbreax.WithMaxConcurrentRequests(opt.MaxConcurrentRequests),
		cirbreax.WithErrorPercentThreshold(opt.ErrorPercentThreshold),
		cirbreax.WithRetryCount(serviceGrpcHTTPRetryCount),
		cirbreax.WithRetrier(retrier),
	)

	gRPCClient := pb.NewCrmServiceClient(conn)

	iCrmService = crmServiceGrpc{
		Option:        opt,
		GrpcClient:    gRPCClient,
		HystrixClient: client,
	}
	return
}

func (o crmServiceGrpc) GetSalesAssignmentList(ctx context.Context, req *pb.GetSalesAssignmentListRequest) (res *pb.GetSalesAssignmentListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesAssignmentList(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) GetSalesAssignmentDetail(ctx context.Context, req *pb.GetSalesAssignmentDetailRequest) (res *pb.GetSalesAssignmentDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesAssignmentDetail(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) GetSalesAssignmentItemList(ctx context.Context, req *pb.GetSalesAssignmentItemListRequest) (res *pb.GetSalesAssignmentItemListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesAssignmentItemList(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) GetSalesAssignmentItemDetail(ctx context.Context, req *pb.GetSalesAssignmentItemDetailRequest) (res *pb.GetSalesAssignmentItemDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesAssignmentItemDetail(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) GetSalesAssignmentObjectiveList(ctx context.Context, req *pb.GetSalesAssignmentObjectiveListRequest) (res *pb.GetSalesAssignmentObjectiveListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesAssignmentObjectiveList(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) GetSalesAssignmentObjectiveDetail(ctx context.Context, req *pb.GetSalesAssignmentObjectiveDetailRequest) (res *pb.GetSalesAssignmentObjectiveDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesAssignmentObjectiveDetail(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) CheckTaskCustomerAcquisitionActive(ctx context.Context, req *pb.CheckTaskCustomerAcquisitionRequest) (res *pb.CheckTaskCustomerAcquisitionResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CheckTaskCustomerAcquisitionActive(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) CheckTaskSalesAssignmentItemActive(ctx context.Context, req *pb.CheckTaskSalesAssignmentItemRequest) (res *pb.CheckTaskSalesAssignmentItemResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CheckTaskSalesAssignmentItemActive(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) UpdateSubmitTaskVisitFU(ctx context.Context, req *pb.UpdateSubmitTaskVisitFURequest) (res *pb.UpdateSubmitTaskVisitFUResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SubmitTaskVisitFU(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) CheckoutTaskVisitFU(ctx context.Context, req *pb.CheckoutTaskVisitFURequest) (res *pb.CheckoutTaskVisitFUResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CheckoutTaskVisitFU(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) BulkCheckoutTaskVisitFU(ctx context.Context, req *pb.BulkCheckoutTaskVisitFURequest) (res *pb.BulkCheckoutTaskVisitFUResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.BulkCheckoutTaskVisitFU(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) SubmitTaskCustomerAcquisition(ctx context.Context, req *pb.SubmitTaskCustomerAcquisitionRequest) (res *pb.SubmitTaskCustomerAcquisitionResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SubmitTaskCustomerAcquisition(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) SubmitTaskFailed(ctx context.Context, req *pb.SubmitTaskFailedRequest) (res *pb.SubmitTaskFailedResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SubmitTaskFailed(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) CreateSalesAssignmentItem(ctx context.Context, req *pb.CreateSalesAssignmentItemRequest) (res *pb.GetSalesAssignmentItemDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateSalesAssignmentItem(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) GetCustomerAcquisitionById(ctx context.Context, req *pb.GetCustomerAcquisitionByIdRequest) (res *pb.GetCustomerAcquisitionDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerAcquisitionById(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) GetCustomerAcquisitionList(ctx context.Context, req *pb.GetCustomerAcquisitionListRequest) (res *pb.GetCustomerAcquisitionListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerAcquisitionList(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) GetCustomerAcquisitionListWithExcludedIds(ctx context.Context, req *pb.GetCustomerAcquisitionListWithExcludedIdsRequest) (res *pb.GetCustomerAcquisitionListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerAcquisitionListWithExcludedIds(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) GetCountCustomerAcquisition(ctx context.Context, req *pb.GetCountCustomerAcquisitionRequest) (res *pb.GetCountCustomerAcquisitionResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCountCustomerAcquisition(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) GetSalesSubmissionList(ctx context.Context, req *pb.GetSalesSubmissionListRequest) (res *pb.GetSalesAssignmentItemListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesSubmissionList(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) GetCustomerDetail(ctx context.Context, req *pb.GetCustomerDetailRequest) (res *pb.GetCustomerDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerDetail(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest) (res *pb.UpdateCustomerResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateCustomer(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (res *pb.CreateCustomerResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateCustomer(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) GetProspectiveCustomerList(ctx context.Context, req *pb.GetProspectiveCustomerListRequest) (res *pb.GetProspectiveCustomerListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetProspectiveCustomerList(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) GetProspectiveCustomerDetail(ctx context.Context, req *pb.GetProspectiveCustomerDetailRequest) (res *pb.GetProspectiveCustomerDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetProspectiveCustomerDetail(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) DeleteProspectiveCustomer(ctx context.Context, req *pb.DeleteProspectiveCustomerRequest) (res *pb.DeleteProspectiveCustomerResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.DeleteProspectiveCustomer(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) CreateProspectiveCustomer(ctx context.Context, req *pb.CreateProspectiveCustomerRequest) (res *pb.CreateProspectiveCustomerResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateProspectiveCustomer(context.TODO(), req)
		return
	})
	return
}
func (o crmServiceGrpc) GetCustomerID(ctx context.Context, req *pb.GetCustomerIDRequest) (res *pb.GetCustomerIDResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerID(context.TODO(), req)
		return
	})
	return
}

func (o crmServiceGrpc) UpdateFixedVa(ctx context.Context, req *pb.UpdateFixedVaRequest) (res *pb.UpdateFixedVaResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateFixedVa(context.TODO(), req)
		return
	})
	return
}