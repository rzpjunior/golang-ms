package repository

import (
	"context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"

	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IDeliveryFeeRepository interface {
	Get(ctx context.Context, req *pb.GetDeliveryFeeListRequest) (DeliveryFeees []*model.DeliveryFee, count int64, err error)
	GetDetail(ctx context.Context, req *pb.GetDeliveryFeeDetailRequest) (DeliveryFee *model.DeliveryFee, err error)
}

type DeliveryFeeRepository struct {
	opt opt.Options
}

func NewDeliveryFeeRepository() IDeliveryFeeRepository {
	return &DeliveryFeeRepository{
		opt: global.Setup.Common,
	}
}

func (r *DeliveryFeeRepository) Get(ctx context.Context, req *pb.GetDeliveryFeeListRequest) (DeliveryFeees []*model.DeliveryFee, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryFeeRepository.Get")
	defer span.End()

	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	// db := r.opt.Database.Read

	// qs := db.QueryTable(new(model.DeliveryFee))

	// cond := orm.NewCondition()

	// if req.Search != "" {
	// 	condGroup := orm.NewCondition()
	// 	condGroup = condGroup.And("description__icontains", req.Search).Or("code__icontains", req.Search)
	// 	cond = cond.AndCond(condGroup)
	// }

	// if len(req.Status) != 0 {
	// 	cond = cond.And("status__in", req.Status)
	// }

	// if req.CustomerTypeId != 0 {
	// 	cond = cond.And("customer_type_id__in", req.CustomerTypeId)
	// }

	// if req.RegionId != 0 {
	// 	cond = cond.And("region_id__in", req.RegionId)
	// }
	// qs = qs.SetCond(cond)

	// if req.OrderBy != "" {
	// 	qs = qs.OrderBy(req.OrderBy)
	// }

	// count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &DeliveryFeees)
	// if err != nil {
	// 	span.RecordError(err)
	// 	return
	// }

	return
}

func (r *DeliveryFeeRepository) GetDetail(ctx context.Context, req *pb.GetDeliveryFeeDetailRequest) (DeliveryFee *model.DeliveryFee, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryFeeRepository.GetDetail")
	defer span.End()

	return r.MockDatas(1)[0], nil

	// DeliveryFee = &model.DeliveryFee{}

	// var cols []string

	// if req.Id != 0 {
	// 	cols = append(cols, "id")
	// 	DeliveryFee.ID = req.Id
	// }

	// if req.Code != "" {
	// 	cols = append(cols, "code")
	// 	DeliveryFee.Code = req.Code
	// }
	// if req.CustomerTypeId != "" {
	// 	cols = append(cols, "customer_type_id")
	// 	DeliveryFee.CutomerTypeId = req.CustomerTypeId
	// }

	// if req.RegionId != 0 {
	// 	cols = append(cols, "region_id")
	// 	DeliveryFee.RegionId = req.RegionId
	// }

	// db := r.opt.Database.Read
	// err = db.ReadWithCtx(ctx, DeliveryFee, cols...)
	// if err != nil {
	// 	span.RecordError(err)
	// 	return
	// }

	return
}

func (r *DeliveryFeeRepository) MockDatas(total int) (mockDatas []*model.DeliveryFee) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.DeliveryFee{
				ID:   int64(i),
				Code: fmt.Sprintf("DummyDF%d", i),
				Name: fmt.Sprintf("Dummy Delivery Fee%d", i),
				//Note:          "",
				Status:        1,
				MinimumOrder:  150000,
				DeliveryFee:   10000,
				RegionId:      1,
				CutomerTypeId: "",
			})
	}
	return
}
