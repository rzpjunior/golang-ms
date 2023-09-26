package repository

import (
	"context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/generator"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type ISubDistrictRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (subDistricts []*model.SubDistrict, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (subDistrict *model.SubDistrict, err error)
}

type SubDistrictRepository struct {
	opt opt.Options
}

func NewSubDistrictRepository() ISubDistrictRepository {
	return &SubDistrictRepository{
		opt: global.Setup.Common,
	}
}

func (r *SubDistrictRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (subDistricts []*model.SubDistrict, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SubDistrictRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.SubDistrict))

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

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &subDistricts)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SubDistrictRepository) GetDetail(ctx context.Context, id int64, code string) (subDistrict *model.SubDistrict, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SubDistrictRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	return r.MockDatas(1)[0], nil

	subDistrict = &model.SubDistrict{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		subDistrict.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		subDistrict.Code = code
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, subDistrict, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SubDistrictRepository) MockDatas(total int) (mockDatas []*model.SubDistrict) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.SubDistrict{
				ID:          int64(i),
				Code:        fmt.Sprintf("SUB%d", i),
				Description: fmt.Sprintf("Dummy SubDistrict %d", i),
				Status:      1,
				CreatedAt:   generator.DummyTime(),
				UpdatedAt:   generator.DummyTime(),
			})
	}
	return
}
