package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/repository"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
)

type IMembershipService interface {
	GetMembershipLevelList(ctx context.Context, req *dto.MembershipLevelRequestGet) (res []*dto.MembershipLevelResponse, total int64, err error)
	GetMembeshipLevelDetail(ctx context.Context, id int64, level int8) (res dto.MembershipLevelResponse, err error)
	GetMembershipCheckpointList(ctx context.Context, req *dto.MembershipCheckpointRequestGet) (res []*dto.MembershipCheckpointResponse, total int64, err error)
	GetMembershipCheckpointDetail(ctx context.Context, id int64, checkpoint int8) (res dto.MembershipCheckpointResponse, err error)
	GetMembershipAdvantageDetail(ctx context.Context, req *campaign_service.GetMembershipAdvantageDetailRequest) (res dto.MembershipAdvantage, err error)
	GetMembershipLevelAdvantageList(ctx context.Context, req *campaign_service.GetMembershipLevelAdvantageListRequest) (res []*dto.MembershipLevelAdvantage, total int64, err error)
	GetMembershipRewardList(ctx context.Context, req *campaign_service.GetMembershipRewardListRequest) (res []*dto.MembershipReward, total int64, err error)
	GetMembershipRewardDetail(ctx context.Context, req *campaign_service.GetMembershipRewardDetailRequest) (res *dto.MembershipReward, err error)
	GetCustomerMembership(ctx context.Context, req *campaign_service.GetCustomerMembershipDetailRequest) (res *dto.CustomerMembership, err error)
}

type MembershipService struct {
	opt                  opt.Options
	RepositoryMembership repository.IMembershipRepository
	ServicesTalon        ITalonService
}

func NewMembershipService() IMembershipService {
	return &MembershipService{
		opt:                  global.Setup.Common,
		RepositoryMembership: repository.NewMembershipRepository(),
		ServicesTalon:        NewTalonService(),
	}
}

func (s *MembershipService) GetMembershipLevelList(ctx context.Context, req *dto.MembershipLevelRequestGet) (res []*dto.MembershipLevelResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "MembershipLevelService.Get")
	defer span.End()

	var membershipLevels []*model.MembershipLevel
	membershipLevels, total, err = s.RepositoryMembership.GetMembershipLevelList(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, membershipLevel := range membershipLevels {
		res = append(res, &dto.MembershipLevelResponse{
			ID:       membershipLevel.ID,
			Code:     membershipLevel.Code,
			Level:    membershipLevel.Level,
			Name:     membershipLevel.Name,
			ImageUrl: membershipLevel.ImageUrl,
			Status:   membershipLevel.Status,
		})
	}

	return
}

func (s *MembershipService) GetMembeshipLevelDetail(ctx context.Context, id int64, level int8) (res dto.MembershipLevelResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "MembershipLevelService.GetByID")
	defer span.End()

	var membershipLevel *model.MembershipLevel

	membershipLevel, err = s.RepositoryMembership.GetMembeshipLevelDetail(ctx, id, level)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.MembershipLevelResponse{
		ID:       membershipLevel.ID,
		Code:     membershipLevel.Code,
		Level:    membershipLevel.Level,
		Name:     membershipLevel.Name,
		ImageUrl: membershipLevel.ImageUrl,
		Status:   membershipLevel.Status,
	}

	return
}

func (s *MembershipService) GetMembershipCheckpointList(ctx context.Context, req *dto.MembershipCheckpointRequestGet) (res []*dto.MembershipCheckpointResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "MembershipCheckpointService.Get")
	defer span.End()

	var membershipCheckpoints []*model.MembershipCheckpoint
	membershipCheckpoints, total, err = s.RepositoryMembership.GetMembershipCheckpointList(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, membershipCheckpoint := range membershipCheckpoints {
		res = append(res, &dto.MembershipCheckpointResponse{
			ID:                membershipCheckpoint.ID,
			Checkpoint:        membershipCheckpoint.Checkpoint,
			TargetAmount:      membershipCheckpoint.TargetAmount,
			MembershipLevelID: membershipCheckpoint.MembershipLevelID,
			Status:            membershipCheckpoint.Status,
		})
	}

	return
}

