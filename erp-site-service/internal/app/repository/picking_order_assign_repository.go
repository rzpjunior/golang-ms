package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/model"
)

type IPickingOrderAssignRepository interface {
	SyncGP(ctx context.Context, pickingOrderAssign *model.PickingOrderAssign) (err error)
	Get(ctx context.Context, req *dto.PickingOrderAssignGetRequest) (pickingOrderAssigns []*model.PickingOrderAssign, count int64, err error)
	GetByID(ctx context.Context, id int64, sopNumber string) (pickingOrderAssign *model.PickingOrderAssign, err error)
	Update(ctx context.Context, pickingOrderAssign *model.PickingOrderAssign, columns ...string) (err error)
}

type PickingOrderAssignRepository struct {
	opt opt.Options
}

func NewPickingOrderAssignRepository() IPickingOrderAssignRepository {
	return &PickingOrderAssignRepository{
		opt: global.Setup.Common,
	}
}

func (r *PickingOrderAssignRepository) SyncGP(ctx context.Context, pickingOrderAssign *model.PickingOrderAssign) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PickingOrderAssignRepository.SyncGP")
	defer span.End()

	db := r.opt.Database.Write

	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, _, err = tx.ReadOrCreateWithCtx(ctx, pickingOrderAssign, "SopNumber", []string{}...)
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

func (r *PickingOrderAssignRepository) Get(ctx context.Context, req *dto.PickingOrderAssignGetRequest) (pickingOrderAssigns []*model.PickingOrderAssign, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PickingOrderAssignRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.PickingOrderAssign))

	cond := orm.NewCondition()

	if len(req.Status) > 0 {
		cond = cond.And("status__in", req.Status)
	}

	if len(req.PickingOrderId) > 0 {
		cond = cond.And("picking_order_id__in", req.PickingOrderId)
	}

	if len(req.SopNumber) > 0 {
		cond = cond.And("sop_number__contains", req.SopNumber)
	}

	if len(req.SiteID) > 0 {
		cond = cond.And("site_id__in", req.SiteID)
	}

	if timex.IsValid(req.DeliveryDateFrom) {
		cond = cond.And("delivery_date__gte", timex.ToStartTime(req.DeliveryDateFrom))
	}

	if timex.IsValid(req.DeliveryDateTo) {
		cond = cond.And("delivery_date__lte", timex.ToLastTime(req.DeliveryDateTo))
	}

	if len(req.CheckerId) > 0 {
		cond = cond.And("checker_id_gp__in", req.CheckerId)
	}

	if len(req.WrtId) > 0 {
		cond = cond.And("wrt_id_gp__in", req.WrtId)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &pickingOrderAssigns)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PickingOrderAssignRepository) GetByID(ctx context.Context, id int64, sopNumber string) (pickingOrderAssign *model.PickingOrderAssign, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PickingOrderAssignRepository.GetByID")
	defer span.End()

	pickingOrderAssign = &model.PickingOrderAssign{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		pickingOrderAssign.ID = id
	}

	if sopNumber != "" {
		cols = append(cols, "sop_number")
		pickingOrderAssign.SopNumber = sopNumber
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, pickingOrderAssign, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PickingOrderAssignRepository) Update(ctx context.Context, model *model.PickingOrderAssign, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PickingOrderAssignRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, model, columns...)
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
