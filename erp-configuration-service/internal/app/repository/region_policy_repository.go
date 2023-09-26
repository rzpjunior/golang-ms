package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/model"
)

type IRegionPolicyRepository interface {
	Get(ctx context.Context, offset int, limit int, search string, orderBy string, regionID string) (regionPolicys []*model.RegionPolicy, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string, regionId string) (regionPolicy *model.RegionPolicy, err error)
	Update(ctx context.Context, RegionPolicy *model.RegionPolicy, columns ...string) (err error)
}

type RegionPolicyRepository struct {
	opt opt.Options
}

func NewRegionPolicyRepository() IRegionPolicyRepository {
	return &RegionPolicyRepository{
		opt: global.Setup.Common,
	}
}

func (r *RegionPolicyRepository) Get(ctx context.Context, offset int, limit int, search string, orderBy string, regionID string) (regionPolicys []*model.RegionPolicy, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "RegionPolicyRepository.Get")
	defer span.End()

	db := r.opt.Database.Read
	qs := db.QueryTable(new(model.RegionPolicy))

	cond := orm.NewCondition()

	if search != "" {
		cond = cond.And("region__icontains", search)
	}

	if regionID != "" {
		cond = cond.And("region_id_gp__icontains", regionID)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &regionPolicys)
	if err != nil {
		span.RecordError(err)
		return
	}

	return regionPolicys, count, nil
}

func (r *RegionPolicyRepository) GetDetail(ctx context.Context, id int64, code string, regionId string) (regionPolicy *model.RegionPolicy, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "RegionPolicyRepository.GetByID")
	defer span.End()

	db := r.opt.Database.Read

	regionPolicy = &model.RegionPolicy{}
	var cols []string
	if id != 0 {
		regionPolicy.ID = id
		cols = append(cols, "id")
	}
	if regionId != "" {
		regionPolicy.RegionIDGP = regionId
		cols = append(cols, "region_id_gp")
	}

	err = db.ReadWithCtx(ctx, regionPolicy, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *RegionPolicyRepository) Update(ctx context.Context, RegionPolicy *model.RegionPolicy, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "RegionPolicyRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, RegionPolicy, columns...)
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
