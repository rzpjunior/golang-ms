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

type ISalesOrderItemService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, salesOrderID int64, itemID int64) (res []dto.SalesOrderItemResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.SalesOrderItemResponse, err error)
}

type SalesOrderItemService struct {
	opt                      opt.Options
	RepositorySalesOrderItem repository.ISalesOrderItemRepository
}

func NewSalesOrderItemService() ISalesOrderItemService {
	return &SalesOrderItemService{
		opt:                      global.Setup.Common,
		RepositorySalesOrderItem: repository.NewSalesOrderItemRepository(),
	}
}

func (s *SalesOrderItemService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, salesOrderID int64, itemID int64) (res []dto.SalesOrderItemResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderItemService.Get")
	defer span.End()

	var salesOrders []*model.SalesOrderItem
	salesOrders, total, err = s.RepositorySalesOrderItem.Get(ctx, offset, limit, status, search, orderBy, salesOrderID, itemID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, salesOrder := range salesOrders {
		res = append(res, dto.SalesOrderItemResponse{
			ID:            salesOrder.ID,
			SalesOrderId:  salesOrder.SalesOrderId,
			ItemId:        salesOrder.ItemId,
			OrderQty:      salesOrder.OrderQty,
			DefaultPrice:  salesOrder.DefaultPrice,
			UnitPrice:     salesOrder.UnitPrice,
			TaxableItem:   salesOrder.TaxableItem,
			TaxPercentage: salesOrder.TaxPercentage,
			ShadowPrice:   salesOrder.ShadowPrice,
			Subtotal:      salesOrder.Subtotal,
			Weight:        salesOrder.Weight,
			Note:          salesOrder.Note,
		})
	}

	return
}

func (s *SalesOrderItemService) GetDetail(ctx context.Context, id int64, code string) (res dto.SalesOrderItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderItemService.GetDetail")
	defer span.End()

	var salesOrder *model.SalesOrderItem
	salesOrder, err = s.RepositorySalesOrderItem.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.SalesOrderItemResponse{
		ID:            salesOrder.ID,
		SalesOrderId:  salesOrder.SalesOrderId,
		ItemId:        salesOrder.ItemId,
		OrderQty:      salesOrder.OrderQty,
		DefaultPrice:  salesOrder.DefaultPrice,
		UnitPrice:     salesOrder.UnitPrice,
		TaxableItem:   salesOrder.TaxableItem,
		TaxPercentage: salesOrder.TaxPercentage,
		ShadowPrice:   salesOrder.ShadowPrice,
		Subtotal:      salesOrder.Subtotal,
		Weight:        salesOrder.Weight,
		Note:          salesOrder.Note,
	}

	return
}
