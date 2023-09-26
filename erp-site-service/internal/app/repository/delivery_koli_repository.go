package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/model"
)

type IDeliveryKoliRepository interface {
	Get(ctx context.Context, req *dto.DeliveryKoliGetRequest) (res []*model.DeliveryKoli, count int64, err error)
	Create(ctx context.Context, model *model.DeliveryKoli) (err error)
	Delete(ctx context.Context, model *model.DeliveryKoli) (err error)
}

type DeliveryKoliRepository struct {
	opt opt.Options
}

func NewDeliveryKoliRepository() IDeliveryKoliRepository {
	return &DeliveryKoliRepository{
		opt: global.Setup.Common,
	}
}

func (r *DeliveryKoliRepository) Get(ctx context.Context, req *dto.DeliveryKoliGetRequest) (res []*model.DeliveryKoli, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryKoliRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.DeliveryKoli))

	cond := orm.NewCondition()

	cond = cond.And("sales_order_code", req.SopNumber)

	qs = qs.SetCond(cond)

	count, err = qs.AllWithCtx(ctx, &res)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *DeliveryKoliRepository) Create(ctx context.Context, model *model.DeliveryKoli) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryKoliRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, model)
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

func (r *DeliveryKoliRepository) Delete(ctx context.Context, model *model.DeliveryKoli) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryKoliRepository.Delete")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.DeleteWithCtx(ctx, model)
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
