package repository

import (
	"context"
	"reflect"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
)

type IMembershipRepository interface {
	GetMembershipLevelList(ctx context.Context, req *dto.MembershipLevelRequestGet) (membershipLevels []*model.MembershipLevel, count int64, err error)
	GetMembeshipLevelDetail(ctx context.Context, id int64, level int8) (membershipLevel *model.MembershipLevel, err error)
	GetMembershipCheckpointList(ctx context.Context, req *dto.MembershipCheckpointRequestGet) (membershipCheckpoints []*model.MembershipCheckpoint, count int64, err error)
	GetMembershipCheckpointDetail(ctx context.Context, id int64, checkpoint int8) (membershipCheckpoint *model.MembershipCheckpoint, err error)
	GetMembershipAdvantageDetail(ctx context.Context, req *campaign_service.GetMembershipAdvantageDetailRequest) (membershipAdvantage *model.MembershipAdvantage, err error)
	GetMembershipLevelAdvantageList(ctx context.Context, req *campaign_service.GetMembershipLevelAdvantageListRequest) (membershipLevelAdvantage []*model.MembershipLevelAdvantage, count int64, err error)
	GetMembershipRewardList(ctx context.Context, req *campaign_service.GetMembershipRewardListRequest) (membershipLevelAdvantage []*model.MembershipReward, count int64, err error)
	GetMembershipRewardDetail(ctx context.Context, req *campaign_service.GetMembershipRewardDetailRequest) (membershipReward *model.MembershipReward, err error)
	GetCustomerMembership(ctx context.Context, req *campaign_service.GetCustomerMembershipDetailRequest, customerProfile *dto.CustomerProfileData) (membership *dto.CustomerMembership, err error)
}

type MembershipRepository struct {
	opt opt.Options
}

func NewMembershipRepository() IMembershipRepository {
	return &MembershipRepository{
		opt: global.Setup.Common,
	}
}

func (r *MembershipRepository) GetMembershipLevelList(ctx context.Context, req *dto.MembershipLevelRequestGet) (membershipLevels []*model.MembershipLevel, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "MembershipLevelRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.MembershipLevel))

	cond := orm.NewCondition()

	if req.Search != "" {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("name__icontains", req.Search).Or("code__icontains", req.Search)
		cond = cond.AndCond(cond1)
	}

	if req.Status != 0 {
		cond = cond.And("status", req.Status)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &membershipLevels)

	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *MembershipRepository) GetMembeshipLevelDetail(ctx context.Context, id int64, level int8) (membershipLevel *model.MembershipLevel, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "MembershipLevelRepository.GetByID")
	defer span.End()

	membershipLevel = &model.MembershipLevel{}

	var cols []string
	if id != 0 {
		membershipLevel.ID = id
		cols = append(cols, "id")
	}
	if level != 0 {
		membershipLevel.Level = level
		cols = append(cols, "level")
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, membershipLevel, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *MembershipRepository) GetMembershipCheckpointList(ctx context.Context, req *dto.MembershipCheckpointRequestGet) (membershipCheckpoints []*model.MembershipCheckpoint, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "MembershipCheckpointRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.MembershipCheckpoint))

	cond := orm.NewCondition()

	if req.Status != 0 {
		cond = cond.And("status", req.Status)
	}

	if req.MembershipLevelID != 0 {
		cond = cond.And("membership_level_id", req.MembershipLevelID)
	}
	if req.ID != 0 {
		cond = cond.And("id", req.ID)
	}
	cond = cond.AndNot("target_amount", 0)
	//
	// if req.TargetAmount != 0 {
	// 	cond = cond.And("target", req.MembershipLevelID)
	// }

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &membershipCheckpoints)

	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *MembershipRepository) GetMembershipCheckpointDetail(ctx context.Context, id int64, checkpoint int8) (membershipCheckpoint *model.MembershipCheckpoint, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "MembershipCheckpointRepository.GetByID")
	defer span.End()

	membershipCheckpoint = &model.MembershipCheckpoint{}

	var cols []string
	if id != 0 {
		membershipCheckpoint.ID = id
		cols = append(cols, "id")
	} else {
		membershipCheckpoint.Checkpoint = checkpoint
		cols = append(cols, "checkpoint")
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, membershipCheckpoint, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *MembershipRepository) GetMembershipAdvantageDetail(ctx context.Context, req *campaign_service.GetMembershipAdvantageDetailRequest) (membershipAdvantage *model.MembershipAdvantage, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "MembershipCheckpointRepository.GetMembershipAdvantageDetail")
	defer span.End()

	membershipAdvantage = &model.MembershipAdvantage{}

	var cols []string
	if req.Id != 0 {
		membershipAdvantage.ID = req.Id
		cols = append(cols, "id")
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, membershipAdvantage, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *MembershipRepository) GetMembershipLevelAdvantageList(ctx context.Context, req *campaign_service.GetMembershipLevelAdvantageListRequest) (membershipLevelAdvantage []*model.MembershipLevelAdvantage, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "MembershipCheckpointRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.MembershipLevelAdvantage))

	cond := orm.NewCondition()

	if req.MembershipLevelId != 0 {
		cond = cond.And("membership_level_id", req.MembershipLevelId)
	}

	if req.MembershipAdvantageId != 0 {
		cond = cond.And("membership_advantage_id", req.MembershipAdvantageId)
	}

	qs = qs.SetCond(cond)

	count, err = qs.AllWithCtx(ctx, &membershipLevelAdvantage)

	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *MembershipRepository) GetMembershipRewardList(ctx context.Context, req *campaign_service.GetMembershipRewardListRequest) (membershipLevelAdvantage []*model.MembershipReward, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "MembershipCheckpointRepository.GetMembershipRewardList")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.MembershipReward))

	cond := orm.NewCondition()

	if req.Id != 0 {
		cond = cond.And("id", req.Id)
	}

	cond = cond.And("status", 1)
	qs = qs.SetCond(cond)
	qs = qs.OrderBy("reward_level")
	count, err = qs.AllWithCtx(ctx, &membershipLevelAdvantage)

	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *MembershipRepository) GetMembershipRewardDetail(ctx context.Context, req *campaign_service.GetMembershipRewardDetailRequest) (membershipReward *model.MembershipReward, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "MembershipCheckpointRepository.GetMembershipRewardList")
	defer span.End()

	membershipReward = &model.MembershipReward{}
	//membershipReward2 := &model.MembershipReward{}

	var cols []string
	if req.Id != 0 {
		membershipReward.ID = req.Id
		cols = append(cols, "id")
	}
	membershipReward.Status = 1
	cols = append(cols, "status")
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, membershipReward, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	var max_ammount float64
	db.Raw("select max_amount "+
		"from membership_reward "+
		"where status = 1 and reward_level < ? "+
		"order by reward_level desc "+
		"limit 1", req.RewardLevel).QueryRow(&max_ammount)

	if req.MembershipRewardAmount == 0 {
		membershipReward.CurrentPercentage = 0
		membershipReward.RemainingAmount = 0
		return
	}

	// count percentage of current level reward
	membershipReward.CurrentPercentage = (req.MembershipRewardAmount - max_ammount) / (membershipReward.MaxAmount - max_ammount) * 100

	// count remaining amount to go to next level of membership reward
	membershipReward.RemainingAmount = membershipReward.MaxAmount - req.MembershipRewardAmount

	// if reward amount has already exceed highest level max amount then set percentage to full
	if req.MembershipRewardAmount > membershipReward.MaxAmount {
		membershipReward.CurrentPercentage = 100
		membershipReward.RemainingAmount = 0
	}

	return
}

