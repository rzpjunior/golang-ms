package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/generator"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type IReceivingItemRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (result []*model.ReceivingItem, count int64, err error)
	GetDetail(ctx context.Context, id int64) (result *model.ReceivingItem, err error)
	Create(ctx context.Context, req *model.ReceivingItem) (result *model.ReceivingItem, err error)
	GetByItemTransferId(ctx context.Context, itemTransferId int64) (result []*model.ReceivingItem, err error)
	GetByReceivingId(ctx context.Context, id int64) (result []*model.ReceivingItem, err error)
}

type ReceivingItemRepository struct {
	opt opt.Options
}

func NewReceivingItemRepository() IReceivingItemRepository {
	return &ReceivingItemRepository{
		opt: global.Setup.Common,
	}
}

func (r *ReceivingItemRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (result []*model.ReceivingItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ReceivingItemRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	count = 10
	for i := 1; i <= int(count); i++ {
		result = append(result, r.MockDatas(int64(i)))
	}
	return
}

func (r *ReceivingItemRepository) GetDetail(ctx context.Context, id int64) (result *model.ReceivingItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ReceivingItemRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	result = r.MockDatas(id)
	return
}

func (r *ReceivingItemRepository) Create(ctx context.Context, req *model.ReceivingItem) (result *model.ReceivingItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ReceivingItemRepository.Create")
	defer span.End()

	// RETURN DUMMIES
	result = r.MockDatas(1)
	return
}

func (r *ReceivingItemRepository) GetByItemTransferId(ctx context.Context, itemTransferId int64) (result []*model.ReceivingItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ReceivingItemRepository.GetByItemTransferId")
	defer span.End()

	// RETURN DUMMIES
	for i := 1; i <= int(5); i++ {
		result = append(result, r.MockDatas(int64(i)))
	}
	return
}

func (r *ReceivingItemRepository) GetByReceivingId(ctx context.Context, id int64) (result []*model.ReceivingItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ReceivingRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	for i := 1; i <= int(5); i++ {
		result = append(result, r.MockDatas(int64(i)))
	}
	return
}

func (r *ReceivingItemRepository) MockDatas(id int64) (mockDatas *model.ReceivingItem) {
	if id == 0 {
		id = 1
	}
	mockDatas = &model.ReceivingItem{
		ID:                  id,
		PurchaseOrderItemID: generator.DummyInt64(1, 2),
		ItemTransferItemID:  generator.DummyInt64(1, 2),
		DeliverQty:          0.25 * float64(generator.DummyInt(1, 50)),
		ReceiveQty:          0.25 * float64(generator.DummyInt(1, 50)),
		RejectQty:           0.25 * float64(generator.DummyInt(1, 50)),
		RejectReason:        int8(generator.DummyInt(1, 2)),
		Weight:              0.25 * float64(generator.DummyInt(1, 50)),
		Note:                "Dummy Note",
	}

	return mockDatas
}
