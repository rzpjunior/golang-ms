package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/generator"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type IItemTransferItemRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (itemTransferItems []*model.ItemTransferItem, count int64, err error)
	GetDetail(ctx context.Context, id int64) (itemTransferItem *model.ItemTransferItem, err error)
	Create(ctx context.Context, req *model.ItemTransferItem) (itemTransferItem *model.ItemTransferItem, err error)
	GetByItemTransferId(ctx context.Context, ItemTransferId int64) (itemTransferItems []*model.ItemTransferItem, err error)
}

type ItemTransferItemRepository struct {
	opt opt.Options
}

func NewItemTransferItemRepository() IItemTransferItemRepository {
	return &ItemTransferItemRepository{
		opt: global.Setup.Common,
	}
}

func (r *ItemTransferItemRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (itemTransferItems []*model.ItemTransferItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemTransferItemRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	count = 10
	for i := 1; i <= int(count); i++ {
		itemTransferItems = append(itemTransferItems, r.MockDatas(int64(i)))
	}
	return
}

func (r *ItemTransferItemRepository) GetDetail(ctx context.Context, id int64) (itemTransferItem *model.ItemTransferItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemTransferItemRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	itemTransferItem = r.MockDatas(id)
	return
}

func (r *ItemTransferItemRepository) Create(ctx context.Context, req *model.ItemTransferItem) (itemTransferItem *model.ItemTransferItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemTransferItemRepository.Create")
	defer span.End()

	// RETURN DUMMIES
	itemTransferItem = r.MockDatas(1)
	return
}

func (r *ItemTransferItemRepository) GetByItemTransferId(ctx context.Context, itemTransferId int64) (itemTransferItems []*model.ItemTransferItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemTransferItemRepository.GetByItemTransferId")
	defer span.End()

	// RETURN DUMMIES
	for i := 1; i <= int(5); i++ {
		itemTransferItems = append(itemTransferItems, r.MockDatas(int64(i)))
	}
	return
}

func (r *ItemTransferItemRepository) MockDatas(id int64) (mockDatas *model.ItemTransferItem) {
	if id == 0 {
		id = 1
	}
	mockDatas = &model.ItemTransferItem{
		ID:             id,
		ItemTransferID: generator.DummyInt64(1, 50),
		DeliverQty:     0.25 * float64(generator.DummyInt(1, 50)),
		ReceiveQty:     0.25 * float64(generator.DummyInt(1, 50)),
		RequestQty:     0.25 * float64(generator.DummyInt(1, 50)),
		ReceiveNote:    "Dummy Receive Note",
		UnitCost:       0.25 * float64(generator.DummyInt(1, 50)),
		Subtotal:       0.25 * float64(generator.DummyInt(1, 50)) * 100000,
		Weight:         0.25 * float64(generator.DummyInt(1, 50)),
		Note:           "Dummy Note",
	}

	return mockDatas
}
