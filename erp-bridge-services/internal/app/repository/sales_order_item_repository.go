package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/generator"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type ISalesOrderItemRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, salesOrderID int64, itemID int64) (salesOrderItems []*model.SalesOrderItem, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (salesOrderItem *model.SalesOrderItem, err error)
}

type SalesOrderItemRepository struct {
	opt opt.Options
}

func NewSalesOrderItemRepository() ISalesOrderItemRepository {
	return &SalesOrderItemRepository{
		opt: global.Setup.Common,
	}
}

func (r *SalesOrderItemRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, salesOrderID int64, itemID int64) (salesOrderItems []*model.SalesOrderItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesOrderItemRepository.Get")
	defer span.End()

	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil
}

func (r *SalesOrderItemRepository) GetDetail(ctx context.Context, id int64, code string) (salesOrderItem *model.SalesOrderItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesOrderItemRepository.GetDetail")
	defer span.End()

	return r.MockDatas(1)[0], nil
}

func (r *SalesOrderItemRepository) MockDatas(total int) (mockDatas []*model.SalesOrderItem) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.SalesOrderItem{
				ID:            int64(i),
				SalesOrderId:  1,
				ItemId:        generator.DummyInt64(1, 10),
				OrderQty:      0.25 * float64(generator.DummyInt(1, 50)),
				DefaultPrice:  generator.DummyFloat64(10000, 100000),
				UnitPrice:     generator.DummyFloat64(10000, 100000),
				TaxableItem:   1,
				TaxPercentage: generator.DummyFloat64(1, 10),
				ShadowPrice:   generator.DummyFloat64(10000, 100000),
				Subtotal:      generator.DummyFloat64(10000, 100000),
				Weight:        0.25 * float64(generator.DummyInt(1, 5)),
				Note:          "note",
			})
	}
	return
}
