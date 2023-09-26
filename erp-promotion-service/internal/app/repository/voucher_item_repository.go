package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/model"
)

type IVoucherItemRepository interface {
	Get(ctx context.Context, req *dto.VoucherItemRequestGet) (voucherItems []*model.VoucherItem, count int64, err error)
	GetDetail(ctx context.Context, id int64) (voucherItem *model.VoucherItem, err error)
	Create(ctx context.Context, voucherItem *model.VoucherItem) (err error)
}

type VoucherItemRepository struct {
	opt opt.Options
}

func NewVoucherItemRepository() IVoucherItemRepository {
	return &VoucherItemRepository{
		opt: global.Setup.Common,
	}
}

func (r *VoucherItemRepository) Get(ctx context.Context, req *dto.VoucherItemRequestGet) (voucherItems []*model.VoucherItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "VoucherItemRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	cond := orm.NewCondition()

	qs := db.QueryTable(new(model.VoucherItem))

	if req.VoucherID != 0 {
		cond = cond.And("voucher_id", req.VoucherID)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &voucherItems)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *VoucherItemRepository) GetDetail(ctx context.Context, id int64) (voucherItem *model.VoucherItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "VoucherItemRepository.GetDetail")
	defer span.End()

	voucherItem = &model.VoucherItem{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, voucherItem, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *VoucherItemRepository) Create(ctx context.Context, voucherItem *model.VoucherItem) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "VoucherItemRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}
	voucherItem.ID, err = tx.InsertWithCtx(ctx, voucherItem)
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
