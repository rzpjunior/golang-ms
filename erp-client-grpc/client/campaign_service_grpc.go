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
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	campaignServiceGrpcCommandName = "campaign.service.grpc"
)

type ICampaignServiceGrpc interface {
	GetBannerList(ctx context.Context, req *pb.GetBannerListRequest) (res *pb.GetBannerListResponse, err error)
	GetBannerDetail(ctx context.Context, req *pb.GetBannerDetailRequest) (res *pb.GetBannerDetailResponse, err error)
	GetItemSectionList(ctx context.Context, req *pb.GetItemSectionListRequest) (res *pb.GetItemSectionListResponse, err error)
	GetItemSectionDetail(ctx context.Context, req *pb.GetItemSectionDetailRequest) (res *pb.GetItemSectionDetailResponse, err error)
	GetMembershipLevelList(ctx context.Context, req *pb.GetMembershipLevelListRequest) (res *pb.GetMembershipLevelListResponse, err error)
	GetMembershipLevelDetail(ctx context.Context, req *pb.GetMembershipLevelDetailRequest) (res *pb.GetMembershipLevelDetailResponse, err error)
	GetMembershipCheckpointList(ctx context.Context, req *pb.GetMembershipCheckpointListRequest) (res *pb.GetMembershipCheckpointListResponse, err error)
	GetMembershipCheckpointDetail(ctx context.Context, req *pb.GetMembershipCheckpointDetailRequest) (res *pb.GetMembershipCheckpointDetailResponse, err error)
	GetMembershipAdvantageDetail(ctx context.Context, req *pb.GetMembershipAdvantageDetailRequest) (res *pb.GetMembershipAdvantageDetailResponse, err error)
	GetMembershipLevelAdvantageList(ctx context.Context, req *pb.GetMembershipLevelAdvantageListRequest) (res *pb.GetMembershipLevelAdvantageListResponse, err error)
	GetMembershipRewardList(ctx context.Context, req *pb.GetMembershipRewardListRequest) (res *pb.GetMembershipRewardListResponse, err error)
	GetMembershipRewardDetail(ctx context.Context, req *pb.GetMembershipRewardDetailRequest) (res *pb.GetMembershipRewardDetailResponse, err error)
	GetCustomerPointLogList(ctx context.Context, req *pb.GetCustomerPointLogListRequest) (res *pb.GetCustomerPointLogListResponse, err error)
	GetCustomerPointLogDetail(ctx context.Context, req *pb.GetCustomerPointLogDetailRequest) (res *pb.GetCustomerPointLogDetailResponse, err error)
	GetCustomerPointLogDetailHistoryMobile(ctx context.Context, req *pb.GetCustomerPointLogDetailRequest) (res *pb.GetCustomerPointLogDetailResponse, err error)
	CreateCustomerPointLog(ctx context.Context, req *pb.CreateCustomerPointLogRequest) (res *pb.CreateCustomerPointLogResponse, err error)
	UpdateCustomerProfileTalon(ctx context.Context, req *pb.UpdateCustomerProfileTalonRequest) (res *pb.UpdateCustomerProfileTalonResponse, err error)
	UpdateCustomerSessionTalon(ctx context.Context, req *pb.UpdateCustomerSessionTalonRequest) (res *pb.UpdateCustomerSessionTalonResponse, err error)
	CreateCustomerPointSummary(ctx context.Context, req *pb.CreateCustomerPointSummaryRequest) (res *pb.CreateCustomerPointSummaryResponse, err error)
	UpdateCustomerPointSummary(ctx context.Context, req *pb.UpdateCustomerPointSummaryRequest) (res *pb.UpdateCustomerPointSummaryResponse, err error)
	GetCustomerPointSummaryDetail(ctx context.Context, req *pb.GetCustomerPointSummaryRequestDetail) (res *pb.GetCustomerPointSummaryDetailResponse, err error)
	GetCustomerMembershipDetail(ctx context.Context, req *pb.GetCustomerMembershipDetailRequest) (res *pb.GetCustomerMembershipDetailResponse, err error)
	GetPushNotificationList(ctx context.Context, req *pb.GetPushNotificationListRequest) (res *pb.GetPushNotificationListResponse, err error)
	GetPushNotificationDetail(ctx context.Context, req *pb.GetPushNotificationDetailRequest) (res *pb.GetPushNotificationDetailResponse, err error)
	GetReferralHistory(ctx context.Context, req *pb.GetReferralHistoryRequest) (res *pb.GetReferralHistoryResponse, err error)
	UpdatePushNotification(ctx context.Context, req *pb.UpdatePushNotificationRequest) (res *pb.UpdatePushNotificationResponse, err error)
	GetCustomerPointExpirationDetail(ctx context.Context, req *pb.GetCustomerPointExpirationDetailRequest) (res *pb.GetCustomerPointExpirationDetailResponse, err error)
	CancelCustomerPointLog(ctx context.Context, req *pb.CancelCustomerPointLogRequest) (res *pb.CancelCustomerPointLogResponse, err error)
}

type CampaignServiceGrpcOption struct {
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

type campaignServiceGrpc struct {
	Option        CampaignServiceGrpcOption
	GrpcClient    pb.CampaignServiceClient
	HystrixClient *cirbreax.Client
}

func NewCampaignServiceGrpc(opt CampaignServiceGrpcOption) (iCampaignService ICampaignServiceGrpc, err error) {
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
		cirbreax.WithCommandName(campaignServiceGrpcCommandName),
		cirbreax.WithHystrixTimeout(opt.Timeout),
		cirbreax.WithMaxConcurrentRequests(opt.MaxConcurrentRequests),
		cirbreax.WithErrorPercentThreshold(opt.ErrorPercentThreshold),
		cirbreax.WithRetryCount(serviceGrpcHTTPRetryCount),
		cirbreax.WithRetrier(retrier),
	)

	gRPCClient := pb.NewCampaignServiceClient(conn)

	iCampaignService = campaignServiceGrpc{
		Option:        opt,
		GrpcClient:    gRPCClient,
		HystrixClient: client,
	}
	return
}

