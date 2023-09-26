package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
)

type ISalesInvoiceItemService interface {
	Get(ctx context.Context, salesInvoiceItemID int64) (res []dto.SalesInvoiceItemResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.SalesInvoiceItemResponse, err error)
}

type SalesInvoiceItemService struct {
	opt                        opt.Options
	RepositorySalesInvoiceItem repository.ISalesInvoiceItemRepository
}

func NewSalesInvoiceItemService() ISalesInvoiceItemService {
	return &SalesInvoiceItemService{
		opt:                        global.Setup.Common,
		RepositorySalesInvoiceItem: repository.NewSalesInvoiceItemRepository(),
	}
}

func (s *SalesInvoiceItemService) Get(ctx context.Context, salesInvoiceItemID int64) (res []dto.SalesInvoiceItemResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesInvoiceItemService.Get")
	defer span.End()

	var salesInvoiceItems []*model.SalesInvoiceItem
	salesInvoiceItems, total, err = s.RepositorySalesInvoiceItem.Get(ctx, salesInvoiceItemID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, salesInvoiceItem := range salesInvoiceItems {
		res = append(res, dto.SalesInvoiceItemResponse{
			ID:               salesInvoiceItem.ID,
			SalesOrderItemID: salesInvoiceItem.SalesOrderItemID,
			SalesInvoiceID:   salesInvoiceItem.SalesInvoiceID,
			ItemID:           salesInvoiceItem.ItemID,
			Note:             salesInvoiceItem.Note,
			InvoiceQty:       salesInvoiceItem.InvoiceQty,
			UnitPrice:        salesInvoiceItem.UnitPrice,
			Subtotal:         salesInvoiceItem.Subtotal,
			TaxableItem:      salesInvoiceItem.TaxableItem,
			TaxPercentage:    salesInvoiceItem.TaxPercentage,
			SkuDiscAmount:    salesInvoiceItem.SkuDiscAmount,
		})
	}

	return
}

func (s *SalesInvoiceItemService) GetDetail(ctx context.Context, id int64, code string) (res dto.SalesInvoiceItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesInvoiceItemService.GetDetail")
	defer span.End()

	var salesInvoiceItem *model.SalesInvoiceItem
	salesInvoiceItem, err = s.RepositorySalesInvoiceItem.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.SalesInvoiceItemResponse{
		ID:               salesInvoiceItem.ID,
		SalesOrderItemID: salesInvoiceItem.SalesOrderItemID,
		SalesInvoiceID:   salesInvoiceItem.SalesInvoiceID,
		ItemID:           salesInvoiceItem.ItemID,
		Note:             salesInvoiceItem.Note,
		InvoiceQty:       salesInvoiceItem.InvoiceQty,
		UnitPrice:        salesInvoiceItem.UnitPrice,
		Subtotal:         salesInvoiceItem.Subtotal,
		TaxableItem:      salesInvoiceItem.TaxableItem,
		TaxPercentage:    salesInvoiceItem.TaxPercentage,
		SkuDiscAmount:    salesInvoiceItem.SkuDiscAmount,
	}

	return
}
