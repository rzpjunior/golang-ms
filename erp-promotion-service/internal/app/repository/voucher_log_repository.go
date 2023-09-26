package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/model"
)

type IVoucherLogRepository interface {
	Get(ctx context.Context, req *dto.VoucherLogRequestGet) (voucherLogs []*model.VoucherLog, count int64, err error)
	GetDetail(ctx context.Context, id int64) (voucherLog *model.VoucherLog, err error)
	Create(ctx context.Context, voucherLog *model.VoucherLog) (err error)
	Update(ctx context.Context, VoucherLog *model.VoucherLog, columns ...string) (err error)
}

type VoucherLogRepository struct {
	opt opt.Options
}

func NewVoucherLogRepository() IVoucherLogRepository {
	return &VoucherLogRepository{
		opt: global.Setup.Common,
	}
}

func (r *VoucherLogRepository) Get(ctx context.Context, req *dto.VoucherLogRequestGet) (voucherLogs []*model.VoucherLog, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "VoucherLogRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	cond := orm.NewCondition()

	qs := db.QueryTable(new(model.VoucherLog))

	if req.VoucherID != 0 {
		cond = cond.And("voucher_id", req.VoucherID)
	}

	if req.SalesOrderIDGP != "" {
		cond = cond.And("sales_order_id_gp", req.SalesOrderIDGP)
	}

	if req.CustomerID != 0 {
		cond = cond.And("customer_id", req.CustomerID)
	}

	if req.AddressIDGP != "" {
		cond = cond.And("address_id_gp", req.AddressIDGP)
	}

	if req.Status != 0 {
		cond = cond.And("status", req.Status)
	}

	if req.Code != "" {
		cond = cond.And("code", req.Code)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &voucherLogs)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *VoucherLogRepository) GetDetail(ctx context.Context, id int64) (voucherLog *model.VoucherLog, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "VoucherLogRepository.GetDetail")
	defer span.End()

	voucherLog = &model.VoucherLog{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, voucherLog, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *VoucherLogRepository) Create(ctx context.Context, voucherLog *model.VoucherLog) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "VoucherLogRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}
	voucherLog.ID, err = tx.InsertWithCtx(ctx, voucherLog)
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

func (r *VoucherLogRepository) Update(ctx context.Context, VoucherLog *model.VoucherLog, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "VoucherLogRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, VoucherLog, columns...)

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
