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

type IUomRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (uoms []*model.Uom, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (uom *model.Uom, err error)
}

type UomRepository struct {
	opt opt.Options
}

func NewUomRepository() IUomRepository {
	return &UomRepository{
		opt: global.Setup.Common,
	}
}

func (r *UomRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (uoms []*model.Uom, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UomRepository.Get")
	defer span.End()

	// RETURN DUMMY
	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.Uom))

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

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &uoms)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *UomRepository) GetDetail(ctx context.Context, id int64, code string) (uom *model.Uom, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UomRepository.GetDetail")
	defer span.End()

	// RETURN DUMMY
	return r.MockDatas(1)[0], nil

	uom = &model.Uom{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		uom.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		uom.Code = code
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, uom, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *UomRepository) MockDatas(total int) (mockDatas []*model.Uom) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.Uom{
				ID:             int64(i),
				Code:           fmt.Sprintf("UOM%d", i),
				Description:    fmt.Sprintf("Dummy UOM %d", i),
				Status:         1,
				DecimalEnabled: 1,
				CreatedAt:      generator.DummyTime(),
				UpdatedAt:      generator.DummyTime(),
			})
	}
	return
}
