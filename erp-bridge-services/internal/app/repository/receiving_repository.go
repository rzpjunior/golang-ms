package repository

import (
	"context"
	"fmt"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/generator"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type IReceivingRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (result []*model.Receiving, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (result *model.Receiving, err error)
	CreateWithItem(ctx context.Context, req *model.Receiving, items []*model.ReceivingItem) (result *model.Receiving, err error)
	Update(ctx context.Context, po *model.Receiving, columns ...string) (err error)
	UpdateWithItem(ctx context.Context, req *model.Receiving, items []*model.ReceivingItem) (result *model.Receiving, err error)
	Confirm(ctx context.Context, receiving *model.Receiving, receivingItem []*model.ReceivingItem) (result *model.Receiving, err error)
	GetByInbound(ctx context.Context, inboundType int32, id int64) (result []*model.Receiving, err error)
}

type ReceivingRepository struct {
	opt opt.Options
}

func NewReceivingRepository() IReceivingRepository {
	return &ReceivingRepository{
		opt: global.Setup.Common,
	}
}

func (r *ReceivingRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (result []*model.Receiving, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ReceivingRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	count = 10
	for i := 1; i <= int(count); i++ {
		result = append(result, r.MockDatas(int64(i)))
	}
	return
}

func (r *ReceivingRepository) GetDetail(ctx context.Context, id int64, code string) (result *model.Receiving, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemTransferRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	result = r.MockDatas(id)
	return
}

func (r *ReceivingRepository) CreateWithItem(ctx context.Context, req *model.Receiving, items []*model.ReceivingItem) (result *model.Receiving, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemTransferRepository.Create")
	defer span.End()

	// RETURN DUMMIES
	result = r.MockDatas(1)
	return
}

func (r *ReceivingRepository) Update(ctx context.Context, po *model.Receiving, columns ...string) (err error) {
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

func (r *ReceivingRepository) UpdateWithItem(ctx context.Context, req *model.Receiving, items []*model.ReceivingItem) (result *model.Receiving, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemTransferRepository.UpdateWithItem")
	defer span.End()

	// RETURN DUMMIES
	result = r.MockDatas(1)
	return
}

func (r *ReceivingRepository) Confirm(ctx context.Context, receiving *model.Receiving, receivingItem []*model.ReceivingItem) (result *model.Receiving, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemTransferRepository.Update")
	defer span.End()

	// RETURN DUMMIES
	result = r.MockDatas(1)
	return
}

func (r *ReceivingRepository) MockDatas(id int64) (mockDatas *model.Receiving) {
	if id == 0 {
		id = 1
	}
	mockDatas = &model.Receiving{
		ID:                  id,
		Code:                fmt.Sprintf("RCV%d", id),
		SiteId:              1,
		PurchaseOrderId:     1,
		ItemTransferId:      1,
		InboundType:         int8(generator.DummyInt(1, 2)),
		ValidSupplierReturn: 1,
		AtaDate:             generator.DummyTime(),
		AtaTime:             generator.DummyTime().Format("15:04"),
		StockType:           int8(generator.DummyInt(1, 2)),
		TotalWeight:         0.25 * float64(generator.DummyInt(1, 50)),
		Note:                "Dummy Note",
		Status:              5,
		Locked:              1,
		LockedBy:            0,
		CreatedAt:           generator.DummyTime(),
		CreatedBy:           1,
		ConfirmedAt:         generator.DummyTime(),
		ConfirmedBy:         1,
		UpdatedAt:           time.Now(),
		UpdatedBy:           1,
	}

	return mockDatas
}

func (r *ReceivingRepository) MockDataItems(total int) (mockDatas []*model.ReceivingItem) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas, &model.ReceivingItem{
			ID:           int64(i),
			DeliverQty:   generator.DummyFloat64(1, 50),
			ReceiveQty:   generator.DummyFloat64(1, 50),
			RejectQty:    generator.DummyFloat64(1, 50),
			RejectReason: int8(generator.DummyInt(1, 8)),
			Weight:       0.25 * float64(generator.DummyInt(1, 50)),
			Note:         "Dummy Note",
		})
	}

	return mockDatas
}

func (r *ReceivingRepository) GetByInbound(ctx context.Context, inboundType int32, id int64) (result []*model.Receiving, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ReceivingRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	for i := 1; i <= int(5); i++ {
		result = append(result, r.MockDatas(int64(i)))
	}
	return
}
