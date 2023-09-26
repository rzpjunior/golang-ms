package repository

import (
	"context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"
)

type ITermConditionRepository interface {
	Get(ctx context.Context, offset int, limit int) (termConditions []*model.TermCondition, count int64, err error)
	GetDetail(ctx context.Context, id int64) (termCondition *model.TermCondition, err error)
}

type TermConditionRepository struct {
	opt opt.Options
}

func NewTermConditionRepository() ITermConditionRepository {
	return &TermConditionRepository{
		opt: global.Setup.Common,
	}
}

func (r *TermConditionRepository) Get(ctx context.Context, offset int, limit int) (termConditions []*model.TermCondition, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "TermConditionRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	// dummies := r.MockDatas(10)
	// return dummies, int64(len(dummies)), nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.TermCondition))

	cond := orm.NewCondition()

	qs = qs.SetCond(cond)

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &termConditions)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *TermConditionRepository) GetDetail(ctx context.Context, id int64) (termCondition *model.TermCondition, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "TermConditionRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	// return r.MockDatas(1)[0], nil

	termCondition = &model.TermCondition{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		termCondition.ID = id
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, termCondition, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *TermConditionRepository) MockDatas(total int) (mockDatas []*model.TermCondition) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.TermCondition{
				ID:          int64(i),
				Application: int8(i),
				Version:     string(i),
				Title:       fmt.Sprintf("Dummy Term and Condition %d", i),
				TitleValue:  string(i),
				Content:     string(i),
			})
	}

	return mockDatas
}
