package service

import (
	"context"
	"encoding/json"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
)

type IItemTransferItemService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.ItemTransferItemResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64) (res dto.ItemTransferItemResponse, err error)
}

type ItemTransferItemService struct {
	opt                        opt.Options
	RepositoryItemTransfer     repository.IItemTransferRepository
	RepositoryItemTransferItem repository.IItemTransferItemRepository
}

func NewItemTransferItemService() IItemTransferItemService {
	return &ItemTransferItemService{
		opt:                        global.Setup.Common,
		RepositoryItemTransfer:     repository.NewItemTransferRepository(),
		RepositoryItemTransferItem: repository.NewItemTransferItemRepository(),
	}
}

func (s *ItemTransferItemService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.ItemTransferItemResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemTransferItemService.Get")
	defer span.End()

	var items []*model.ItemTransferItem
	items, total, err = s.RepositoryItemTransferItem.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, item := range items {
		res = append(res, dto.ItemTransferItemResponse{
			ID:          item.ID,
			DeliverQty:  item.DeliverQty,
			RequestQty:  item.RequestQty,
			ReceiveQty:  item.ReceiveQty,
			ReceiveNote: item.ReceiveNote,
			UnitCost:    item.UnitCost,
			Subtotal:    item.Subtotal,
			Weight:      item.Weight,
			Note:        item.Note,
		})
	}

	jsonRes, _ := json.Marshal(res)
	fmt.Println(jsonRes)

	return
}

func (s *ItemTransferItemService) GetDetail(ctx context.Context, id int64) (res dto.ItemTransferItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemTransferItemService.GetDetail")
	defer span.End()

	var (
		itemTransfer *model.ItemTransfer
		item         *model.ItemTransferItem
	)
	item, err = s.RepositoryItemTransferItem.GetDetail(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	itemTransfer, err = s.RepositoryItemTransfer.GetDetail(ctx, item.ItemTransferID, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ItemTransferItemResponse{
		ID:          item.ID,
		DeliverQty:  item.DeliverQty,
		RequestQty:  item.RequestQty,
		ReceiveQty:  item.ReceiveQty,
		ReceiveNote: item.ReceiveNote,
		UnitCost:    item.UnitCost,
		Subtotal:    item.Subtotal,
		Weight:      item.Weight,
		Note:        item.Note,
		ItemTransfer: &dto.ItemTransferResponse{
			ID:                 itemTransfer.ID,
			Code:               itemTransfer.Code,
			RequestDate:        itemTransfer.RequestDate,
			RecognitionDate:    itemTransfer.RecognitionDate,
			EtaDate:            itemTransfer.EtaDate,
			EtaTime:            itemTransfer.EtaTime,
			AtaDate:            itemTransfer.AtaDate,
			AtaTime:            itemTransfer.AtaTime,
			AdditionalCost:     itemTransfer.AdditionalCost,
			AdditionalCostNote: itemTransfer.AdditionalCostNote,
			StockType:          itemTransfer.StockType,
			TotalCost:          itemTransfer.TotalCost,
			TotalCharge:        itemTransfer.TotalCharge,
			TotalWeight:        itemTransfer.TotalWeight,
			Note:               itemTransfer.Note,
			Status:             itemTransfer.Status,
			Locked:             itemTransfer.Locked,
			LockedBy:           itemTransfer.LockedBy,
			TotalSku:           itemTransfer.TotalSku,
			UpdatedAt:          itemTransfer.UpdatedAt,
			UpdatedBy:          itemTransfer.UpdatedBy,
		},
	}

	jsonRes, _ := json.Marshal(res)
	fmt.Println(jsonRes)

	return
}
