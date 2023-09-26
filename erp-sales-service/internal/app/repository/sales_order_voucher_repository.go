package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/model"
)

type ISalesOrderVoucherRepository interface {
	GetList(ctx context.Context, req *dto.GetSalesOrderVoucherListRequest) (SalesOrderVoucheres []*model.SalesOrderVoucher, count int64, err error)
}

type SalesOrderVoucherRepository struct {
	opt opt.Options
}

func NewSalesOrderVoucherRepository() ISalesOrderVoucherRepository {
	return &SalesOrderVoucherRepository{
		opt: global.Setup.Common,
	}
}

func (r *SalesOrderVoucherRepository) GetList(ctx context.Context, req *dto.GetSalesOrderVoucherListRequest) (SalesOrderVouchers []*model.SalesOrderVoucher, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesOrderVoucherRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.SalesOrderVoucher))

	cond := orm.NewCondition()

	if req.SalesOrderID != 0 {
		cond = cond.And("sales_order_id", req.SalesOrderID)
	}

	qs = qs.SetCond(cond)

	count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &SalesOrderVouchers)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
