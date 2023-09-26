package service

import (
	"context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
)

type IPurchasePlanItemService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.PurchasePlanItemResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.PurchasePlanItemResponse, err error)
}

type PurchasePlanItemService struct {
	opt                        opt.Options
	RepositoryPurchasePlanItem repository.IPurchasePlanItemRepository
}

func NewPurchasePlanItemService() IPurchasePlanItemService {
	return &PurchasePlanItemService{
		opt:                        global.Setup.Common,
		RepositoryPurchasePlanItem: repository.NewPurchasePlanItemRepository(),
	}
}

func (s *PurchasePlanItemService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.PurchasePlanItemResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchasePlanItemService.Get")
	defer span.End()
	fmt.Print("----")

	var purchasePlanItems []*model.PurchasePlanItem
	purchasePlanItems, total, err = s.RepositoryPurchasePlanItem.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, purchasePlanItem := range purchasePlanItems {
		res = append(res, dto.PurchasePlanItemResponse{
			ID:              purchasePlanItem.ID,
			PurchasePlanID:  purchasePlanItem.PurchasePlanID,
			ItemID:          purchasePlanItem.ItemID,
			PurchasePlanQty: purchasePlanItem.PurchasePlanQty,
			PurchaseQty:     purchasePlanItem.PurchaseQty,
			UnitPrice:       purchasePlanItem.UnitPrice,
			Subtotal:        purchasePlanItem.Subtotal,
			Weight:          purchasePlanItem.Weight,
		})
	}

	return
}

func (s *PurchasePlanItemService) GetDetail(ctx context.Context, id int64, code string) (res dto.PurchasePlanItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchasePlanItemService.GetDetail")
	defer span.End()

	var purchasePlanItem *model.PurchasePlanItem
	purchasePlanItem, err = s.RepositoryPurchasePlanItem.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.PurchasePlanItemResponse{
		ID:              purchasePlanItem.ID,
		PurchasePlanID:  purchasePlanItem.PurchasePlanID,
		ItemID:          purchasePlanItem.ItemID,
		PurchasePlanQty: purchasePlanItem.PurchasePlanQty,
		PurchaseQty:     purchasePlanItem.PurchaseQty,
		UnitPrice:       purchasePlanItem.UnitPrice,
		Subtotal:        purchasePlanItem.Subtotal,
		Weight:          purchasePlanItem.Weight,
	}

	return
}
