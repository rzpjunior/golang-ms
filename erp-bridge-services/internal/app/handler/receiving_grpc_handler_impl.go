package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *BridgeGrpcHandler) GetReceivingList(ctx context.Context, req *bridgeService.GetReceivingListRequest) (res *bridgeService.GetReceivingListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetReceivingList")
	defer span.End()

	var grs []dto.ReceivingResponse
	grs, _, err = h.ServiceReceiving.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.Receiving
	for _, gr := range grs {
		data = append(data, &bridgeService.Receiving{
			Id:                  gr.ID,
			Code:                gr.Code,
			SiteId:              gr.SiteId,
			PurchaseOrderId:     gr.PurchaseOrderId,
			ItemTransferId:      gr.ItemTransferId,
			InboundType:         int32(gr.InboundType),
			ValidSupplierReturn: int32(gr.ValidSupplierReturn),
			AtaDate:             timestamppb.New(gr.AtaDate),
			AtaTime:             gr.AtaTime,
			StockType:           int32(gr.StockType),
			TotalWeight:         gr.TotalWeight,
			Note:                gr.Note,
			Status:              int32(gr.Status),
			CreatedAt:           timestamppb.New(gr.CreatedAt),
			CreatedBy:           gr.CreatedBy,
			ConfirmedAt:         timestamppb.New(gr.ConfirmedAt),
			ConfirmedBy:         gr.ConfirmedBy,
			UpdatedAt:           timestamppb.New(gr.UpdatedAt),
			UpdatedBy:           gr.UpdatedBy,
			Locked:              int32(gr.Locked),
			LockedBy:            gr.LockedBy,
		})
	}

	res = &bridgeService.GetReceivingListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetReceivingDetail(ctx context.Context, req *bridgeService.GetReceivingDetailRequest) (res *bridgeService.GetReceivingDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemTransferDetail")
	defer span.End()

	var (
		gr    dto.ReceivingResponse
		items []*bridgeService.ReceivingItem
		po    *bridgeService.PurchaseOrder
		it    *bridgeService.ItemTransfer
	)

	gr, err = h.ServiceReceiving.GetDetail(ctx, req.Id, "")
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	if gr.InboundType == 1 {
		po = &bridgeService.PurchaseOrder{
			Id:                     gr.PurchaseOrder.ID,
			Code:                   gr.PurchaseOrder.Code,
			VendorId:               gr.PurchaseOrder.VendorID,
			SiteId:                 gr.PurchaseOrder.SiteID,
			TermPaymentPurId:       gr.PurchaseOrder.TermPaymentPurID,
			VendorClassificationId: gr.PurchaseOrder.VendorClassificationID,
			PurchasePlanId:         gr.PurchaseOrder.PurchasePlanID,
			ConsolidatedShipmentId: gr.PurchaseOrder.ConsolidatedShipmentID,
			Status:                 gr.PurchaseOrder.Status,
			RecognitionDate:        timestamppb.New(gr.PurchaseOrder.RecognitionDate),
			EtaDate:                timestamppb.New(gr.PurchaseOrder.EtaDate),
			SiteAddress:            gr.PurchaseOrder.SiteAddress,
			EtaTime:                gr.PurchaseOrder.EtaTime,
			TaxPct:                 gr.PurchaseOrder.TaxPct,
			DeliveryFee:            gr.PurchaseOrder.DeliveryFee,
			TotalPrice:             gr.PurchaseOrder.TotalPrice,
			TaxAmount:              gr.PurchaseOrder.TaxAmount,
			TotalCharge:            gr.PurchaseOrder.TotalCharge,
			TotalInvoice:           gr.PurchaseOrder.TotalInvoice,
			TotalWeight:            gr.PurchaseOrder.TotalWeight,
			Note:                   gr.PurchaseOrder.Note,
			DeltaPrint:             gr.PurchaseOrder.DeltaPrint,
			Latitude:               gr.PurchaseOrder.Latitude,
			Longitude:              gr.PurchaseOrder.Longitude,
			UpdatedAt:              timestamppb.New(gr.PurchaseOrder.UpdatedAt),
			UpdatedBy:              gr.PurchaseOrder.UpdatedBy,
			CreatedFrom:            gr.PurchaseOrder.CreatedFrom,
			HasFinishedGr:          int32(gr.PurchaseOrder.HasFinishedGr),
			CreatedAt:              timestamppb.New(gr.PurchaseOrder.CreatedAt),
			CreatedBy:              gr.PurchaseOrder.CreatedBy,
			CommittedAt:            timestamppb.New(gr.PurchaseOrder.CommittedAt),
			CommittedBy:            gr.PurchaseOrder.CommittedBy,
			AssignedTo:             gr.PurchaseOrder.AssignedTo,
			AssignedBy:             gr.PurchaseOrder.AssignedBy,
			AssignedAt:             timestamppb.New(gr.PurchaseOrder.AssignedAt),
			Locked:                 gr.PurchaseOrder.Locked,
			LockedBy:               gr.PurchaseOrder.LockedBy,
		}
	} else if gr.InboundType == 2 {
		it = &bridgeService.ItemTransfer{
			Id:                 gr.ItemTransfer.ID,
			Code:               gr.ItemTransfer.Code,
			RequestDate:        timestamppb.New(gr.ItemTransfer.RequestDate),
			RecognitionDate:    timestamppb.New(gr.ItemTransfer.RecognitionDate),
			EtaDate:            timestamppb.New(gr.ItemTransfer.EtaDate),
			EtaTime:            gr.ItemTransfer.EtaTime,
			AtaDate:            timestamppb.New(gr.ItemTransfer.AtaDate),
			AtaTime:            gr.ItemTransfer.AtaTime,
			AdditionalCost:     gr.ItemTransfer.AdditionalCost,
			AdditionalCostNote: gr.ItemTransfer.AdditionalCostNote,
			StockType:          int32(gr.ItemTransfer.StockType),
			TotalCost:          gr.ItemTransfer.TotalCost,
			TotalSku:           gr.ItemTransfer.TotalSku,
			TotalCharge:        gr.ItemTransfer.TotalCharge,
			TotalWeight:        gr.ItemTransfer.TotalWeight,
			Note:               gr.ItemTransfer.Note,
			Status:             int32(gr.ItemTransfer.Status),
			UpdatedAt:          timestamppb.New(gr.ItemTransfer.UpdatedAt),
			UpdatedBy:          gr.ItemTransfer.UpdatedBy,
			Locked:             int32(gr.ItemTransfer.Locked),
			LockedBy:           gr.ItemTransfer.LockedBy,
		}
	}

	for _, item := range gr.ReceivingItems {
		items = append(items, &bridgeService.ReceivingItem{
			Id:                  item.ID,
			PurchaseOrderItemId: item.PurchaseOrderItemID,
			ItemTransferItemId:  item.ItemTransferItemID,
			RejectQty:           item.RejectQty,
			RejectReason:        int32(item.RejectReason),
			IsDisabled:          int32(item.IsDisabled),
			DeliverQty:          item.DeliverQty,
			ReceiveQty:          item.ReceiveQty,
			Weight:              item.Weight,
			Note:                item.Note,
			Item: &bridgeService.Item{
				Id:                      item.Item.ID,
				Code:                    item.Item.Code,
				UomId:                   item.Item.UomID,
				ClassId:                 item.Item.ClassID,
				ItemCategoryId:          item.Item.ItemCategoryID,
				Description:             item.Item.Description,
				UnitWeightConversion:    item.Item.UnitWeightConversion,
				OrderMinQty:             item.Item.OrderMinQty,
				OrderMaxQty:             item.Item.OrderMaxQty,
				ItemType:                item.Item.ItemType,
				Packability:             item.Item.Packability,
				Capitalize:              item.Item.Capitalize,
				ExcludeArchetype:        item.Item.ExcludeArchetype,
				MaxDayDeliveryDate:      int32(item.Item.MaxDayDeliveryDate),
				FragileGoods:            item.Item.FragileGoods,
				Taxable:                 item.Item.Taxable,
				OrderChannelRestriction: item.Item.OrderChannelRestriction,
				Note:                    item.Item.Note,
				Status:                  int32(item.Item.Status),
				CreatedAt:               timestamppb.New(item.Item.CreatedAt),
				UpdatedAt:               timestamppb.New(item.Item.UpdatedAt),
				Uom: &bridgeService.Uom{
					Id:             item.Item.Uom.ID,
					Code:           item.Item.Uom.Code,
					Description:    item.Item.Uom.Description,
					Status:         int32(item.Item.Uom.Status),
					DecimalEnabled: int32(item.Item.Uom.DecimalEnabled),
					CreatedAt:      timestamppb.New(item.Item.Uom.CreatedAt),
					UpdatedAt:      timestamppb.New(item.Item.Uom.UpdatedAt),
				},
			},
		})
	}

	res = &bridgeService.GetReceivingDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.Receiving{
			Id:                  gr.ID,
			Code:                gr.Code,
			SiteId:              gr.SiteId,
			PurchaseOrderId:     gr.PurchaseOrderId,
			InboundType:         int32(gr.InboundType),
			ValidSupplierReturn: int32(gr.ValidSupplierReturn),
			AtaDate:             timestamppb.New(gr.AtaDate),
			AtaTime:             gr.AtaTime,
			StockType:           int32(gr.StockType),
			TotalWeight:         gr.TotalWeight,
			Note:                gr.Note,
			Status:              int32(gr.Status),
			CreatedAt:           timestamppb.New(gr.CreatedAt),
			CreatedBy:           gr.CreatedBy,
			ConfirmedAt:         timestamppb.New(gr.ConfirmedAt),
			ConfirmedBy:         gr.ConfirmedBy,
			UpdatedAt:           timestamppb.New(gr.UpdatedAt),
			UpdatedBy:           gr.UpdatedBy,
			Locked:              int32(gr.Locked),
			LockedBy:            gr.LockedBy,
			ReceivingItems:      items,
			PurchaseOrder:       po,
			ItemTransfer:        it,
		},
	}
	return
}

