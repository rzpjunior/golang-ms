package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type ISalesInvoiceItemRepository interface {
	Get(ctx context.Context, salesInvoiceID int64) (salesInvoiceItems []*model.SalesInvoiceItem, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (salesInvoiceItem *model.SalesInvoiceItem, err error)
}

type SalesInvoiceItemRepository struct {
	opt opt.Options
}

func NewSalesInvoiceItemRepository() ISalesInvoiceItemRepository {
	return &SalesInvoiceItemRepository{
		opt: global.Setup.Common,
	}
}

func (r *SalesInvoiceItemRepository) Get(ctx context.Context, salesInvoiceID int64) (salesInvoiceItems []*model.SalesInvoiceItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesInvoiceItemRepository.Get")
	defer span.End()

	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil
}

func (r *SalesInvoiceItemRepository) GetDetail(ctx context.Context, id int64, code string) (salesInvoiceItem *model.SalesInvoiceItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesInvoiceItemRepository.GetDetail")
	defer span.End()

	return r.MockDatas(1)[0], nil
}

func (r *SalesInvoiceItemRepository) MockDatas(total int) (mockDatas []*model.SalesInvoiceItem) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.SalesInvoiceItem{
				ID:               int64(i),
				SalesInvoiceID:   1,
				SalesOrderItemID: 1,
				ItemID:           1,
				InvoiceQty:       3,
				UnitPrice:        10000,
				TaxableItem:      2,
				TaxPercentage:    0,
				Subtotal:         30000,
				Note:             "note",
			})
	}
	return
}
