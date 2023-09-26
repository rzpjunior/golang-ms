package repository

import (
	"context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/generator"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IDeliveryOrderRepository interface {
	GetDetail(ctx context.Context, req *bridge_service.GetDeliveryOrderDetailRequest) (DeliveryOrder *model.DeliveryOrder, err error)
}

type DeliveryOrderRepository struct {
	opt opt.Options
}

func NewDeliveryOrderRepository() IDeliveryOrderRepository {
	return &DeliveryOrderRepository{
		opt: global.Setup.Common,
	}
}

func (r *DeliveryOrderRepository) GetDetail(ctx context.Context, req *bridge_service.GetDeliveryOrderDetailRequest) (DeliveryOrder *model.DeliveryOrder, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryOrderRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	return r.MockDatas(1)[0], nil

	return
}

func (r *DeliveryOrderRepository) MockDatas(total int) (mockDatas []*model.DeliveryOrder) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.DeliveryOrder{
				ID:          1,
				Code:        fmt.Sprintf("DUMMY-DO%d", i),
				CustomerID:  1,
				WrtID:       1,
				SiteID:      1,
				Status:      1,
				CreatedDate: generator.DummyTime(),
			})
	}
	return mockDatas
}
