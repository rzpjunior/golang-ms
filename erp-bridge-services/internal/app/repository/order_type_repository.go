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

type IOrderTypeRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (orderTypes []*model.OrderType, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (orderType *model.OrderType, err error)
}

type OrderTypeRepository struct {
	opt opt.Options
}

func NewOrderTypeRepository() IOrderTypeRepository {
	return &OrderTypeRepository{
		opt: global.Setup.Common,
	}
}

func (r *OrderTypeRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (orderTypes []*model.OrderType, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "OrderTypeRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.OrderType))

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

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &orderTypes)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *OrderTypeRepository) GetDetail(ctx context.Context, id int64, code string) (orderType *model.OrderType, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "OrderTypeRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	return r.MockDatas(1)[0], nil

	orderType = &model.OrderType{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		orderType.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		orderType.Code = code
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, orderType, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *OrderTypeRepository) MockDatas(total int) (mockDatas []*model.OrderType) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.OrderType{
				ID:          int64(i),
				Code:        fmt.Sprintf("REG%d", i),
				Description: fmt.Sprintf("Dummy OrderType %d", i),
				Status:      1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			})
	}

	return mockDatas
}
