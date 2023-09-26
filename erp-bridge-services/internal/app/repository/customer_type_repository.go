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

type ICustomerTypeRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (CustomerTypes []*model.CustomerType, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (CustomerType *model.CustomerType, err error)
}

type CustomerTypeRepository struct {
	opt opt.Options
}

func NewCustomerTypeRepository() ICustomerTypeRepository {
	return &CustomerTypeRepository{
		opt: global.Setup.Common,
	}
}

func (r *CustomerTypeRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (CustomerTypes []*model.CustomerType, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerTypeRepository.Get")
	defer span.End()

	// DUMMY RETURN
	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.CustomerType))

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

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &CustomerTypes)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CustomerTypeRepository) GetDetail(ctx context.Context, id int64, code string) (CustomerType *model.CustomerType, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerTypeRepository.GetDetail")
	defer span.End()

	// DUMMY RETURN
	return r.MockDatas(1)[0], nil

	CustomerType = &model.CustomerType{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		CustomerType.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		CustomerType.Code = code
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, CustomerType, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CustomerTypeRepository) MockDatas(total int) (mockDatas []*model.CustomerType) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.CustomerType{
				ID:           int64(i),
				Code:         fmt.Sprintf("CSTY%d", i),
				Description:  fmt.Sprintf("Dummy CustomerType %d", i),
				GroupType:    "Dummy GroupType",
				Abbreviation: "Dummy Abbreviation",
				Status:       1,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			})
	}
	return
}
