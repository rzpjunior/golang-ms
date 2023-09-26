package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *BridgeGrpcHandler) GetItemTransferList(ctx context.Context, req *bridgeService.GetItemTransferListRequest) (res *bridgeService.GetItemTransferListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemTransferList")
	defer span.End()

	var itemTransfers []dto.ItemTransferResponse
	itemTransfers, _, err = h.ServiceItemTransfer.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.ItemTransfer
	for _, itemTransfer := range itemTransfers {
		var grs []*bridgeService.ReceivingListinDetailResponse
		for _, gr := range itemTransfer.Receiving {
			grs = append(grs, &bridgeService.ReceivingListinDetailResponse{
				Id:     gr.ID,
				Code:   gr.Code,
				Status: int32(gr.Status),
			})
		}
		data = append(data, &bridgeService.ItemTransfer{
			Id:                 itemTransfer.ID,
			Code:               itemTransfer.Code,
			RequestDate:        timestamppb.New(itemTransfer.RequestDate),
			RecognitionDate:    timestamppb.New(itemTransfer.RecognitionDate),
			EtaDate:            timestamppb.New(itemTransfer.EtaDate),
			EtaTime:            itemTransfer.EtaTime,
			AtaDate:            timestamppb.New(itemTransfer.AtaDate),
			AtaTime:            itemTransfer.AtaTime,
			AdditionalCost:     itemTransfer.AdditionalCost,
			AdditionalCostNote: itemTransfer.AdditionalCostNote,
			StockType:          int32(itemTransfer.StockType),
			TotalCost:          itemTransfer.TotalCost,
			TotalSku:           itemTransfer.TotalSku,
			TotalCharge:        itemTransfer.TotalCharge,
			TotalWeight:        itemTransfer.TotalWeight,
			Note:               itemTransfer.Note,
			Status:             int32(itemTransfer.Status),
			UpdatedAt:          timestamppb.New(itemTransfer.UpdatedAt),
			UpdatedBy:          itemTransfer.UpdatedBy,
			Locked:             int32(itemTransfer.Locked),
			LockedBy:           itemTransfer.LockedBy,
			Receiving:          grs,
			SiteOrigin: &bridgeService.Site{
				Id:          itemTransfer.SiteOrigin.ID,
				Code:        itemTransfer.SiteOrigin.Code,
				Description: itemTransfer.SiteOrigin.Description,
				Status:      int32(itemTransfer.SiteOrigin.Status),
				CreatedAt:   timestamppb.New(itemTransfer.SiteOrigin.CreatedAt),
				UpdatedAt:   timestamppb.New(itemTransfer.SiteOrigin.UpdatedAt),
			},
			SiteDestination: &bridgeService.Site{
				Id:          itemTransfer.SiteDestination.ID,
				Code:        itemTransfer.SiteDestination.Code,
				Description: itemTransfer.SiteDestination.Description,
				Status:      int32(itemTransfer.SiteDestination.Status),
				CreatedAt:   timestamppb.New(itemTransfer.SiteDestination.CreatedAt),
				UpdatedAt:   timestamppb.New(itemTransfer.SiteDestination.UpdatedAt),
			},
		})
	}

	res = &bridgeService.GetItemTransferListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetItemTransferDetail(ctx context.Context, req *bridgeService.GetItemTransferDetailRequest) (res *bridgeService.GetItemTransferDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemTransferDetail")
	defer span.End()

	var (
		itemTransfer dto.ItemTransferResponse
		items        []*bridgeService.ItemTransferItem
		grs          []*bridgeService.ReceivingListinDetailResponse
	)
	itemTransfer, err = h.ServiceItemTransfer.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	for _, item := range itemTransfer.ItemTransferItem {
		items = append(items, &bridgeService.ItemTransferItem{
			Id:          item.ID,
			DeliverQty:  item.DeliverQty,
			ReceiveQty:  item.ReceiveQty,
			RequestQty:  item.RequestQty,
			ReceiveNote: item.ReceiveNote,
			UnitCost:    item.UnitCost,
			Subtotal:    item.Subtotal,
			Weight:      item.Weight,
			Note:        item.Note,
			Item: &bridgeService.Item{
				Id:                   item.Item.ID,
				Code:                 item.Item.Code,
				UomId:                item.Item.UomID,
				ClassId:              item.Item.ClassID,
				ItemCategoryId:       item.Item.ItemCategoryID,
				Description:          item.Item.Description,
				UnitWeightConversion: item.Item.UnitWeightConversion,
				OrderMinQty:          item.Item.OrderMinQty,
				OrderMaxQty:          item.Item.OrderMaxQty,
				ItemType:             item.Item.ItemType,
				Uom: &bridgeService.Uom{
					Id:             item.Item.Uom.ID,
					Code:           item.Item.Uom.Code,
					Description:    item.Item.Uom.Description,
					Status:         int32(item.Item.Uom.Status),
					CreatedAt:      timestamppb.New(item.Item.Uom.CreatedAt),
					UpdatedAt:      timestamppb.New(item.Item.Uom.UpdatedAt),
					DecimalEnabled: int32(item.Item.Uom.DecimalEnabled),
				},
			},
		})
	}

	for _, gr := range itemTransfer.Receiving {
		grs = append(grs, &bridgeService.ReceivingListinDetailResponse{
			Id:     gr.ID,
			Code:   gr.Code,
			Status: int32(gr.Status),
		})
	}

	res = &bridgeService.GetItemTransferDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.ItemTransfer{
			Id:                 itemTransfer.ID,
			Code:               itemTransfer.Code,
			RequestDate:        timestamppb.New(itemTransfer.RequestDate),
			RecognitionDate:    timestamppb.New(itemTransfer.RecognitionDate),
			EtaDate:            timestamppb.New(itemTransfer.EtaDate),
			EtaTime:            itemTransfer.EtaTime,
			AtaDate:            timestamppb.New(itemTransfer.AtaDate),
			AtaTime:            itemTransfer.AtaTime,
			AdditionalCost:     itemTransfer.AdditionalCost,
			AdditionalCostNote: itemTransfer.AdditionalCostNote,
			StockType:          int32(itemTransfer.StockType),
			TotalCost:          itemTransfer.TotalCost,
			TotalSku:           itemTransfer.TotalSku,
			TotalCharge:        itemTransfer.TotalCharge,
			TotalWeight:        itemTransfer.TotalWeight,
			Note:               itemTransfer.Note,
			Status:             int32(itemTransfer.Status),
			UpdatedAt:          timestamppb.New(itemTransfer.UpdatedAt),
			UpdatedBy:          itemTransfer.UpdatedBy,
			Locked:             int32(itemTransfer.Locked),
			LockedBy:           itemTransfer.LockedBy,
			SiteOriginId:       itemTransfer.SiteOriginID,
			SiteDestinationId:  itemTransfer.SiteDestinationID,
			ItemTransferItem:   items,
			SiteOrigin: &bridgeService.Site{
				Id:          itemTransfer.SiteOrigin.ID,
				Code:        itemTransfer.SiteOrigin.Code,
				Description: itemTransfer.SiteOrigin.Description,
				Status:      int32(itemTransfer.SiteOrigin.Status),
				CreatedAt:   timestamppb.New(itemTransfer.SiteOrigin.CreatedAt),
				UpdatedAt:   timestamppb.New(itemTransfer.SiteOrigin.UpdatedAt),
				Region: &bridgeService.Region{
					Id:          itemTransfer.SiteOrigin.Region.ID,
					Code:        itemTransfer.SiteOrigin.Region.Code,
					Description: itemTransfer.SiteOrigin.Region.Description,
					Status:      int32(itemTransfer.SiteOrigin.Region.Status),
					CreatedAt:   timestamppb.New(itemTransfer.SiteOrigin.Region.CreatedAt),
					UpdatedAt:   timestamppb.New(itemTransfer.SiteOrigin.Region.UpdatedAt),
				},
			},
			SiteDestination: &bridgeService.Site{
				Id:          itemTransfer.SiteDestination.ID,
				Code:        itemTransfer.SiteDestination.Code,
				Description: itemTransfer.SiteDestination.Description,
				Status:      int32(itemTransfer.SiteDestination.Status),
				CreatedAt:   timestamppb.New(itemTransfer.SiteDestination.CreatedAt),
				UpdatedAt:   timestamppb.New(itemTransfer.SiteDestination.UpdatedAt),
				Region: &bridgeService.Region{
					Id:          itemTransfer.SiteDestination.Region.ID,
					Code:        itemTransfer.SiteDestination.Region.Code,
					Description: itemTransfer.SiteDestination.Region.Description,
					Status:      int32(itemTransfer.SiteDestination.Region.Status),
					CreatedAt:   timestamppb.New(itemTransfer.SiteDestination.Region.CreatedAt),
					UpdatedAt:   timestamppb.New(itemTransfer.SiteDestination.Region.UpdatedAt),
				},
			},
			Receiving: grs,
		},
	}
	return
}

