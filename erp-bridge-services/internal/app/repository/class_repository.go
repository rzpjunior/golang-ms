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

type IClassRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (classes []*model.Class, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (class *model.Class, err error)
}

type ClassRepository struct {
	opt opt.Options
}

func NewClassRepository() IClassRepository {
	return &ClassRepository{
		opt: global.Setup.Common,
	}
}

func (r *ClassRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (classes []*model.Class, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ClassRepository.Get")
	defer span.End()

	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.Class))

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

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &classes)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ClassRepository) GetDetail(ctx context.Context, id int64, code string) (class *model.Class, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ClassRepository.GetDetail")
	defer span.End()

	return r.MockDatas(1)[0], nil

	class = &model.Class{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		class.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		class.Code = code
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, class, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ClassRepository) MockDatas(total int) (mockDatas []*model.Class) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.Class{
				ID:          int64(i),
				Code:        fmt.Sprintf("CRV%d", i),
				Description: fmt.Sprintf("Dummy Class %d", i),
				Status:      1,
				CreatedAt:   generator.DummyTime(),
				UpdatedAt:   generator.DummyTime(),
			})
	}
	return
}
