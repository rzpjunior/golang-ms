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

type IPurchaseOrderItemService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.PurchaseOrderItemResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.PurchaseOrderItemResponse, err error)
}

type PurchaseOrderItemService struct {
	opt                         opt.Options
	RepositoryPurchaseOrderItem repository.IPurchaseOrderItemRepository
}

func NewPurchaseOrderItemService() IPurchaseOrderItemService {
	return &PurchaseOrderItemService{
		opt:                         global.Setup.Common,
		RepositoryPurchaseOrderItem: repository.NewPurchaseOrderItemRepository(),
	}
}

func (s *PurchaseOrderItemService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.PurchaseOrderItemResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderItemService.Get")
	defer span.End()

	var purchaseOrderItems []*model.PurchaseOrderItem
	purchaseOrderItems, total, err = s.RepositoryPurchaseOrderItem.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, purchaseOrderItem := range purchaseOrderItems {
		res = append(res, dto.PurchaseOrderItemResponse{
			ID:                 purchaseOrderItem.ID,
			PurchaseOrderID:    purchaseOrderItem.PurchaseOrderID,
			PurchasePlanItemID: purchaseOrderItem.PurchasePlanItemID,
			ItemID:             purchaseOrderItem.ItemID,
			OrderQty:           purchaseOrderItem.OrderQty,
			UnitPrice:          purchaseOrderItem.UnitPrice,
			TaxableItem:        purchaseOrderItem.TaxableItem,
			IncludeTax:         purchaseOrderItem.IncludeTax,
			TaxPercentage:      purchaseOrderItem.TaxPercentage,
			TaxAmount:          purchaseOrderItem.TaxAmount,
			UnitPriceTax:       purchaseOrderItem.UnitPriceTax,
			Subtotal:           purchaseOrderItem.Subtotal,
			Weight:             purchaseOrderItem.Weight,
			Note:               purchaseOrderItem.Note,
			PurchaseQty:        purchaseOrderItem.PurchaseQty,
		})
	}

	return
}

func (s *PurchaseOrderItemService) GetDetail(ctx context.Context, id int64, code string) (res dto.PurchaseOrderItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderItemService.GetDetail")
	defer span.End()

	var purchaseOrderItem *model.PurchaseOrderItem
	purchaseOrderItem, err = s.RepositoryPurchaseOrderItem.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.PurchaseOrderItemResponse{
		ID:                 purchaseOrderItem.ID,
		PurchaseOrderID:    purchaseOrderItem.PurchaseOrderID,
		PurchasePlanItemID: purchaseOrderItem.PurchasePlanItemID,
		ItemID:             purchaseOrderItem.ItemID,
		OrderQty:           purchaseOrderItem.OrderQty,
		UnitPrice:          purchaseOrderItem.UnitPrice,
		TaxableItem:        purchaseOrderItem.TaxableItem,
		IncludeTax:         purchaseOrderItem.IncludeTax,
		TaxPercentage:      purchaseOrderItem.TaxPercentage,
		TaxAmount:          purchaseOrderItem.TaxAmount,
		UnitPriceTax:       purchaseOrderItem.UnitPriceTax,
		Subtotal:           purchaseOrderItem.Subtotal,
		Weight:             purchaseOrderItem.Weight,
		Note:               purchaseOrderItem.Note,
		PurchaseQty:        purchaseOrderItem.PurchaseQty,
	}

	return
}
