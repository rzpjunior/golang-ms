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

type IPurchaseOrderImageRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, recognitionDateFrom time.Time, recognitionDateTo time.Time) (purchaseOrderImages []*model.PurchaseOrderImage, count int64, err error)
	GetByID(ctx context.Context, id int64) (purchaseOrderImage *model.PurchaseOrderImage, err error)
	GetImageByPurchaseOrderID(ctx context.Context, id string) (purchaseOrderImage []*model.PurchaseOrderImage, count int64, err error)
	Create(ctx context.Context, purchaseOrderImage *model.PurchaseOrderImage) (err error)
	Delete(ctx context.Context, purchaseOrderIDGP string) (err error)
}

type PurchaseOrderImageRepository struct {
	opt opt.Options
}

func NewPurchaseOrderImageRepository() IPurchaseOrderImageRepository {
	return &PurchaseOrderImageRepository{
		opt: global.Setup.Common,
	}
}

func (r *PurchaseOrderImageRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, recognitionDateFrom time.Time, recognitionDateTo time.Time) (purchaseOrderImages []*model.PurchaseOrderImage, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderImageRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.PurchaseOrderImage))

	cond := orm.NewCondition()

	if search != "" {
		cond = cond.And("code__icontains", search)
	}

	if status != 0 {
		cond = cond.And("status", status)
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

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &purchaseOrderImages)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PurchaseOrderImageRepository) GetImageByPurchaseOrderID(ctx context.Context, id string) (purchaseOrderImages []*model.PurchaseOrderImage, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderImageRepository.GetImageByPurchaseOrderID")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.PurchaseOrderImage))

	cond := orm.NewCondition()

	if id != "" {
		cond = cond.And("purchase_order_id_gp", id)
	}

	qs = qs.SetCond(cond)

	count, err = qs.AllWithCtx(ctx, &purchaseOrderImages)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PurchaseOrderImageRepository) GetByID(ctx context.Context, id int64) (purchaseOrderImage *model.PurchaseOrderImage, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderImageRepository.GetByID")
	defer span.End()

	purchaseOrderImage = &model.PurchaseOrderImage{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, purchaseOrderImage, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PurchaseOrderImageRepository) Create(ctx context.Context, purchaseOrderImage *model.PurchaseOrderImage) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderImageRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, purchaseOrderImage)
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

func (r *PurchaseOrderImageRepository) Delete(ctx context.Context, purchaseOrderIDGP string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderImageRepository.Delete")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	purchaseOrderImage := &model.PurchaseOrderImage{
		PurchaseOrderIDGP: purchaseOrderIDGP,
	}

	_, err = tx.DeleteWithCtx(ctx, purchaseOrderImage, "purchase_order_id_gp")
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
