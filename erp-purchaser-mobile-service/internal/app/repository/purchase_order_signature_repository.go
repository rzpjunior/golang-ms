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

type IPurchaseOrderSignatureRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, recognitionDateFrom time.Time, recognitionDateTo time.Time) (purchaseOrderSignatures []*model.PurchaseOrderSignature, count int64, err error)
	GetByID(ctx context.Context, id int64) (purchaseOrderSignature *model.PurchaseOrderSignature, err error)
	GetSignatureByPurchaseOrderID(ctx context.Context, id string) (purchaseOrderSignature []*model.PurchaseOrderSignature, count int64, err error)
	Create(ctx context.Context, purchaseOrderSignature *model.PurchaseOrderSignature) (err error)
	Update(ctx context.Context, purchaseOrderSignature *model.PurchaseOrderSignature, columns ...string) (err error)
}

type PurchaseOrderSignatureRepository struct {
	opt opt.Options
}

func NewPurchaseOrderSignatureRepository() IPurchaseOrderSignatureRepository {
	return &PurchaseOrderSignatureRepository{
		opt: global.Setup.Common,
	}
}

func (r *PurchaseOrderSignatureRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, recognitionDateFrom time.Time, recognitionDateTo time.Time) (purchaseOrderSignatures []*model.PurchaseOrderSignature, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderSignatureRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.PurchaseOrderSignature))

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

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &purchaseOrderSignatures)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PurchaseOrderSignatureRepository) GetSignatureByPurchaseOrderID(ctx context.Context, id string) (purchaseOrderSignatures []*model.PurchaseOrderSignature, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderSignatureRepository.GetSignatureByPurchaseOrderID")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.PurchaseOrderSignature))

	cond := orm.NewCondition()

	if id != "" {
		cond = cond.And("purchase_order_id_gp", id)
	}

	qs = qs.SetCond(cond)

	count, err = qs.AllWithCtx(ctx, &purchaseOrderSignatures)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PurchaseOrderSignatureRepository) GetByID(ctx context.Context, id int64) (purchaseOrderSignature *model.PurchaseOrderSignature, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderSignatureRepository.GetByID")
	defer span.End()

	purchaseOrderSignature = &model.PurchaseOrderSignature{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, purchaseOrderSignature, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PurchaseOrderSignatureRepository) Create(ctx context.Context, purchaseOrderSignature *model.PurchaseOrderSignature) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderSignatureRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, purchaseOrderSignature)
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

func (r *PurchaseOrderSignatureRepository) Update(ctx context.Context, purchaseOrderSignature *model.PurchaseOrderSignature, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderSignatureRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, purchaseOrderSignature, columns...)
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
