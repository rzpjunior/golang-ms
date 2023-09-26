package repository

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/model"
)

type IPurchaseOrderItemRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID int64, salespersonID int64, submitDateFrom time.Time, submitDateTo time.Time) (purchaseOrderItems []*model.PurchaseOrderItem, count int64, err error)
	GetByPurchaseOrderID(ctx context.Context, purchaseOrderID int64) (purchaseOrderItems []*model.PurchaseOrderItem, count int64, err error)
	Create(ctx context.Context, purchaseOrderItem *model.PurchaseOrderItem) (err error)
	Update(ctx context.Context, purchaseOrderItem *model.PurchaseOrderItem) (err error)
	Delete(ctx context.Context, purchaseOrderID int64) (err error)
}

type PurchaseOrderItemRepository struct {
	opt opt.Options
}

func NewPurchaseOrderItemRepository() IPurchaseOrderItemRepository {
	return &PurchaseOrderItemRepository{
		opt: global.Setup.Common,
	}
}

func (r *PurchaseOrderItemRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID int64, salespersonID int64, submitDateFrom time.Time, submitDateTo time.Time) (purchaseOrderItems []*model.PurchaseOrderItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderItemRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.PurchaseOrderItem))

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

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &purchaseOrderItems)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PurchaseOrderItemRepository) GetByPurchaseOrderID(ctx context.Context, purchaseOrderID int64) (purchaseOrderItems []*model.PurchaseOrderItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderItemRepository.GetByID")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.PurchaseOrderItem))

	cond := orm.NewCondition()

	cond = cond.And("purchase_order_id", purchaseOrderID)

	qs = qs.SetCond(cond)

	count, err = qs.AllWithCtx(ctx, &purchaseOrderItems)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PurchaseOrderItemRepository) Create(ctx context.Context, purchaseOrderItem *model.PurchaseOrderItem) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderItemRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, purchaseOrderItem)
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PurchaseOrderItemRepository) Update(ctx context.Context, purchaseOrderItem *model.PurchaseOrderItem) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderItemRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, purchaseOrderItem)
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PurchaseOrderItemRepository) Delete(ctx context.Context, purchaseOrderID int64) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderItemRepository.Delete")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	purchaseOrderItem := &model.PurchaseOrderItem{
		PurchaseOrderID: purchaseOrderID,
	}

	_, err = tx.DeleteWithCtx(ctx, purchaseOrderItem, "purchase_order_id")
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return
	}
	err = tx.Commit()

	return
}
