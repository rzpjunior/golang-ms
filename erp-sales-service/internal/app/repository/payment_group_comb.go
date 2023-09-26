package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/model"
)

type IPaymentGroupCombRepository interface {
	GetPaymentGroupCombList(ctx context.Context, req *pb.GetPaymentGroupCombListRequest) (paymentGroupComb []*model.PaymentGroupComb, err error)
}

type PaymentGroupCombRepository struct {
	opt opt.Options
}

func NewPaymentGroupCombRepository() IPaymentGroupCombRepository {
	return &PaymentGroupCombRepository{
		opt: global.Setup.Common,
	}
}

func (r *PaymentGroupCombRepository) GetPaymentGroupCombList(ctx context.Context, req *pb.GetPaymentGroupCombListRequest) (paymentGroupComb []*model.PaymentGroupComb, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "GetPaymentGroupCombList.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.PaymentGroupComb))

	cond := orm.NewCondition()

	if req.PaymentGroupSls != "" {
		cond = cond.And("payment_group_sls", req.PaymentGroupSls)
	}

	if req.TermPaymentSls != "" {
		cond = cond.And("term_payment_sls", req.TermPaymentSls)
	}
	qs = qs.SetCond(cond)

	qs = qs.OrderBy("id")
	_, err = qs.AllWithCtx(ctx, &paymentGroupComb)

	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
