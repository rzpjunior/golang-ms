package repository

import (
	"context"
	"fmt"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type IRegionRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (regions []*model.Region, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (region *model.Region, err error)
}

type RegionRepository struct {
	opt opt.Options
}

func NewRegionRepository() IRegionRepository {
	return &RegionRepository{
		opt: global.Setup.Common,
	}
}

func (r *RegionRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (regions []*model.Region, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "RegionRepository.Get")
	defer span.End()

	count = 10
	for i := 1; i <= int(count); i++ {
		regions = append(regions, r.MockDatas(int64(i)))
	}
	return

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.Region))

	cond := orm.NewCondition()

	if search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("description__icontains", search).Or("code__icontains", search)
		cond = cond.AndCond(condGroup)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &regions)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *RegionRepository) GetDetail(ctx context.Context, id int64, code string) (region *model.Region, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "RegionRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	region = r.MockDatas(id)
	return

	region = &model.Region{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		region.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		region.Code = code
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, region, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *RegionRepository) MockDatas(id int64) (mockDatas *model.Region) {
	if id == 0 {
		id = 1
	}
	mockDatas = &model.Region{
		ID:          int64(id),
		Code:        fmt.Sprintf("REG%d", id),
		Description: fmt.Sprintf("Dummy Region %d", id),
		Status:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return mockDatas
}
