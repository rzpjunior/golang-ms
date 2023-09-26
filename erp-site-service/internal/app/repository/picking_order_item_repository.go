package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/model"
)

type IPickingOrderItemRepository interface {
	SyncGP(ctx context.Context, pickingOrderItem *model.PickingOrderItem) (err error)
	Get(ctx context.Context, req *dto.PickingOrderItemGetRequest) (pickingOrderItems []*model.PickingOrderItem, count int64, err error)
	GetByID(ctx context.Context, req *dto.PickingOrderItemGetDetailRequest) (pickingOrderItem *model.PickingOrderItem, err error)
	Update(ctx context.Context, pickingOrderItem *model.PickingOrderItem, columns ...string) (err error)
}

type PickingOrderItemRepository struct {
	opt opt.Options
}

func NewPickingOrderItemRepository() IPickingOrderItemRepository {
	return &PickingOrderItemRepository{
		opt: global.Setup.Common,
	}
}

func (r *PickingOrderItemRepository) SyncGP(ctx context.Context, pickingOrderItem *model.PickingOrderItem) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PickingOrderItem.SyncGP")
	defer span.End()

	db := r.opt.Database.Write

	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, _, err = tx.ReadOrCreateWithCtx(ctx, pickingOrderItem, "PickingOrderAssignID", "ItemNumber")
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

func (r *PickingOrderItemRepository) Get(ctx context.Context, req *dto.PickingOrderItemGetRequest) (pickingOrderItems []*model.PickingOrderItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PickingOrderItemRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.PickingOrderItem))

	cond := orm.NewCondition()

	if len(req.Status) > 0 {
		cond = cond.And("status__in", req.Status)
	}

	if len(req.PickingOrderAssignId) > 0 {
		cond = cond.And("picking_order_assign_id__in", req.PickingOrderAssignId)
	}

	if len(req.ItemNumber) > 0 {
		cond = cond.And("item_number__in", req.ItemNumber)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	count, err = qs.AllWithCtx(ctx, &pickingOrderItems)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PickingOrderItemRepository) GetByID(ctx context.Context, req *dto.PickingOrderItemGetDetailRequest) (pickingOrderItem *model.PickingOrderItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PickingOrderItemRepository.GetByID")
	defer span.End()

	pickingOrderItem = &model.PickingOrderItem{}

	var cols []string

	if req.Id != 0 {
		cols = append(cols, "id")
		pickingOrderItem.Id = req.Id
	}

	if req.PickingOrderAssignId != 0 {
		cols = append(cols, "picking_order_assign_id")
		pickingOrderItem.PickingOrderAssignId = req.PickingOrderAssignId
	}

	if req.ItemNumber != "" {
		cols = append(cols, "item_number")
		pickingOrderItem.ItemNumber = req.ItemNumber
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, pickingOrderItem, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PickingOrderItemRepository) Update(ctx context.Context, model *model.PickingOrderItem, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PickingOrderItemRepository.Update")
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
