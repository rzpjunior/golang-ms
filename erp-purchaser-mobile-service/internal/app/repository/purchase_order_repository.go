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

type IPurchaseOrderRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, recognitionDateFrom time.Time, recognitionDateTo time.Time, code []string, isNotConsolidated bool) (purchaseOrders []*model.PurchaseOrder, count int64, err error)
	GetByID(ctx context.Context, id string) (purchaseOrder *model.PurchaseOrder, err error)
	GetByConsolidatedShipmentID(ctx context.Context, consolidatedShipmentID int64) (purchaseOrder []*model.PurchaseOrder, count int64, err error)
	Create(ctx context.Context, purchaseOrder *model.PurchaseOrder) (err error)
	Update(ctx context.Context, purchaseOrder *model.PurchaseOrder, columns ...string) (err error)
	DeleteConsolidatedShipmentID(ctx context.Context, consolidatedShipmentID int64) (err error)
}

type PurchaseOrderRepository struct {
	opt opt.Options
}

func NewPurchaseOrderRepository() IPurchaseOrderRepository {
	return &PurchaseOrderRepository{
		opt: global.Setup.Common,
	}
}

func (r *PurchaseOrderRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, recognitionDateFrom time.Time, recognitionDateTo time.Time, code []string, isNotConsolidated bool) (purchaseOrders []*model.PurchaseOrder, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.PurchaseOrder))

	cond := orm.NewCondition()

	if search != "" {
		cond = cond.And("code__icontains", search)
	}

	if len(code) != 0 {
		cond = cond.And("code__in", code)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	if isNotConsolidated {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("consolidated_shipment_id", 0).Or("consolidated_shipment_id__isnull", true)
		cond = cond.AndCond(condGroup)
	}

	if timex.IsValid(recognitionDateFrom) {
		cond = cond.And("created_at__gte", timex.ToStartTime(recognitionDateFrom))
	}
	if timex.IsValid(recognitionDateTo) {
		cond = cond.And("created_at__lte", timex.ToLastTime(recognitionDateTo))
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	if offset != 0 {
		qs = qs.Offset(offset)
	}

	if limit != 0 {
		qs = qs.Limit(limit)
	}

	count, err = qs.AllWithCtx(ctx, &purchaseOrders)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PurchaseOrderRepository) GetByConsolidatedShipmentID(ctx context.Context, consolidatedShipmentID int64) (purchaseOrders []*model.PurchaseOrder, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderRepository.GetByConsolidatedShipmentID")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.PurchaseOrder))

	cond := orm.NewCondition()

	if consolidatedShipmentID != 0 {
		cond = cond.And("consolidated_shipment_id", consolidatedShipmentID)
	}

	qs = qs.SetCond(cond)

	count, err = qs.AllWithCtx(ctx, &purchaseOrders)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PurchaseOrderRepository) GetByID(ctx context.Context, id string) (purchaseOrder *model.PurchaseOrder, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderRepository.GetByID")
	defer span.End()

	purchaseOrder = &model.PurchaseOrder{
		PurchaseOrderIDGP: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, purchaseOrder, "purchase_order_id_gp")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PurchaseOrderRepository) Create(ctx context.Context, purchaseOrder *model.PurchaseOrder) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, purchaseOrder)
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

func (r *PurchaseOrderRepository) Update(ctx context.Context, purchaseOrder *model.PurchaseOrder, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, purchaseOrder, columns...)
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

func (r *PurchaseOrderRepository) DeleteConsolidatedShipmentID(ctx context.Context, consolidatedShipmentID int64) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderRepository.DeleteConsolidatedShipmentID")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	qs := db.QueryTable(new(model.PurchaseOrder))
	_, err = qs.Filter("consolidated_shipment_id", consolidatedShipmentID).UpdateWithCtx(ctx, orm.Params{"consolidated_shipment_id": nil})
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
