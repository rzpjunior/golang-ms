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

type IPurchaseOrderRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (purchaseOrders []*model.PurchaseOrder, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (purchaseOrder *model.PurchaseOrder, err error)
	CreateWithItem(ctx context.Context, req *model.PurchaseOrder, items []*model.PurchaseOrderItem) (purchaseOrder *model.PurchaseOrder, err error)
	Update(ctx context.Context, po *model.PurchaseOrder, columns ...string) (err error)
	UpdateWithItem(ctx context.Context, req *model.PurchaseOrder, items []*model.PurchaseOrderItem) (purchaseOrder *model.PurchaseOrder, err error)
	CommitPurchaseOrder(ctx context.Context, po *model.PurchaseOrder) (err error)
	CancelPurchaseOrder(ctx context.Context, po *model.PurchaseOrder) (err error)
	UpdateProduct(ctx context.Context, req *model.PurchaseOrder, items []*model.PurchaseOrderItem) (purchaseOrder *model.PurchaseOrder, err error)
}

type PurchaseOrderRepository struct {
	opt opt.Options
}

func NewPurchaseOrderRepository() IPurchaseOrderRepository {
	return &PurchaseOrderRepository{
		opt: global.Setup.Common,
	}
}

func (r *PurchaseOrderRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (purchaseOrders []*model.PurchaseOrder, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	count = 10
	for i := 1; i <= int(count); i++ {
		purchaseOrders = append(purchaseOrders, r.MockDatas(int64(i)))
	}
	return
}

func (r *PurchaseOrderRepository) GetDetail(ctx context.Context, id int64, code string) (purchaseOrder *model.PurchaseOrder, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	purchaseOrder = r.MockDatas(id)
	return
}

func (r *PurchaseOrderRepository) CreateWithItem(ctx context.Context, req *model.PurchaseOrder, items []*model.PurchaseOrderItem) (purchaseOrder *model.PurchaseOrder, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderRepository.Create")
	defer span.End()

	// RETURN DUMMIES
	purchaseOrder = r.MockDatas(1)
	return
}

func (r *PurchaseOrderRepository) Update(ctx context.Context, po *model.PurchaseOrder, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderRepository.Update")
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

func (r *PurchaseOrderRepository) UpdateWithItem(ctx context.Context, req *model.PurchaseOrder, items []*model.PurchaseOrderItem) (purchaseOrder *model.PurchaseOrder, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderRepository.UpdateWithItem")
	defer span.End()

	// RETURN DUMMIES
	purchaseOrder = r.MockDatas(1)
	return
}

func (r *PurchaseOrderRepository) UpdateProduct(ctx context.Context, req *model.PurchaseOrder, items []*model.PurchaseOrderItem) (purchaseOrder *model.PurchaseOrder, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderRepository.UpdateProduct")
	defer span.End()

	// RETURN DUMMIES
	purchaseOrder = r.MockDatas(1)
	return
}

func (r *PurchaseOrderRepository) MockDatas(id int64) (mockDatas *model.PurchaseOrder) {
	if id == 0 {
		id = 1
	}
	mockDatas = &model.PurchaseOrder{
		ID:                     id,
		Code:                   fmt.Sprintf("PO%d", id),
		VendorID:               generator.DummyInt64(1, 10),
		SiteID:                 generator.DummyInt64(1, 10),
		TermPaymentPurID:       generator.DummyInt64(1, 10),
		VendorClassificationID: generator.DummyInt64(1, 10),
		PurchasePlanID:         generator.DummyInt64(1, 10),
		ConsolidatedShipmentID: generator.DummyInt64(1, 10),
		Status:                 int32(rand.Intn(2)*4 + 1),
		RecognitionDate:        generator.DummyTime(),
		EtaDate:                generator.DummyTime(),
		SiteAddress:            "Dummy SiteAddress",
		EtaTime:                "Dummy EtaTime",
		TaxPct:                 0,
		DeliveryFee:            0,
		TotalPrice:             0,
		TaxAmount:              0.25 * float64(generator.DummyInt(1, 50)) * 10000,
		TotalCharge:            0.25 * float64(generator.DummyInt(1, 50)),
		TotalInvoice:           0.25 * float64(generator.DummyInt(1, 50)),
		TotalWeight:            0.25 * float64(generator.DummyInt(1, 50)),
		Note:                   "Dummy Note",
		DeltaPrint:             0,
		Latitude:               1,
		Longitude:              1,
		UpdatedAt:              time.Now(),
		UpdatedBy:              1,
		CreatedAt:              time.Now(),
		CreatedBy:              1,
		CommittedAt:            time.Now(),
		CommittedBy:            generator.DummyInt64(1, 10),
		AssignedTo:             generator.DummyInt64(1, 10),
		AssignedBy:             generator.DummyInt64(1, 10),
		AssignedAt:             time.Now(),
		Locked:                 1,
		LockedBy:               0,
		HasFinishedGr:          int8(generator.DummyInt(0, 1)),
		PurchaseOrderItem:      r.MockDataItems(3),
	}

	return mockDatas
}

func (r *PurchaseOrderRepository) CommitPurchaseOrder(ctx context.Context, po *model.PurchaseOrder) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderRepository.Update")
	defer span.End()

	// TODO: integrate to GP
	return
}

func (r *PurchaseOrderRepository) CancelPurchaseOrder(ctx context.Context, po *model.PurchaseOrder) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PurchaseOrderRepository.Update")
	defer span.End()

	// TODO: integrate to GP
	return
}

func (r *PurchaseOrderRepository) MockDataItems(total int) (mockDatas []*model.PurchaseOrderItem) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas, &model.PurchaseOrderItem{
			ID:                 int64(i),
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
		})
	}

	return mockDatas
}
