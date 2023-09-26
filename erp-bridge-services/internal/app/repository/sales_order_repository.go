package repository

import (
	"context"
	"fmt"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/generator"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type ISalesOrderRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, addressID int64, customerID int64, salespersonID int64, orderDateFrom time.Time, orderDateTo time.Time) (salesOrders []*model.SalesOrder, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (salesOrder *model.SalesOrder, err error)
}

type SalesOrderRepository struct {
	opt opt.Options
}

func NewSalesOrderRepository() ISalesOrderRepository {
	return &SalesOrderRepository{
		opt: global.Setup.Common,
	}
}

func (r *SalesOrderRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, addressID int64, customerID int64, salespersonID int64, orderDateFrom time.Time, orderDateTo time.Time) (salesOrders []*model.SalesOrder, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesOrderRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.SalesOrder))

	cond := orm.NewCondition()

	if search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("code__icontains", search)
		cond = cond.AndCond(condGroup)
	}

	if addressID != 0 {
		cond = cond.And("address_id", addressID)
	}

	if customerID != 0 {
		cond = cond.And("customer_id", customerID)
	}

	if salespersonID != 0 {
		cond = cond.And("salesperson_id", salespersonID)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	if timex.IsValid(orderDateFrom) {
		cond = cond.And("order_date__gte", timex.ToStartTime(orderDateFrom))
	}
	if timex.IsValid(orderDateTo) {
		cond = cond.And("order_date__lte", timex.ToLastTime(orderDateTo))
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &salesOrders)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesOrderRepository) GetDetail(ctx context.Context, id int64, code string) (salesOrder *model.SalesOrder, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesOrderRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	return r.MockDatas(1)[0], nil

	salesOrder = &model.SalesOrder{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		salesOrder.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		salesOrder.Code = code
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, salesOrder, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesOrderRepository) MockDatas(total int) (mockDatas []*model.SalesOrder) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.SalesOrder{
				ID:            int64(i),
				Code:          fmt.Sprintf("DUMMY-SO%d", i),
				DocNumber:     fmt.Sprintf("DUMMY-SODOC%d", i),
				AddressID:     1,
				CustomerID:    1,
				SalespersonID: 1,
				WrtID:         1,
				OrderTypeID:   1,
				SiteID:        1,
				Application:   1,
				Status:        1,
				OrderDate:     generator.DummyTime(),
				Total:         999999,
				CreatedDate:   generator.DummyTime(),
				ModifiedDate:  generator.DummyTime(),
				FinishedDate:  generator.DummyTime(),
				CreatedAt:     generator.DummyTime(),
				UpdatedAt:     generator.DummyTime(),
			})
	}
	return mockDatas
}
