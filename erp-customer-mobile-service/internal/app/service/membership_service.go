package service

import (
	"context"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
)

type IMembershipService interface {
	Get(ctx context.Context, req dto.RequestGetMembershipList) (res []dto.ResponseMembershipLevelList, err error)
	GetReward(ctx context.Context, req dto.RequestGetRewardList) (res []dto.MembershipRewardList, err error)
	GetRewardDetail(ctx context.Context, req dto.RequestGetRewardList) (res dto.MembershipRewardList, err error)
}

type MembershipService struct {
	opt opt.Options
	//RepositoryOTPOutgoing repository.IOtpOutgoingRepository
}

func NewMembershipService() IMembershipService {
	return &MembershipService{
		opt: global.Setup.Common,
		//RepositoryOTPOutgoing: repository.NewOtpOutgoingRepository(),
	}
}

func (s *MembershipService) Get(ctx context.Context, req dto.RequestGetMembershipList) (res []dto.ResponseMembershipLevelList, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "MembershipService.Get")
	defer span.End()

	level, _ := strconv.Atoi(req.Data.Level)
	checkpoint, _ := strconv.Atoi(req.Data.Checkpoint)
	//get all
	if level == 0 {
		if checkpoint != 0 {
			//throw error here
		}

		memberships, _ := s.opt.Client.CampaignServiceGrpc.GetMembershipLevelList(ctx, &campaign_service.GetMembershipLevelListRequest{})
		for _, v := range memberships.Data {
			var tempMA []*dto.MembershipAdvantage
			var tempCP []*dto.MembershipCheckpoint

			checkpoints, _ := s.opt.Client.CampaignServiceGrpc.GetMembershipCheckpointList(ctx, &campaign_service.GetMembershipCheckpointListRequest{
				MembershipLevelId: v.Id,
			})
			for _, y := range checkpoints.Data {
				tempCP = append(tempCP, &dto.MembershipCheckpoint{
					ID:                y.Id,
					Checkpoint:        strconv.Itoa(int(y.Checkpoint)),
					TargetAmount:      strconv.FormatFloat(y.TargetAmount, 'f', 1, 64),
					Status:            strconv.Itoa(int(y.Status)),
					MembershipLevelID: strconv.Itoa(int(y.MembershipLevelId)),
				})
			}

			levelAdvantages, _ := s.opt.Client.CampaignServiceGrpc.GetMembershipLevelAdvantageList(ctx, &campaign_service.GetMembershipLevelAdvantageListRequest{
				MembershipLevelId: int64(v.Level),
			})
			for _, x := range levelAdvantages.Data {
				advantage, _ := s.opt.Client.CampaignServiceGrpc.GetMembershipAdvantageDetail(ctx, &campaign_service.GetMembershipAdvantageDetailRequest{
					Id: x.MembershipAdvantageId,
				})
				tempMA = append(tempMA, &dto.MembershipAdvantage{
					ID:          advantage.Data.Id,
					Name:        advantage.Data.Name,
					Description: advantage.Data.Description,
					ImageUrl:    advantage.Data.ImageUrl,
					LinkUrl:     advantage.Data.LinkUrl,
					Status:      strconv.Itoa(int(advantage.Data.Status)),
				})

			}
			res = append(res, dto.ResponseMembershipLevelList{
				ID:                    v.Id,
				Code:                  v.Code,
				Name:                  v.Name,
				Level:                 strconv.Itoa(int(v.Level)),
				ImageUrl:              v.ImageUrl,
				Status:                strconv.Itoa(int(v.Status)),
				MembershipAdvantages:  tempMA,
				MembershipCheckpoints: tempCP,
			})
		}
		return res, err
		//return res
	} else {
		//get by level
		var tempMA []*dto.MembershipAdvantage
		var tempCP []*dto.MembershipCheckpoint

		memberships, _ := s.opt.Client.CampaignServiceGrpc.GetMembershipLevelDetail(ctx, &campaign_service.GetMembershipLevelDetailRequest{
			Level: int64(level),
		})
		checkpoints, _ := s.opt.Client.CampaignServiceGrpc.GetMembershipCheckpointList(ctx, &campaign_service.GetMembershipCheckpointListRequest{
			MembershipLevelId: memberships.Data.Id,
			Id:                int64(checkpoint),
		})
		for _, y := range checkpoints.Data {
			tempCP = append(tempCP, &dto.MembershipCheckpoint{
				ID:                y.Id,
				Checkpoint:        strconv.Itoa(int(y.Checkpoint)),
				TargetAmount:      strconv.FormatFloat(y.TargetAmount, 'f', 1, 64),
				Status:            strconv.Itoa(int(y.Status)),
				MembershipLevelID: strconv.Itoa(int(y.MembershipLevelId)),
			})
		}

		levelAdvantages, _ := s.opt.Client.CampaignServiceGrpc.GetMembershipLevelAdvantageList(ctx, &campaign_service.GetMembershipLevelAdvantageListRequest{
			MembershipLevelId: int64(memberships.Data.Level),
		})
		for _, x := range levelAdvantages.Data {
			advantage, _ := s.opt.Client.CampaignServiceGrpc.GetMembershipAdvantageDetail(ctx, &campaign_service.GetMembershipAdvantageDetailRequest{
				Id: x.MembershipAdvantageId,
			})
			tempMA = append(tempMA, &dto.MembershipAdvantage{
				ID:          advantage.Data.Id,
				Name:        advantage.Data.Name,
				Description: advantage.Data.Description,
				ImageUrl:    advantage.Data.ImageUrl,
				LinkUrl:     advantage.Data.LinkUrl,
				Status:      strconv.Itoa(int(advantage.Data.Status)),
			})
		}
		res = append(res, dto.ResponseMembershipLevelList{
			ID:                    memberships.Data.Id,
			Code:                  memberships.Data.Code,
			Name:                  memberships.Data.Name,
			Level:                 strconv.Itoa(int(memberships.Data.Level)),
			ImageUrl:              memberships.Data.ImageUrl,
			Status:                strconv.Itoa(int(memberships.Data.Status)),
			MembershipAdvantages:  tempMA,
			MembershipCheckpoints: tempCP,
		})
		return res, err
	}

	// return
}