func (o campaignServiceGrpc) GetBannerList(ctx context.Context, req *pb.GetBannerListRequest) (res *pb.GetBannerListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetBannerList(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetBannerDetail(ctx context.Context, req *pb.GetBannerDetailRequest) (res *pb.GetBannerDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetBannerDetail(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetItemSectionList(ctx context.Context, req *pb.GetItemSectionListRequest) (res *pb.GetItemSectionListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemSectionList(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetItemSectionDetail(ctx context.Context, req *pb.GetItemSectionDetailRequest) (res *pb.GetItemSectionDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemSectionDetail(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetMembershipLevelList(ctx context.Context, req *pb.GetMembershipLevelListRequest) (res *pb.GetMembershipLevelListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetMembershipLevelList(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetMembershipLevelDetail(ctx context.Context, req *pb.GetMembershipLevelDetailRequest) (res *pb.GetMembershipLevelDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetMembershipLevelDetail(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetMembershipCheckpointList(ctx context.Context, req *pb.GetMembershipCheckpointListRequest) (res *pb.GetMembershipCheckpointListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetMembershipCheckpointList(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetMembershipCheckpointDetail(ctx context.Context, req *pb.GetMembershipCheckpointDetailRequest) (res *pb.GetMembershipCheckpointDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetMembershipCheckpointDetail(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetCustomerPointLogList(ctx context.Context, req *pb.GetCustomerPointLogListRequest) (res *pb.GetCustomerPointLogListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerPointLogList(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetCustomerPointLogDetail(ctx context.Context, req *pb.GetCustomerPointLogDetailRequest) (res *pb.GetCustomerPointLogDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerPointLogDetail(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetCustomerPointLogDetailHistoryMobile(ctx context.Context, req *pb.GetCustomerPointLogDetailRequest) (res *pb.GetCustomerPointLogDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerPointLogDetailHistoryMobile(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) CreateCustomerPointLog(ctx context.Context, req *pb.CreateCustomerPointLogRequest) (res *pb.CreateCustomerPointLogResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateCustomerPointLog(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) UpdateCustomerProfileTalon(ctx context.Context, req *pb.UpdateCustomerProfileTalonRequest) (res *pb.UpdateCustomerProfileTalonResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateCustomerProfileTalon(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) UpdateCustomerSessionTalon(ctx context.Context, req *pb.UpdateCustomerSessionTalonRequest) (res *pb.UpdateCustomerSessionTalonResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateCustomerSessionTalon(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) CreateCustomerPointSummary(ctx context.Context, req *pb.CreateCustomerPointSummaryRequest) (res *pb.CreateCustomerPointSummaryResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateCustomerPointSummary(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) UpdateCustomerPointSummary(ctx context.Context, req *pb.UpdateCustomerPointSummaryRequest) (res *pb.UpdateCustomerPointSummaryResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateCustomerPointSummary(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetCustomerPointSummaryDetail(ctx context.Context, req *pb.GetCustomerPointSummaryRequestDetail) (res *pb.GetCustomerPointSummaryDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerPointSummaryDetail(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetMembershipAdvantageDetail(ctx context.Context, req *pb.GetMembershipAdvantageDetailRequest) (res *pb.GetMembershipAdvantageDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetMembershipAdvantageDetail(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetMembershipLevelAdvantageList(ctx context.Context, req *pb.GetMembershipLevelAdvantageListRequest) (res *pb.GetMembershipLevelAdvantageListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetMembershipLevelAdvantageList(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetMembershipRewardList(ctx context.Context, req *pb.GetMembershipRewardListRequest) (res *pb.GetMembershipRewardListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetMembershipRewardList(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetMembershipRewardDetail(ctx context.Context, req *pb.GetMembershipRewardDetailRequest) (res *pb.GetMembershipRewardDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetMembershipRewardDetail(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetCustomerMembershipDetail(ctx context.Context, req *pb.GetCustomerMembershipDetailRequest) (res *pb.GetCustomerMembershipDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerMembershipDetail(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetPushNotificationList(ctx context.Context, req *pb.GetPushNotificationListRequest) (res *pb.GetPushNotificationListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPushNotificationList(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetPushNotificationDetail(ctx context.Context, req *pb.GetPushNotificationDetailRequest) (res *pb.GetPushNotificationDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPushNotificationDetail(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetReferralHistory(ctx context.Context, req *pb.GetReferralHistoryRequest) (res *pb.GetReferralHistoryResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetReferralHistory(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) UpdatePushNotification(ctx context.Context, req *pb.UpdatePushNotificationRequest) (res *pb.UpdatePushNotificationResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdatePushNotification(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) GetCustomerPointExpirationDetail(ctx context.Context, req *pb.GetCustomerPointExpirationDetailRequest) (res *pb.GetCustomerPointExpirationDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerPointExpirationDetail(context.TODO(), req)
		return
	})
	return
}

func (o campaignServiceGrpc) CancelCustomerPointLog(ctx context.Context, req *pb.CancelCustomerPointLogRequest) (res *pb.CancelCustomerPointLogResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CancelCustomerPointLog(context.TODO(), req)
		return
	})
	return
}