func (h *BridgeGrpcHandler) CreateItemTransfer(ctx context.Context, req *bridgeService.CreateItemTransferRequest) (res *bridgeService.GetItemTransferDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateItemTransfer")
	defer span.End()

	var (
		itemTransfer dto.ItemTransferResponse
		items        []*dto.CreateItemTransferItemRequest
		itemsRes     []*bridgeService.ItemTransferItem
	)

	for _, item := range req.ItemTransferItems {
		items = append(items, &dto.CreateItemTransferItemRequest{
			ItemID:      item.ItemId,
			TransferQty: item.TransferQty,
			RequestQty:  item.RequestQty,
			ReceiveQty:  item.ReceiveQty,
			ReceiveNote: item.ReceiveNote,
			UnitCost:    item.UnitCost,
			Note:        item.Note,
		})
	}

	itemTransfer, err = h.ServicesItemTransfer.Create(ctx, &dto.CreateItemTransferRequest{
		RequestDateStr:    req.RequestDateStr,
		SiteOriginID:      req.SiteOriginId,
		SiteDestinationID: req.SiteDestinationId,
		StockTypeID:       int8(req.StockTypeId),
		Note:              req.Note,
		ItemTransferItems: items,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	for _, item := range itemTransfer.ItemTransferItem {
		itemsRes = append(itemsRes, &bridgeService.ItemTransferItem{
			Id:          item.ID,
			DeliverQty:  item.DeliverQty,
			ReceiveQty:  item.ReceiveQty,
			RequestQty:  item.RequestQty,
			ReceiveNote: item.ReceiveNote,
			UnitCost:    item.UnitCost,
			Subtotal:    item.Subtotal,
			Weight:      item.Weight,
			Note:        item.Note,
		})
	}

	res = &bridgeService.GetItemTransferDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.ItemTransfer{
			Id:                 itemTransfer.ID,
			Code:               itemTransfer.Code,
			RequestDate:        timestamppb.New(itemTransfer.RequestDate),
			RecognitionDate:    timestamppb.New(itemTransfer.RecognitionDate),
			EtaDate:            timestamppb.New(itemTransfer.EtaDate),
			EtaTime:            itemTransfer.EtaTime,
			AtaDate:            timestamppb.New(itemTransfer.AtaDate),
			AtaTime:            itemTransfer.AtaTime,
			AdditionalCost:     itemTransfer.AdditionalCost,
			AdditionalCostNote: itemTransfer.AdditionalCostNote,
			StockType:          int32(itemTransfer.StockType),
			TotalCost:          itemTransfer.TotalCost,
			TotalSku:           itemTransfer.TotalSku,
			TotalCharge:        itemTransfer.TotalCharge,
			TotalWeight:        itemTransfer.TotalWeight,
			Note:               itemTransfer.Note,
			Status:             int32(itemTransfer.Status),
			UpdatedAt:          timestamppb.New(itemTransfer.UpdatedAt),
			UpdatedBy:          itemTransfer.UpdatedBy,
			Locked:             int32(itemTransfer.Locked),
			LockedBy:           itemTransfer.LockedBy,
			ItemTransferItem:   itemsRes,
		},
	}
	return
}

func (h *BridgeGrpcHandler) UpdateItemTransfer(ctx context.Context, req *bridgeService.UpdateItemTransferRequest) (res *bridgeService.GetItemTransferDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.UpdatePurchaseOrder")
	defer span.End()

	var (
		itemTransfer dto.ItemTransferResponse
		items        []*dto.UpdateItemTransferItemRequest
		itemsRes     []*bridgeService.ItemTransferItem
	)

	for _, item := range req.ItemTransferItems {
		items = append(items, &dto.UpdateItemTransferItemRequest{
			Id:          item.Id,
			ItemID:      item.ItemId,
			TransferQty: item.TransferQty,
			RequestQty:  item.RequestQty,
			ReceiveQty:  item.ReceiveQty,
			ReceiveNote: item.ReceiveNote,
			UnitCost:    item.UnitCost,
			Note:        item.Note,
		})
	}

	itemTransfer, err = h.ServicesItemTransfer.Update(ctx, &dto.UpdateItemTransferRequest{
		Id:                 req.Id,
		RecognitionDateStr: req.RecognitionDateStr,
		EtaDateStr:         req.EtaDateStr,
		EtaTimeStr:         req.EtaTimeStr,
		AdditionalCost:     req.AdditionalCost,
		AdditionalCostNote: req.AdditionalCostNote,
		RequestDateStr:     req.RequestDateStr,
		SiteOriginID:       req.SiteOriginId,
		SiteDestinationID:  req.SiteDestinationId,
		Note:               req.Note,
		ItemTransferItems:  items,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	for _, item := range itemTransfer.ItemTransferItem {
		itemsRes = append(itemsRes, &bridgeService.ItemTransferItem{
			Id:          item.ID,
			DeliverQty:  item.DeliverQty,
			ReceiveQty:  item.ReceiveQty,
			RequestQty:  item.RequestQty,
			ReceiveNote: item.ReceiveNote,
			UnitCost:    item.UnitCost,
			Subtotal:    item.Subtotal,
			Weight:      item.Weight,
			Note:        item.Note,
		})
	}

	res = &bridgeService.GetItemTransferDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.ItemTransfer{
			Id:                 itemTransfer.ID,
			Code:               itemTransfer.Code,
			RequestDate:        timestamppb.New(itemTransfer.RequestDate),
			RecognitionDate:    timestamppb.New(itemTransfer.RecognitionDate),
			EtaDate:            timestamppb.New(itemTransfer.EtaDate),
			EtaTime:            itemTransfer.EtaTime,
			AtaDate:            timestamppb.New(itemTransfer.AtaDate),
			AtaTime:            itemTransfer.AtaTime,
			AdditionalCost:     itemTransfer.AdditionalCost,
			AdditionalCostNote: itemTransfer.AdditionalCostNote,
			StockType:          int32(itemTransfer.StockType),
			TotalCost:          itemTransfer.TotalCost,
			TotalSku:           itemTransfer.TotalSku,
			TotalCharge:        itemTransfer.TotalCharge,
			TotalWeight:        itemTransfer.TotalWeight,
			Note:               itemTransfer.Note,
			Status:             int32(itemTransfer.Status),
			UpdatedAt:          timestamppb.New(itemTransfer.UpdatedAt),
			UpdatedBy:          itemTransfer.UpdatedBy,
			Locked:             int32(itemTransfer.Locked),
			LockedBy:           itemTransfer.LockedBy,
			ItemTransferItem:   itemsRes,
		},
	}
	return
}

func (h *BridgeGrpcHandler) CommitItemTransfer(ctx context.Context, req *bridgeService.CommitItemTransferRequest) (res *bridgeService.GetItemTransferDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CommitPurchaseOrder")
	defer span.End()

	var (
		itemTransfer dto.ItemTransferResponse
		items        []*dto.UpdateItemTransferItemRequest
		itemsRes     []*bridgeService.ItemTransferItem
	)

	for _, item := range req.ItemTransferItems {
		items = append(items, &dto.UpdateItemTransferItemRequest{
			Id:          item.Id,
			ItemID:      item.ItemId,
			TransferQty: item.TransferQty,
			RequestQty:  item.RequestQty,
			ReceiveQty:  item.ReceiveQty,
			ReceiveNote: item.ReceiveNote,
			UnitCost:    item.UnitCost,
			Note:        item.Note,
		})
	}

	itemTransfer, err = h.ServicesItemTransfer.Commit(ctx, &dto.CommitItemTransferRequest{
		Id:                 req.Id,
		RecognitionDateStr: req.RecognitionDate,
		EtaDateStr:         req.EtaDate,
		EtaTimeStr:         req.EtaTime,
		AdditionalCost:     req.AdditionalCost,
		AdditionalCostNote: req.AdditionalCostNote,
		ItemTransferItems:  items,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	for _, item := range itemTransfer.ItemTransferItem {
		itemsRes = append(itemsRes, &bridgeService.ItemTransferItem{
			Id:          item.ID,
			DeliverQty:  item.DeliverQty,
			ReceiveQty:  item.ReceiveQty,
			RequestQty:  item.RequestQty,
			ReceiveNote: item.ReceiveNote,
			UnitCost:    item.UnitCost,
			Subtotal:    item.Subtotal,
			Weight:      item.Weight,
			Note:        item.Note,
		})
	}

	res = &bridgeService.GetItemTransferDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.ItemTransfer{
			Id:                 itemTransfer.ID,
			Code:               itemTransfer.Code,
			RequestDate:        timestamppb.New(itemTransfer.RequestDate),
			RecognitionDate:    timestamppb.New(itemTransfer.RecognitionDate),
			EtaDate:            timestamppb.New(itemTransfer.EtaDate),
			EtaTime:            itemTransfer.EtaTime,
			AtaDate:            timestamppb.New(itemTransfer.AtaDate),
			AtaTime:            itemTransfer.AtaTime,
			AdditionalCost:     itemTransfer.AdditionalCost,
			AdditionalCostNote: itemTransfer.AdditionalCostNote,
			StockType:          int32(itemTransfer.StockType),
			TotalCost:          itemTransfer.TotalCost,
			TotalSku:           itemTransfer.TotalSku,
			TotalCharge:        itemTransfer.TotalCharge,
			TotalWeight:        itemTransfer.TotalWeight,
			Note:               itemTransfer.Note,
			Status:             int32(itemTransfer.Status),
			UpdatedAt:          timestamppb.New(itemTransfer.UpdatedAt),
			UpdatedBy:          itemTransfer.UpdatedBy,
			Locked:             int32(itemTransfer.Locked),
			LockedBy:           itemTransfer.LockedBy,
			ItemTransferItem:   itemsRes,
		},
	}
	return
}

func (h *BridgeGrpcHandler) GetInTransitTransferGPList(ctx context.Context, req *bridgeService.GetInTransitTransferGPListRequest) (res *bridgeService.GetInTransitTransferGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetInTransitTransferGPList")
	defer span.End()

	res, err = h.ServicesItemTransfer.GetInTransitTransferGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) GetInTransitTransferGPDetail(ctx context.Context, req *bridgeService.GetInTransitTransferGPDetailRequest) (res *bridgeService.GetInTransitTransferGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetInTransitTransferGPDetail")
	defer span.End()

	res, err = h.ServicesItemTransfer.GetInTransitTransferDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) GetTransferRequestGPList(ctx context.Context, req *bridgeService.GetTransferRequestGPListRequest) (res *bridgeService.GetTransferRequestGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetTransferRequestGPList")
	defer span.End()

	res, err = h.ServicesItemTransfer.GetTransferRequestGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) GetTransferRequestGPDetail(ctx context.Context, req *bridgeService.GetTransferRequestGPDetailRequest) (res *bridgeService.GetTransferRequestGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetTransferRequestGPDetail")
	defer span.End()

	res, err = h.ServicesItemTransfer.GetTransferRequestDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) CreateTransferRequestGP(ctx context.Context, req *bridgeService.CreateTransferRequestGPRequest) (res *bridgeService.CreateTransferRequestGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateTransferRequestGP")
	defer span.End()

	var (
		bodyReq dto.CreateTransferRequestGPRequest
		details []dto.CreateTransferRequestDetailGPRequest
	)

	for _, detail := range req.Detail {
		details = append(details, dto.CreateTransferRequestDetailGPRequest{
			Lnitmseq:      int(detail.Lnitmseq),
			Itemnmbr:      detail.Itemnmbr,
			Uofm:          detail.Uofm,
			IvmQtyRequest: detail.IvmQtyRequest,
			IvmQtyFulfill: detail.IvmQtyFulfill,
		})
	}
	bodyReq = dto.CreateTransferRequestGPRequest{
		Interid:         global.EnvDatabaseGP,
		Docnumbr:        req.Docnumbr,
		Docdate:         req.Docdate,
		IvmTrType:       int(req.IvmTrType),
		RequestDate:     req.RequestDate,
		IvmReqEta:       req.IvmReqEta,
		IvmLocncodeFrom: req.IvmLocncodeFrom,
		IvmLocncodeTo:   req.IvmLocncodeTo,
		ReasonCode:      req.Reason_Code,
		Detail:          details,
	}

	res, err = h.ServicesItemTransfer.CreateGP(ctx, &bodyReq)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) UpdateTransferRequestGP(ctx context.Context, req *bridgeService.UpdateTransferRequestGPRequest) (res *bridgeService.CreateTransferRequestGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.UpdateTransferRequestGP")
	defer span.End()

	var (
		bodyReq dto.UpdateTransferRequestGPRequest
		details []dto.UpdateTransferRequestDetailGPRequest
	)

	for _, detail := range req.Detail {
		details = append(details, dto.UpdateTransferRequestDetailGPRequest{
			Lnitmseq:      int(detail.Lnitmseq),
			Itemnmbr:      detail.Itemnmbr,
			IvmQtyRequest: detail.IvmQtyRequest,
		})
	}
	bodyReq = dto.UpdateTransferRequestGPRequest{
		Interid:     global.EnvDatabaseGP,
		Docnumbr:    req.Docnumbr,
		Docdate:     req.Docdate,
		RequestDate: req.RequestDate,
		IvmReqEta:   req.IvmReqEta,
		ReasonCode:  req.Reason_Code,
		Detail:      details,
	}

	res, err = h.ServicesItemTransfer.UpdateGP(ctx, &bodyReq)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) UpdateInTransitTransferGP(ctx context.Context, req *bridgeService.UpdateInTransitTransferGPRequest) (res *bridgeService.UpdateInTransitTransferGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.UpdateInTransitTransferGP")
	defer span.End()

	var (
		bodyReq dto.UpdateInTransiteTransferGPRequest
		details []dto.UpdateInTransitTransferDetailGPRequest
	)

	for _, detail := range req.Detail {
		details = append(details, dto.UpdateInTransitTransferDetailGPRequest{
			Lnitmseq:   int(detail.Lnitmseq),
			Itemnmbr:   detail.Itemnmbr,
			ReasonCode: detail.ReasonCode,
			Qtyfulfi:   detail.Qtyfulfi,
		})
	}
	bodyReq = dto.UpdateInTransiteTransferGPRequest{
		Interid:     global.EnvDatabaseGP,
		Orddocid:    req.Orddocid,
		IvmTrNumber: req.IvmTrNumber,
		Ordrdate:    req.Ordrdate,
		Etadte:      req.Etadte,
		Etatime:     req.Eta,
		Note:        req.Note,
		Detail:      details,
	}

	res, err = h.ServicesItemTransfer.UpdateITTGP(ctx, &bodyReq)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) CommitTransferRequestGP(ctx context.Context, req *bridgeService.CommitTransferRequestGPRequest) (res *bridgeService.CommitTransferRequestGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateTransferRequestGP")
	defer span.End()

	var (
		bodyReq dto.CommitTransferRequestGPRequest
		details []dto.CommitTransferRequestDetailGPRequest
	)

	for _, detail := range req.Detail {
		details = append(details, dto.CommitTransferRequestDetailGPRequest{
			Lnitmseq:      int(detail.Lnitmseq),
			IvmQtyFulfill: detail.IvmQtyFulfill,
		})
	}
	bodyReq = dto.CommitTransferRequestGPRequest{
		Interid:  global.EnvDatabaseGP,
		Docnumbr: req.Docnumbr,
		ItLocn:   "WRHINT0001", // currently hardcoded, because via site gonna be this always
		Detail:   details,
	}

	res, err = h.ServicesItemTransfer.CommitGP(ctx, &bodyReq)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}
