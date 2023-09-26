package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/model"
)

type IPaymentChannelRepository interface {
	GetPaymentChannelList(ctx context.Context, req *pb.GetPaymentChannelListRequest) (paymentChannel []*model.PaymentChannel, err error)
}

type PaymentChannelRepository struct {
	opt opt.Options
}

func NewPaymentChannelRepository() IPaymentChannelRepository {
	return &PaymentChannelRepository{
		opt: global.Setup.Common,
	}
}

func (r *PaymentChannelRepository) GetPaymentChannelList(ctx context.Context, req *pb.GetPaymentChannelListRequest) (paymentChannel []*model.PaymentChannel, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "GetPaymentChannelList.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.PaymentChannel))

	cond := orm.NewCondition()

	if req.Status != 0 {
		cond = cond.And("status", req.Status)
	}

	if req.PublishFva != 0 {
		cond = cond.And("publish_fva", req.PublishFva)
	}

	if req.PublishIva != 0 {
		cond = cond.And("publish_iva", req.PublishIva)
	}

	if req.Value != "" {
		cond = cond.And("value", req.Value)
	}

	qs = qs.SetCond(cond)

	qs = qs.OrderBy("id")
	_, err = qs.AllWithCtx(ctx, &paymentChannel)

	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
