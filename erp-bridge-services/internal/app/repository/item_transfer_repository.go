package repository

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/generator"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type IItemTransferRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (itemTransfers []*model.ItemTransfer, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (ItemTransfer *model.ItemTransfer, err error)
	CreateWithItem(ctx context.Context, req *model.ItemTransfer, items []*model.ItemTransferItem) (itemTransfer *model.ItemTransfer, err error)
	Update(ctx context.Context, po *model.ItemTransfer, columns ...string) (err error)
	UpdateWithItem(ctx context.Context, req *model.ItemTransfer, items []*model.ItemTransferItem) (itemTransfer *model.ItemTransfer, err error)
	CommitItemTransfer(ctx context.Context, itemTransfer *model.ItemTransfer, itemTransferItem []*model.ItemTransferItem) (result *model.ItemTransfer, err error)
}

type ItemTransferRepository struct {
	opt opt.Options
}

func NewItemTransferRepository() IItemTransferRepository {
	return &ItemTransferRepository{
		opt: global.Setup.Common,
	}
}

func (r *ItemTransferRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (itemTransfers []*model.ItemTransfer, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemTransferRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	count = 10
	for i := 1; i <= int(count); i++ {
		itemTransfers = append(itemTransfers, r.MockDatas(int64(i)))
	}
	return
}

func (r *ItemTransferRepository) GetDetail(ctx context.Context, id int64, code string) (itemTransfer *model.ItemTransfer, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemTransferRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	itemTransfer = r.MockDatas(id)
	return
}

func (r *ItemTransferRepository) CreateWithItem(ctx context.Context, req *model.ItemTransfer, items []*model.ItemTransferItem) (itemTransfer *model.ItemTransfer, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemTransferRepository.Create")
	defer span.End()

	// RETURN DUMMIES
	itemTransfer = r.MockDatas(1)
	return
}

func (r *ItemTransferRepository) Update(ctx context.Context, po *model.ItemTransfer, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemTransferRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, po, columns...)

	if err == nil {
		err = tx.Commit()
		if err != nil {
			span.RecordError(err)
			return
		}
	} else {
		span.RecordError(err)
		err = tx.Rollback()
		if err != nil {
			span.RecordError(err)
			return
		}
	}
	return
}

func (r *ItemTransferRepository) UpdateWithItem(ctx context.Context, req *model.ItemTransfer, items []*model.ItemTransferItem) (itemTransfer *model.ItemTransfer, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemTransferRepository.UpdateWithItem")
	defer span.End()

	// RETURN DUMMIES
	itemTransfer = r.MockDatas(1)
	return
}

func (r *ItemTransferRepository) CommitItemTransfer(ctx context.Context, itemTransfer *model.ItemTransfer, itemTransferItem []*model.ItemTransferItem) (result *model.ItemTransfer, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemTransferRepository.Update")
	defer span.End()

	// RETURN DUMMIES
	result = r.MockDatas(1)
	return
}

func (r *ItemTransferRepository) MockDatas(id int64) (mockDatas *model.ItemTransfer) {
	if id == 0 {
		id = 1
	}
	mockDatas = &model.ItemTransfer{
		ID:                 id,
		Code:               fmt.Sprintf("IT%d", id),
		RequestDate:        generator.DummyTime(),
		RecognitionDate:    generator.DummyTime(),
		EtaDate:            generator.DummyTime(),
		EtaTime:            generator.DummyTime().Format("15:04"),
		AtaDate:            generator.DummyTime(),
		AtaTime:            generator.DummyTime().Format("15:04"),
		AdditionalCost:     float64(generator.DummyInt(1, 50)) * 10000,
		AdditionalCostNote: "Dummy Additional Cost Note",
		StockType:          int8(generator.DummyInt(1, 2)),
		TotalCost:          float64(generator.DummyInt(1, 50)) * 10000,
		TotalCharge:        float64(generator.DummyInt(1, 50)) * 10000,
		TotalWeight:        0.25 * float64(generator.DummyInt(1, 50)),
		Note:               "Dummy Note",
		Status:             int8(rand.Intn(2)*4 + 1),
		Locked:             1,
		LockedBy:           0,
		TotalSku:           generator.DummyInt64(1, 50),
		UpdatedAt:          time.Now(),
		UpdatedBy:          1,
		SiteOriginID:       1,
		SiteDestinationID:  2,
		HasFinishedGr:      int8(generator.DummyInt(0, 1)),
		ItemTransferItem:   r.MockDataItems(3),
	}

	return mockDatas
}

func (r *ItemTransferRepository) MockDataItems(total int) (mockDatas []*model.ItemTransferItem) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas, &model.ItemTransferItem{
			ID:          int64(i),
			DeliverQty:  generator.DummyFloat64(1, 50),
			ReceiveQty:  generator.DummyFloat64(1, 50),
			RequestQty:  generator.DummyFloat64(1, 50),
			ReceiveNote: "Dummy Receive Note",
			UnitCost:    generator.DummyFloat64(1, 50),
			Subtotal:    0.25 * float64(generator.DummyInt(1, 50)) * 100000,
			Weight:      0.25 * float64(generator.DummyInt(1, 50)),
			Note:        "Dummy Note",
		})
	}

	return mockDatas
}