func (s *MembershipService) GetMembershipCheckpointDetail(ctx context.Context, id int64, checkpoint int8) (res dto.MembershipCheckpointResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "MembershipCheckpointService.GetByID")
	defer span.End()

	var membershipCheckpoint *model.MembershipCheckpoint

	membershipCheckpoint, err = s.RepositoryMembership.GetMembershipCheckpointDetail(ctx, id, checkpoint)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.MembershipCheckpointResponse{
		ID:                membershipCheckpoint.ID,
		Checkpoint:        membershipCheckpoint.Checkpoint,
		TargetAmount:      membershipCheckpoint.TargetAmount,
		MembershipLevelID: membershipCheckpoint.MembershipLevelID,
		Status:            membershipCheckpoint.Status,
	}

	return
}

func (s *MembershipService) GetMembershipAdvantageDetail(ctx context.Context, req *campaign_service.GetMembershipAdvantageDetailRequest) (res dto.MembershipAdvantage, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "MembershipCheckpointService.GetByID")
	defer span.End()

	var membershipAdvantage *model.MembershipAdvantage

	membershipAdvantage, err = s.RepositoryMembership.GetMembershipAdvantageDetail(ctx, req)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.MembershipAdvantage{
		ID:          membershipAdvantage.ID,
		Name:        membershipAdvantage.Name,
		Description: membershipAdvantage.Description,
		ImageUrl:    membershipAdvantage.ImageUrl,
		LinkUrl:     membershipAdvantage.LinkUrl,
		Status:      membershipAdvantage.Status,
	}

	return
}

func (s *MembershipService) GetMembershipLevelAdvantageList(ctx context.Context, req *campaign_service.GetMembershipLevelAdvantageListRequest) (res []*dto.MembershipLevelAdvantage, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "MembershipCheckpointService.GetByID")
	defer span.End()

	var membershipLevelAdvantage []*model.MembershipLevelAdvantage

	membershipLevelAdvantage, total, err = s.RepositoryMembership.GetMembershipLevelAdvantageList(ctx, req)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, v := range membershipLevelAdvantage {
		res = append(res, &dto.MembershipLevelAdvantage{
			ID:                    v.ID,
			MembershipLevelID:     v.MembershipLevelID,
			MembershipAdvantageID: v.MembershipAdvantageID,
		})
	}

	return
}

func (s *MembershipService) GetMembershipRewardList(ctx context.Context, req *campaign_service.GetMembershipRewardListRequest) (res []*dto.MembershipReward, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "MembershipCheckpointService.GetByID")
	defer span.End()

	var membershipReward []*model.MembershipReward

	membershipReward, total, err = s.RepositoryMembership.GetMembershipRewardList(ctx, req)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, v := range membershipReward {
		res = append(res, &dto.MembershipReward{
			ID:                 v.ID,
			OpenedImageUrl:     v.OpenedImageUrl,
			ClosedImageUrl:     v.ClosedImageUrl,
			BackgroundImageUrl: v.BackgroundImageUrl,
			RewardLevel:        v.RewardLevel,
			MaxAmount:          v.MaxAmount,
			Status:             v.Status,
			Description:        v.Description,
			IsPassed:           v.IsPassed,
			CurrentPercentage:  v.CurrentPercentage,
			RemainingAmount:    v.RemainingAmount,
		})
	}

	return
}

func (s *MembershipService) GetMembershipRewardDetail(ctx context.Context, req *campaign_service.GetMembershipRewardDetailRequest) (res *dto.MembershipReward, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "MembershipCheckpointService.GetByID")
	defer span.End()

	var membershipReward *model.MembershipReward

	membershipReward, err = s.RepositoryMembership.GetMembershipRewardDetail(ctx, req)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.MembershipReward{
		ID:                 membershipReward.ID,
		OpenedImageUrl:     membershipReward.OpenedImageUrl,
		ClosedImageUrl:     membershipReward.ClosedImageUrl,
		BackgroundImageUrl: membershipReward.BackgroundImageUrl,
		RewardLevel:        membershipReward.RewardLevel,
		MaxAmount:          membershipReward.MaxAmount,
		Status:             membershipReward.Status,
		Description:        membershipReward.Description,
		IsPassed:           membershipReward.IsPassed,
		CurrentPercentage:  membershipReward.CurrentPercentage,
		RemainingAmount:    membershipReward.RemainingAmount,
	}

	return
}

func (s *MembershipService) GetCustomerMembership(ctx context.Context, req *campaign_service.GetCustomerMembershipDetailRequest) (res *dto.CustomerMembership, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "MembershipCheckpointService.GetByID")
	defer span.End()

	//var profilecode string
	customerProfile, err := s.ServicesTalon.GetCustomerProfile(req.ProfileCode)
	if err != nil {
		return nil, err
	}
	res, err = s.RepositoryMembership.GetCustomerMembership(ctx, req, customerProfile)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
