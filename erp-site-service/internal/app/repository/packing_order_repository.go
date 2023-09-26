package repository

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/model"
)

type IPackingOrderRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, siteID string, deliveryDateFrom time.Time, deliveryDateTo time.Time) (packingOrders []*model.PackingOrder, count int64, err error)
	GetByID(ctx context.Context, id int64) (packingOrder *model.PackingOrder, err error)
	CheckExisted(ctx context.Context, siteID string, deliveryDate time.Time) (packingOrder *model.PackingOrder, err error)
	Create(ctx context.Context, packingOrder *model.PackingOrder) (err error)
	Update(ctx context.Context, packingOrder *model.PackingOrder, columns ...string) (err error)
}

type PackingOrderRepository struct {
	opt opt.Options
}

func NewPackingOrderRepository() IPackingOrderRepository {
	return &PackingOrderRepository{
		opt: global.Setup.Common,
	}
}

func (r *PackingOrderRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, siteID string, deliveryDateFrom time.Time, deliveryDateTo time.Time) (packingOrders []*model.PackingOrder, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PackingOrderRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.PackingOrder))

	cond := orm.NewCondition()

	if search != "" {
		cond = cond.And("code__icontains", search)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	if siteID != "" {
		cond = cond.And("site_id_gp", siteID)
	}

	if timex.IsValid(deliveryDateFrom) {
		cond = cond.And("delivery_date__gte", timex.ToStartTime(deliveryDateFrom))
	}
	if timex.IsValid(deliveryDateTo) {
		cond = cond.And("delivery_date__lte", timex.ToLastTime(deliveryDateTo))
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &packingOrders)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PackingOrderRepository) GetByID(ctx context.Context, id int64) (packingOrder *model.PackingOrder, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PackingOrderRepository.GetByID")
	defer span.End()

	packingOrder = &model.PackingOrder{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, packingOrder, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PackingOrderRepository) CheckExisted(ctx context.Context, siteID string, deliveryDate time.Time) (packingOrder *model.PackingOrder, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PackingOrderRepository.GetByID")
	defer span.End()

	packingOrder = &model.PackingOrder{
		SiteIDGP:     siteID,
		DeliveryDate: deliveryDate,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, packingOrder, "site_id_gp", "delivery_date")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PackingOrderRepository) Create(ctx context.Context, packingOrder *model.PackingOrder) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PackingOrderRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, packingOrder)
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return
	}

	return
}

func (r *PackingOrderRepository) Update(ctx context.Context, packingOrder *model.PackingOrder, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PackingOrderRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, packingOrder, columns...)

	if err == nil {
		err = tx.Commit()
		if err != nil {
			span.RecordError(err)
			return
		}
	} else {
		span.RecordError(err)
		err = tx.Rollback()
		if err != nil {
			span.RecordError(err)
			return
		}
	}

	return
}
