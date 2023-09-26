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

type ISalespersonRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (salesPersons []*model.Salesperson, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (salesperson *model.Salesperson, err error)
}

type SalespersonRepository struct {
	opt opt.Options
}

func NewSalespersonRepository() ISalespersonRepository {
	return &SalespersonRepository{
		opt: global.Setup.Common,
	}
}

func (r *SalespersonRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (salesPersons []*model.Salesperson, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalespersonRepository.Get")
	defer span.End()

	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.Salesperson))

	cond := orm.NewCondition()

	if search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("firstname__icontains", search).Or("middlename__icontains", search).Or("lastname__icontains", search)
		cond = cond.AndCond(condGroup)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &salesPersons)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalespersonRepository) GetDetail(ctx context.Context, id int64, code string) (salesperson *model.Salesperson, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalespersonRepository.GetDetail")
	defer span.End()

	return r.MockDatas(1)[0], nil

	salesperson = &model.Salesperson{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		salesperson.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		salesperson.Code = code
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, salesperson, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalespersonRepository) MockDatas(total int) (mockDatas []*model.Salesperson) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.Salesperson{
				ID:         int64(i),
				Code:       fmt.Sprintf("SPR%d", i),
				FirstName:  fmt.Sprintf("First%d", i),
				MiddleName: fmt.Sprintf("Middle%d", i),
				LastName:   fmt.Sprintf("Last%d", i),
				Status:     1,
				CreatedAt:  generator.DummyTime(),
				UpdatedAt:  generator.DummyTime(),
			})
	}
	return
}
