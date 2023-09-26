package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/generator"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type IPurchaseOrderItemRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (purchaseOrderItems []*model.PurchaseOrderItem, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (purchaseOrderItem *model.PurchaseOrderItem, err error)
	Create(ctx context.Context, req *model.PurchaseOrderItem) (purchaseOrderItem *model.PurchaseOrderItem, err error)
	GetByPurchaseOrderId(ctx context.Context, purchaseOrderId int64) (purchaseOrderItems []*model.PurchaseOrderItem, err error)
}

type PurchaseOrderItemRepository struct {
	opt opt.Options
}

func NewPurchaseOrderItemRepository() IPurchaseOrderItemRepository {
	return &PurchaseOrderItemRepository{
		opt: global.Setup.Common,
	}
}

func (r *PurchaseOrderItemRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (purchaseOrderItems []*model.PurchaseOrderItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderItemRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	count = 10
	for i := 1; i <= int(count); i++ {
		purchaseOrderItems = append(purchaseOrderItems, r.MockDatas(int64(i)))
	}
	return
}

func (r *PurchaseOrderItemRepository) GetDetail(ctx context.Context, id int64, code string) (purchaseOrderItem *model.PurchaseOrderItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderItemRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	purchaseOrderItem = r.MockDatas(id)
	return
}

func (r *PurchaseOrderItemRepository) Create(ctx context.Context, req *model.PurchaseOrderItem) (purchaseOrderItem *model.PurchaseOrderItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderItemRepository.Create")
	defer span.End()

	// RETURN DUMMIES
	purchaseOrderItem = r.MockDatas(1)
	return
}

func (r *PurchaseOrderItemRepository) GetByPurchaseOrderId(ctx context.Context, purchaseOrderId int64) (purchaseOrderItems []*model.PurchaseOrderItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderItemRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	for i := 1; i <= int(5); i++ {
		purchaseOrderItems = append(purchaseOrderItems, r.MockDatas(int64(i)))
	}
	return
}

func (r *PurchaseOrderItemRepository) MockDatas(id int64) (mockDatas *model.PurchaseOrderItem) {
	if id == 0 {
		id = 1
	}
	mockDatas = &model.PurchaseOrderItem{
		ID:                 id,
		PurchaseOrderID:    generator.DummyInt64(1, 10),
		PurchasePlanItemID: generator.DummyInt64(1, 10),
		ItemID:             generator.DummyInt64(1, 10),
		OrderQty:           0.25 * float64(generator.DummyInt(1, 50)),
		UnitPrice:          0.25 * float64(generator.DummyInt(1, 50)),
		TaxableItem:        1,
		IncludeTax:         1,
		TaxPercentage:      0.25 * float64(generator.DummyInt(1, 50)),
		TaxAmount:          0.25 * float64(generator.DummyInt(1, 50)),
		UnitPriceTax:       0.25 * float64(generator.DummyInt(1, 50)) * 100000,
		Subtotal:           0.25 * float64(generator.DummyInt(1, 50)) * 100000,
		Weight:             0.25 * float64(generator.DummyInt(1, 50)),
		Note:               "Dummy Note",
		PurchaseQty:        0.25 * float64(generator.DummyInt(1, 50)),
	}

	return mockDatas
}
