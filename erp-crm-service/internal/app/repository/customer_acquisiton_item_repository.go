package repository

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/model"
)

type ICustomerAcquisitionItemRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID int, salespersonID int, submitDateFrom time.Time, submitDateTo time.Time) (customerAcquisitionItems []*model.CustomerAcquisitionItem, count int64, err error)
	GetByCustomerAcquisitionID(ctx context.Context, customerAcquisitionID int64) (customerAcquisitionItems []*model.CustomerAcquisitionItem, count int64, err error)
	GetByID(ctx context.Context, id int64) (customerAcquisition *model.CustomerAcquisitionItem, err error)
}

type CustomerAcquisitionItemRepository struct {
	opt opt.Options
}

func NewCustomerAcquisitionItemRepository() ICustomerAcquisitionItemRepository {
	return &CustomerAcquisitionItemRepository{
		opt: global.Setup.Common,
	}
}

func (r *CustomerAcquisitionItemRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID int, salespersonID int, submitDateFrom time.Time, submitDateTo time.Time) (customerAcquisitionItems []*model.CustomerAcquisitionItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerAcquisitionItemRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.CustomerAcquisitionItem))

	cond := orm.NewCondition()

	if search != "" {
		cond = cond.And("name__icontains", search)
	}

	if territoryID != 0 {
		cond = cond.And("territory_id", territoryID)
	}

	if salespersonID != 0 {
		cond = cond.And("salesperson_id", salespersonID)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	if timex.IsValid(submitDateFrom) {
		cond = cond.And("submit_date__gte", timex.ToStartTime(submitDateFrom))
	}
	if timex.IsValid(submitDateTo) {
		cond = cond.And("submit_date__lte", timex.ToLastTime(submitDateTo))
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &customerAcquisitionItems)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CustomerAcquisitionItemRepository) GetByCustomerAcquisitionID(ctx context.Context, customerAcquisitionID int64) (customerAcquisitionItems []*model.CustomerAcquisitionItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerAcquisitionItemRepository.GetByID")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.CustomerAcquisitionItem))

	count, err = qs.Filter("customer_acquisition_id", customerAcquisitionID).AllWithCtx(ctx, &customerAcquisitionItems)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CustomerAcquisitionItemRepository) GetByID(ctx context.Context, id int64) (customerAcquisition *model.CustomerAcquisitionItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerAcquisitionItemRepository.GetByID")
	defer span.End()

	customerAcquisition = &model.CustomerAcquisitionItem{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, customerAcquisition, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