func (h *BridgeGrpcHandler) CreateReceiving(ctx context.Context, req *bridgeService.CreateReceivingRequest) (res *bridgeService.GetReceivingDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateItemTransfer")
	defer span.End()

	var (
		gr       dto.ReceivingResponse
		items    []*dto.CreateReceivingItemRequest
		itemsRes []*bridgeService.ReceivingItem
	)

	for _, item := range req.ReceivingItem {
		items = append(items, &dto.CreateReceivingItemRequest{
			ReceiveQty: item.ReceiveQty,
			Note:       item.Note,
		})
	}

	gr, err = h.ServiceReceiving.Create(ctx, &dto.CreateReceivingRequest{

		Note:          req.Note,
		ReceivingItem: items,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	for _, item := range gr.ReceivingItems {
		itemsRes = append(itemsRes, &bridgeService.ReceivingItem{
			Id:           item.ID,
			DeliverQty:   item.DeliverQty,
			ReceiveQty:   item.ReceiveQty,
			RejectQty:    item.RejectQty,
			RejectReason: int32(item.RejectReason),
			IsDisabled:   int32(item.IsDisabled),
			Weight:       item.Weight,
			Note:         item.Note,
		})
	}

	res = &bridgeService.GetReceivingDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.Receiving{
			Id:                  gr.ID,
			Code:                gr.Code,
			SiteId:              gr.SiteId,
			PurchaseOrderId:     gr.PurchaseOrderId,
			InboundType:         int32(gr.InboundType),
			ValidSupplierReturn: int32(gr.ValidSupplierReturn),
			AtaDate:             timestamppb.New(gr.AtaDate),
			AtaTime:             gr.AtaTime,
			StockType:           int32(gr.StockType),
			TotalWeight:         gr.TotalWeight,
			Note:                gr.Note,
			Status:              int32(gr.Status),
			CreatedAt:           timestamppb.New(gr.CreatedAt),
			CreatedBy:           gr.CreatedBy,
			ConfirmedAt:         timestamppb.New(gr.ConfirmedAt),
			ConfirmedBy:         gr.ConfirmedBy,
			UpdatedAt:           timestamppb.New(gr.UpdatedAt),
			UpdatedBy:           gr.UpdatedBy,
			Locked:              int32(gr.Locked),
			LockedBy:            gr.LockedBy,
			ReceivingItems:      itemsRes,
		},
	}
	return
}

func (h *BridgeGrpcHandler) ConfirmReceiving(ctx context.Context, req *bridgeService.ConfirmReceivingRequest) (res *bridgeService.GetReceivingDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.ConfirmItemTransfer")
	defer span.End()

	var (
		gr       dto.ReceivingResponse
		itemsRes []*bridgeService.ReceivingItem
	)

	gr, err = h.ServiceReceiving.Confirm(ctx, &dto.ConfirmReceivingRequest{
		InboundType: req.InboundType,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	for _, item := range gr.ReceivingItems {
		itemsRes = append(itemsRes, &bridgeService.ReceivingItem{
			Id:           item.ID,
			DeliverQty:   item.DeliverQty,
			ReceiveQty:   item.ReceiveQty,
			RejectQty:    item.RejectQty,
			RejectReason: int32(item.RejectReason),
			IsDisabled:   int32(item.IsDisabled),
			Weight:       item.Weight,
			Note:         item.Note,
		})
	}

	res = &bridgeService.GetReceivingDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.Receiving{
			Id:                  gr.ID,
			Code:                gr.Code,
			SiteId:              gr.SiteId,
			PurchaseOrderId:     gr.PurchaseOrderId,
			InboundType:         int32(gr.InboundType),
			ValidSupplierReturn: int32(gr.ValidSupplierReturn),
			AtaDate:             timestamppb.New(gr.AtaDate),
			AtaTime:             gr.AtaTime,
			StockType:           int32(gr.StockType),
			TotalWeight:         gr.TotalWeight,
			Note:                gr.Note,
			Status:              int32(gr.Status),
			CreatedAt:           timestamppb.New(gr.CreatedAt),
			CreatedBy:           gr.CreatedBy,
			ConfirmedAt:         timestamppb.New(gr.ConfirmedAt),
			ConfirmedBy:         gr.ConfirmedBy,
			UpdatedAt:           timestamppb.New(gr.UpdatedAt),
			UpdatedBy:           gr.UpdatedBy,
			Locked:              int32(gr.Locked),
			LockedBy:            gr.LockedBy,
			ReceivingItems:      itemsRes,
		},
	}
	return
}

func (h *BridgeGrpcHandler) GetGoodsReceiptGPList(ctx context.Context, req *bridgeService.GetGoodsReceiptGPListRequest) (res *bridgeService.GetGoodsReceiptGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetGoodsReceiptGPList")
	defer span.End()

	res, err = h.ServiceReceiving.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) GetGoodsReceiptGPDetail(ctx context.Context, req *bridgeService.GetGoodsReceiptGPDetailRequest) (res *bridgeService.GetGoodsReceiptGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetGoodsReceiptGPDetail")
	defer span.End()

	res, err = h.ServiceReceiving.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) CreateGoodsReceiptGP(ctx context.Context, req *bridgeService.CreateGoodsReceiptGPRequest) (res *bridgeService.CreateTransferRequestGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateGoodsReceiptGP")
	defer span.End()

	var resAPi *bridgeService.CreateTransferRequestGPResponse
	resAPi, err = h.ServiceReceiving.CreateGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.CreateTransferRequestGPResponse{
		Code:    resAPi.Code,
		Message: resAPi.Message,
	}
	return
}

func (h *BridgeGrpcHandler) UpdateGoodsReceiptGP(ctx context.Context, req *bridgeService.UpdateGoodsReceiptGPRequest) (res *bridgeService.CreateTransferRequestGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.UpdateGoodsReceiptGP")
	defer span.End()

	var resAPi *bridgeService.CreateTransferRequestGPResponse
	resAPi, err = h.ServiceReceiving.UpdateGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.CreateTransferRequestGPResponse{
		Code:    resAPi.Code,
		Message: resAPi.Message,
	}
	return
}
