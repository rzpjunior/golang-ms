package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/model"
)

type IPriceTieringLogRepository interface {
	Get(ctx context.Context, req *dto.PriceTieringLogRequestGet) (PriceTieringLogs []*model.PriceTieringLog, err error)
	GetDetail(ctx context.Context, id int64) (PriceTieringLog *model.PriceTieringLog, err error)
	Create(ctx context.Context, PriceTieringLog *model.PriceTieringLog) (err error)
	Update(ctx context.Context, PriceTieringLog *model.PriceTieringLog, columns ...string) (err error)
}

type PriceTieringLogRepository struct {
	opt opt.Options
}

func NewPriceTieringLogRepository() IPriceTieringLogRepository {
	return &PriceTieringLogRepository{
		opt: global.Setup.Common,
	}
}

func (r *PriceTieringLogRepository) Get(ctx context.Context, req *dto.PriceTieringLogRequestGet) (PriceTieringLogs []*model.PriceTieringLog, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PriceTieringLogRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	cond := orm.NewCondition()

	qs := db.QueryTable(new(model.PriceTieringLog))

	if req.PriceTieringIDGP != "" {
		cond = cond.And("price_tiering_id_gp", req.PriceTieringIDGP)
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

	if req.ItemID != 0 {
		cond = cond.And("item_id", req.ItemID)
	}

	if req.Status != 0 {
		cond = cond.And("status", req.Status)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	_, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &PriceTieringLogs)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PriceTieringLogRepository) GetDetail(ctx context.Context, id int64) (PriceTieringLog *model.PriceTieringLog, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PriceTieringLogRepository.GetDetail")
	defer span.End()

	PriceTieringLog = &model.PriceTieringLog{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, PriceTieringLog, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PriceTieringLogRepository) Create(ctx context.Context, PriceTieringLog *model.PriceTieringLog) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PriceTieringLogRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}
	PriceTieringLog.ID, err = tx.InsertWithCtx(ctx, PriceTieringLog)
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

func (r *PriceTieringLogRepository) Update(ctx context.Context, PriceTieringLog *model.PriceTieringLog, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PriceTieringLogRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, PriceTieringLog, columns...)

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
