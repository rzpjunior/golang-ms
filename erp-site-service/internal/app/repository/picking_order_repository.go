package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/model"
)

type IPickingOrderRepository interface {
	SyncGP(ctx context.Context, pickingOrder *model.PickingOrder) (err error)
	GetDetail(ctx context.Context, id int64, docNumber string) (pickingOrder *model.PickingOrder, err error)
	Update(ctx context.Context, pickingOrder *model.PickingOrder, columns ...string) (err error)
}

type PickingOrderRepository struct {
	opt opt.Options
}

func NewPickingOrderRepository() IPickingOrderRepository {
	return &PickingOrderRepository{
		opt: global.Setup.Common,
	}
}

func (r *PickingOrderRepository) SyncGP(ctx context.Context, pickingOrder *model.PickingOrder) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PickingOrderRepository.SyncGP")
	defer span.End()

	db := r.opt.Database.Write

	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, _, err = tx.ReadOrCreateWithCtx(ctx, pickingOrder, "DocNumber", []string{}...)
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

func (r *PickingOrderRepository) GetDetail(ctx context.Context, id int64, docNumber string) (pickingOrder *model.PickingOrder, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PickingOrderRepository.GetDetail")
	defer span.End()

	pickingOrder = &model.PickingOrder{}

	var cols []string

	if id != 0 {
		cols = append(cols, "id")
		pickingOrder.Id = id
	}

	if docNumber != "" {
		cols = append(cols, "doc_number")
		pickingOrder.DocNumber = docNumber
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, pickingOrder, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PickingOrderRepository) Update(ctx context.Context, model *model.PickingOrder, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PickingOrderRepository.Update")
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
