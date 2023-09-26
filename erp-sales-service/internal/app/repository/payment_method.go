package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/model"
)

type IPaymentMethodRepository interface {
	GetPaymentMethodList(ctx context.Context, req *pb.GetPaymentMethodListRequest) (paymentMethod []*model.PaymentMethod, err error)
}

type PaymentMethodRepository struct {
	opt opt.Options
}

func NewPaymentMethodRepository() IPaymentMethodRepository {
	return &PaymentMethodRepository{
		opt: global.Setup.Common,
	}
}

func (r *PaymentMethodRepository) GetPaymentMethodList(ctx context.Context, req *pb.GetPaymentMethodListRequest) (paymentMethod []*model.PaymentMethod, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "GetPaymentMethodList.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.PaymentMethod))

	cond := orm.NewCondition()

	if req.Status != 0 {
		cond = cond.And("status", req.Status)
	}

	if req.Maintenance != 0 {
		cond = cond.And("maintenance", req.Status)
	}

	if req.Publish != 0 {
		cond = cond.And("publish", req.Status)
	}

	if req.Search != "" {
		cond = cond.And("name__icontains", req.Search)
	}

	if req.Id != "" {
		cond = cond.And("code", req.Id)
	}
	qs = qs.SetCond(cond)

	qs = qs.OrderBy("id")
	_, err = qs.AllWithCtx(ctx, &paymentMethod)

	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