func (r *MembershipRepository) GetCustomerMembership(ctx context.Context, req *campaign_service.GetCustomerMembershipDetailRequest, customerProfile *dto.CustomerProfileData) (membership *dto.CustomerMembership, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "MembershipCheckpointRepository.GetCustomerMembership")
	defer span.End()
	var (
		//customerProfile                                *model.CustomerProfileData
		q, membershipName                              string
		currentAmount                                  float64
		membershipLevel                                *model.MembershipLevel
		mcData                                         []*model.MembershipCheckpoint
		targetAmount, prevTarget, checkpointPercentage float64
		checkpoint, maxLevel                           int8
	)
	db := r.opt.Database.Read
	q = "select ml.* " +
		"from membership_level ml " +
		"where ml.id = ? and ml.status = 1"
	if err = db.Raw(q, req.MembershipLevelId).QueryRow(&membershipLevel); err != nil {
		return nil, err
	}

	attributes := reflect.ValueOf(customerProfile.Profile.Attributes)
	for _, v := range attributes.MapKeys() {
		if v.String() == "membership_level" {
			membershipName = attributes.MapIndex(v).Interface().(string)
			continue
		}

		if v.String() == "fresh_product_revenue" {
			currentAmount = attributes.MapIndex(v).Interface().(float64)
			continue
		}
	}

	q = "select target_amount from membership_checkpoint where membership_level_id = ? order by id desc limit 1"
	if err = db.Raw(q, req.MembershipLevelId).QueryRow(&targetAmount); err != nil {
		return nil, err
	}

	q = "select * from membership_checkpoint where status = 1"
	if _, err = db.Raw(q).QueryRows(&mcData); err != nil {
		return nil, err
	}

	for i, v := range mcData[:len(mcData)-1] {
		if currentAmount >= v.TargetAmount {
			continue
		}

		if i-1 >= 0 {
			prevTarget = mcData[i-1].TargetAmount
		}

		checkpoint = v.Checkpoint
		checkpointPercentage = (currentAmount - prevTarget) / (v.TargetAmount - prevTarget) * 100
		break
	}

	q = "select max(level) from membership_level where status = 1"
	if err = db.Raw(q).QueryRow(&maxLevel); err != nil {
		return nil, err
	}

	if membershipLevel.Level == maxLevel && checkpoint == 0 {
		checkpoint = mcData[len(mcData)-1].Checkpoint
		checkpointPercentage = 100
	}

	membership = &dto.CustomerMembership{
		MembershipLevel:      membershipLevel.Level,
		MembershipLevelName:  membershipName,
		MembershipCheckpoint: checkpoint,
		CurrentAmount:        currentAmount,
		TargetAmount:         targetAmount,
		CheckpointPercentage: checkpointPercentage,
	}
	return
}
