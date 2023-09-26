package service

import (
	"context"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
)

type ICustomerService interface {
	GetCustomerMembership(ctx context.Context, req dto.RequestGetPostSession) (res dto.CustomerMembership, err error)
	GetReferralHistory(ctx context.Context, req dto.RequestGetPostSession) (res dto.ReferralHistoryReturn, err error)
}

type CustomerService struct {
	opt opt.Options
	//RepositoryOTPOutgoing repository.IOtpOutgoingRepository
}

func NewCustomerService() ICustomerService {
	return &CustomerService{
		opt: global.Setup.Common,
		//RepositoryOTPOutgoing: repository.NewOtpOutgoingRepository(),
	}
}

func (s *CustomerService) GetCustomerMembership(ctx context.Context, req dto.RequestGetPostSession) (res dto.CustomerMembership, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerService.Get")
	defer span.End()

	custMembership, err := s.opt.Client.CampaignServiceGrpc.GetCustomerMembershipDetail(ctx, &campaign_service.GetCustomerMembershipDetailRequest{
		MembershipLevelId: req.Session.Customer.MembershipLevel.ID,
		ProfileCode:       req.Session.Customer.ProfileCode,
	})
	if err != nil {
		return res, err
	}

	res = dto.CustomerMembership{
		MembershipLevel:      strconv.Itoa(int(custMembership.Data.MembershipLevel)),
		MembershipLevelName:  custMembership.Data.MembershipLevelName,
		MembershipCheckpoint: strconv.Itoa(int(custMembership.Data.MembershipCheckpoint)),
		CheckpointPercentage: strconv.FormatFloat(custMembership.Data.CheckpointPercentage, 'f', 1, 64),
		CurrentAmount:        strconv.FormatFloat(custMembership.Data.CurrentAmount, 'f', 1, 64),
		TargetAmount:         strconv.FormatFloat(custMembership.Data.TargetAmount, 'f', 1, 64),
	}

	return
}

func (s *CustomerService) GetReferralHistory(ctx context.Context, req dto.RequestGetPostSession) (res dto.ReferralHistoryReturn, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerService.GetReferralHistory")
	defer span.End()

	customerID, _ := strconv.Atoi(req.Session.Customer.ID)
	referralHistory, err := s.opt.Client.CampaignServiceGrpc.GetReferralHistory(ctx, &campaign_service.GetReferralHistoryRequest{
		ReferrerId: int64(customerID),
	})
	res.Summary.TotalPoint = referralHistory.TotalPoint
	res.Summary.TotalReferral = referralHistory.TotalReferral
	layout := "2006-01-02 15:04:05"

	for _, v := range referralHistory.DataReferral {
		res.Detail.ReferralList = append(res.Detail.ReferralList, &dto.ReferralList{
			Name:      v.Name,
			CreatedAt: v.CreatedAt.AsTime().Format(layout),
		})
	}
	for _, v := range referralHistory.DataReferralPoint {
		res.Detail.ReferralPointList = append(res.Detail.ReferralPointList, &dto.ReferralPointList{
			Name:       v.Name,
			CreatedAt:  v.CreatedAt.AsTime().Format(layout),
			PointValue: strconv.FormatFloat(float64(v.PointValue), 'f', 1, 64),
		})
	}

	return
}
