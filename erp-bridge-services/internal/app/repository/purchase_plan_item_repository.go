package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/generator"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type IPurchasePlanItemRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (purchasePlanItems []*model.PurchasePlanItem, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (purchasePlanItem *model.PurchasePlanItem, err error)
}

type PurchasePlanItemRepository struct {
	opt opt.Options
}

func NewPurchasePlanItemRepository() IPurchasePlanItemRepository {
	return &PurchasePlanItemRepository{
		opt: global.Setup.Common,
	}
}

func (r *PurchasePlanItemRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (purchasePlanItems []*model.PurchasePlanItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchasePlanItemRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	count = 10
	for i := 1; i <= int(count); i++ {
		purchasePlanItems = append(purchasePlanItems, r.MockDatas(int64(i)))
	}
	return
}

func (r *PurchasePlanItemRepository) GetDetail(ctx context.Context, id int64, code string) (purchasePlanItem *model.PurchasePlanItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchasePlanItemRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	purchasePlanItem = r.MockDatas(id)
	return
}

func (r *PurchasePlanItemRepository) MockDatas(id int64) (mockDatas *model.PurchasePlanItem) {
	if id == 0 {
		id = 1
	}
	mockDatas = &model.PurchasePlanItem{
		ID:              id,
		PurchasePlanID:  generator.DummyInt64(1, 10),
		ItemID:          generator.DummyInt64(1, 10),
		PurchasePlanQty: 0.25 * float64(generator.DummyInt(1, 50)),
		PurchaseQty:     0.25 * float64(generator.DummyInt(1, 50)),
		UnitPrice:       generator.DummyFloat64(10000, 100000),
		Subtotal:        generator.DummyFloat64(10000, 100000),
		Weight:          0.25 * float64(generator.DummyInt(1, 50)),
	}

	return mockDatas
}