func (s *MembershipService) GetReward(ctx context.Context, req dto.RequestGetRewardList) (res []dto.MembershipRewardList, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "MembershipService.Get")
	defer span.End()

	reward, err := s.opt.Client.CampaignServiceGrpc.GetMembershipRewardList(ctx, &campaign_service.GetMembershipRewardListRequest{})
	if err != nil {
		return nil, err
	}
	for _, v := range reward.Data {
		res = append(res, dto.MembershipRewardList{
			ID:                 strconv.Itoa(int(v.Id)),
			OpenedImageUrl:     v.OpenedImageUrl,
			ClosedImageUrl:     v.ClosedImageUrl,
			BackgroundImageUrl: v.BackgroundImageUrl,
			RewardLevel:        strconv.Itoa(int(v.RewardLevel)),
			MaxAmount:          strconv.FormatFloat(v.MaxAmount, 'f', 1, 64),
			Status:             strconv.Itoa(int(v.Status)),
			Description:        v.Description,
			IsPassed:           strconv.Itoa(int(v.IsPassed)),
			CurrentPercentage:  strconv.FormatFloat(v.CurrentPercentage, 'f', 1, 64),
			RemainingAmount:    strconv.FormatFloat(v.RemainingAmount, 'f', 1, 64),
		})
	}
	membershipRewardID, _ := strconv.Atoi(req.Session.Customer.MembershipRewardID)
	//membershipRewardID := req.Session.Customer.MembershipReward.ID
	if membershipRewardID == 0 {
		return res, err
	}

	rewards, e := s.opt.Client.CampaignServiceGrpc.GetMembershipRewardList(ctx, &campaign_service.GetMembershipRewardListRequest{
		Id: int64(membershipRewardID),
	})
	if e != nil {
		return res, e
	}
	if len(rewards.Data) == 0 {
		return res, e
	}
	req.Session.Customer.MembershipReward = &model.MembershipReward{
		ID:                 rewards.Data[0].Id,
		OpenedImageUrl:     rewards.Data[0].OpenedImageUrl,
		ClosedImageUrl:     rewards.Data[0].ClosedImageUrl,
		BackgroundImageUrl: rewards.Data[0].BackgroundImageUrl,
		RewardLevel:        int8(rewards.Data[0].RewardLevel),
		MaxAmount:          rewards.Data[0].MaxAmount,
		Status:             int8(rewards.Data[0].Status),
		Description:        rewards.Data[0].Description,
		IsPassed:           int8(rewards.Data[0].IsPassed),
		CurrentPercentage:  rewards.Data[0].CurrentPercentage,
		RemainingAmount:    rewards.Data[0].RemainingAmount,
	}

	membershipRewardAmount, _ := strconv.ParseFloat(req.Session.Customer.MembershipRewardAmount, 64)

	// compare current merchant reward level with membership reward level
	for i, v := range res {
		maxAmount, _ := strconv.ParseFloat(v.MaxAmount, 64)

		if maxAmount <= membershipRewardAmount {
			res[i].IsPassed = "1"
			continue
		}
		break
	}

	return
}

func (s *MembershipService) GetRewardDetail(ctx context.Context, req dto.RequestGetRewardList) (res dto.MembershipRewardList, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "MembershipService.Get")
	defer span.End()
	membershipRewardID, _ := strconv.Atoi(req.Session.Customer.MembershipRewardID)
	// membershipRewardID := req.Session.Customer.MembershipReward.ID
	membershipRewardAmount, _ := strconv.ParseFloat(req.Session.Customer.MembershipRewardAmount, 64)

	reward, err := s.opt.Client.CampaignServiceGrpc.GetMembershipRewardDetail(ctx, &campaign_service.GetMembershipRewardDetailRequest{
		Id:                     int64(membershipRewardID),
		RewardLevel:            int64(membershipRewardID),
		MembershipRewardAmount: membershipRewardAmount,
	})
	if err != nil {
		return res, err
	}

	res = dto.MembershipRewardList{
		ID:                 strconv.Itoa(int(reward.Data.Id)),
		OpenedImageUrl:     reward.Data.OpenedImageUrl,
		ClosedImageUrl:     reward.Data.ClosedImageUrl,
		BackgroundImageUrl: reward.Data.BackgroundImageUrl,
		RewardLevel:        strconv.Itoa(int(reward.Data.RewardLevel)),
		MaxAmount:          strconv.FormatFloat(reward.Data.MaxAmount, 'f', 1, 64),
		Status:             strconv.Itoa(int(reward.Data.Status)),
		Description:        reward.Data.Description,
		IsPassed:           strconv.Itoa(int(reward.Data.IsPassed)),
		CurrentPercentage:  strconv.FormatFloat(reward.Data.CurrentPercentage, 'f', 1, 64),
		RemainingAmount:    strconv.FormatFloat(reward.Data.RemainingAmount, 'f', 1, 64),
	}

	return
}
