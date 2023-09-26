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

type IConsolidatedShipmentRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, siteCode string, createdAtFrom time.Time, createdAtTo time.Time, createdBy int64) (consolidatedShipments []*model.ConsolidatedShipment, count int64, err error)
	GetByID(ctx context.Context, id int64) (consolidatedShipment *model.ConsolidatedShipment, err error)
	GetByPurchaseOrderID(ctx context.Context, id string) (consolidatedShipment *model.ConsolidatedShipment, err error)
	Create(ctx context.Context, consolidatedShipment *model.ConsolidatedShipment) (err error)
	Update(ctx context.Context, consolidatedShipment *model.ConsolidatedShipment, columns ...string) (err error)
}

type ConsolidatedShipmentRepository struct {
	opt opt.Options
}

func NewConsolidatedShipmentRepository() IConsolidatedShipmentRepository {
	return &ConsolidatedShipmentRepository{
		opt: global.Setup.Common,
	}
}

func (r *ConsolidatedShipmentRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, siteCode string, createdAtFrom time.Time, createdAtTo time.Time, createdBy int64) (consolidatedShipments []*model.ConsolidatedShipment, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ConsolidatedShipmentRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.ConsolidatedShipment))

	cond := orm.NewCondition()

	if search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("driver_name__icontains", search).Or("vehicle_number__icontains", search).Or("code__icontains", search)
		cond = cond.AndCond(condGroup)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	if siteCode != "" {
		cond = cond.And("site_name", siteCode)
	}

	if timex.IsValid(createdAtFrom) {
		cond = cond.And("created_at__gte", timex.ToStartTime(createdAtFrom))
	}
	if timex.IsValid(createdAtTo) {
		cond = cond.And("created_at__lte", timex.ToLastTime(createdAtTo))
	}

	if createdBy != 0 {
		cond = cond.And("created_by", createdBy)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &consolidatedShipments)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ConsolidatedShipmentRepository) GetByID(ctx context.Context, id int64) (consolidatedShipment *model.ConsolidatedShipment, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ConsolidatedShipmentRepository.GetByID")
	defer span.End()

	consolidatedShipment = &model.ConsolidatedShipment{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, consolidatedShipment, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ConsolidatedShipmentRepository) GetByPurchaseOrderID(ctx context.Context, code string) (consolidatedShipment *model.ConsolidatedShipment, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ConsolidatedShipmentRepository.GetByID")
	defer span.End()

	consolidatedShipment = &model.ConsolidatedShipment{
		Code: code,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, consolidatedShipment, "code")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ConsolidatedShipmentRepository) Create(ctx context.Context, consolidatedShipment *model.ConsolidatedShipment) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ConsolidatedShipmentRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, consolidatedShipment)
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

func (r *ConsolidatedShipmentRepository) Update(ctx context.Context, consolidatedShipment *model.ConsolidatedShipment, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ConsolidatedShipmentRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, consolidatedShipment, columns...)
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
