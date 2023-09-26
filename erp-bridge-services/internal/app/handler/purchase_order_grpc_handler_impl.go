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

func (h *BridgeGrpcHandler) GetPurchaseOrderList(ctx context.Context, req *bridgeService.GetPurchaseOrderListRequest) (res *bridgeService.GetPurchaseOrderListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPurchaseOrderList")
	defer span.End()

	var purchaseOrders []dto.PurchaseOrderResponse
	purchaseOrders, _, err = h.ServicesPurchaseOrder.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.PurchaseOrder
	for _, purchaseOrder := range purchaseOrders {
		var grs []*bridgeService.ReceivingListinDetailResponse
		for _, gr := range purchaseOrder.Receiving {
			grs = append(grs, &bridgeService.ReceivingListinDetailResponse{
				Id:     gr.ID,
				Code:   gr.Code,
				Status: int32(gr.Status),
			})
		}

		data = append(data, &bridgeService.PurchaseOrder{
			Id:                     purchaseOrder.ID,
			Code:                   purchaseOrder.Code,
			VendorId:               purchaseOrder.VendorID,
			SiteId:                 purchaseOrder.SiteID,
			TermPaymentPurId:       purchaseOrder.TermPaymentPurID,
			VendorClassificationId: purchaseOrder.VendorClassificationID,
			PurchasePlanId:         purchaseOrder.PurchasePlanID,
			ConsolidatedShipmentId: purchaseOrder.ConsolidatedShipmentID,
			Status:                 purchaseOrder.Status,
			RecognitionDate:        timestamppb.New(purchaseOrder.RecognitionDate),
			EtaDate:                timestamppb.New(purchaseOrder.EtaDate),
			SiteAddress:            purchaseOrder.SiteAddress,
			EtaTime:                purchaseOrder.EtaTime,
			TaxPct:                 purchaseOrder.TaxPct,
			DeliveryFee:            purchaseOrder.DeliveryFee,
			TotalPrice:             purchaseOrder.TotalPrice,
			TaxAmount:              purchaseOrder.TaxAmount,
			TotalCharge:            purchaseOrder.TotalCharge,
			TotalInvoice:           purchaseOrder.TotalInvoice,
			TotalWeight:            purchaseOrder.TotalWeight,
			Note:                   purchaseOrder.Note,
			DeltaPrint:             purchaseOrder.DeltaPrint,
			Latitude:               purchaseOrder.Latitude,
			Longitude:              purchaseOrder.Longitude,
			UpdatedAt:              timestamppb.New(purchaseOrder.UpdatedAt),
			UpdatedBy:              purchaseOrder.UpdatedBy,
			CreatedFrom:            purchaseOrder.CreatedFrom,
			HasFinishedGr:          int32(purchaseOrder.HasFinishedGr),
			Receiving:              grs,
			CreatedAt:              timestamppb.New(purchaseOrder.CreatedAt),
			CreatedBy:              purchaseOrder.CreatedBy,
			CommittedAt:            timestamppb.New(purchaseOrder.CommittedAt),
			CommittedBy:            purchaseOrder.CommittedBy,
			AssignedTo:             purchaseOrder.AssignedTo,
			AssignedBy:             purchaseOrder.AssignedBy,
			AssignedAt:             timestamppb.New(purchaseOrder.AssignedAt),
			Locked:                 purchaseOrder.Locked,
			LockedBy:               purchaseOrder.LockedBy,
		})
	}

	res = &bridgeService.GetPurchaseOrderListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetPurchaseOrderDetail(ctx context.Context, req *bridgeService.GetPurchaseOrderDetailRequest) (res *bridgeService.GetPurchaseOrderDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPurchaseOrderDetail")
	defer span.End()

	var (
		purchaseOrder dto.PurchaseOrderResponse
		items         []*bridgeService.PurchaseOrderItem
		grs           []*bridgeService.ReceivingListinDetailResponse
	)
	purchaseOrder, err = h.ServicesPurchaseOrder.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	for _, item := range purchaseOrder.PurchaseOrderItems {
		items = append(items, &bridgeService.PurchaseOrderItem{
			Id:                 item.ID,
			PurchaseOrderId:    item.PurchaseOrderID,
			PurchasePlanItemId: item.PurchasePlanItemID,
			ItemId:             item.ItemID,
			OrderQty:           item.OrderQty,
			UnitPrice:          item.UnitPrice,
			TaxableItem:        item.TaxableItem,
			IncludeTax:         item.IncludeTax,
			TaxPercentage:      item.TaxPercentage,
			TaxAmount:          item.TaxAmount,
			UnitPriceTax:       item.UnitPriceTax,
			Subtotal:           item.Subtotal,
			Weight:             item.Weight,
			Note:               item.Note,
			PurchaseQty:        item.PurchaseQty,
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

	for _, gr := range purchaseOrder.Receiving {
		grs = append(grs, &bridgeService.ReceivingListinDetailResponse{
			Id:     gr.ID,
			Code:   gr.Code,
			Status: int32(gr.Status),
		})
	}

	res = &bridgeService.GetPurchaseOrderDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.PurchaseOrder{
			Id:                     purchaseOrder.ID,
			Code:                   purchaseOrder.Code,
			VendorId:               purchaseOrder.VendorID,
			SiteId:                 purchaseOrder.SiteID,
			TermPaymentPurId:       purchaseOrder.TermPaymentPurID,
			VendorClassificationId: purchaseOrder.VendorClassificationID,
			PurchasePlanId:         purchaseOrder.PurchasePlanID,
			ConsolidatedShipmentId: purchaseOrder.ConsolidatedShipmentID,
			Status:                 purchaseOrder.Status,
			RecognitionDate:        timestamppb.New(purchaseOrder.RecognitionDate),
			EtaDate:                timestamppb.New(purchaseOrder.EtaDate),
			SiteAddress:            purchaseOrder.SiteAddress,
			EtaTime:                purchaseOrder.EtaTime,
			TaxPct:                 purchaseOrder.TaxPct,
			DeliveryFee:            purchaseOrder.DeliveryFee,
			TotalPrice:             purchaseOrder.TotalPrice,
			TaxAmount:              purchaseOrder.TaxAmount,
			TotalCharge:            purchaseOrder.TotalCharge,
			TotalInvoice:           purchaseOrder.TotalInvoice,
			TotalWeight:            purchaseOrder.TotalWeight,
			Note:                   purchaseOrder.Note,
			DeltaPrint:             purchaseOrder.DeltaPrint,
			Latitude:               purchaseOrder.Latitude,
			Longitude:              purchaseOrder.Longitude,
			UpdatedAt:              timestamppb.New(purchaseOrder.UpdatedAt),
			UpdatedBy:              purchaseOrder.UpdatedBy,
			CreatedFrom:            purchaseOrder.CreatedFrom,
			HasFinishedGr:          int32(purchaseOrder.HasFinishedGr),
			CreatedAt:              timestamppb.New(purchaseOrder.CreatedAt),
			CreatedBy:              purchaseOrder.CreatedBy,
			CommittedAt:            timestamppb.New(purchaseOrder.CommittedAt),
			CommittedBy:            purchaseOrder.CommittedBy,
			AssignedTo:             purchaseOrder.AssignedTo,
			AssignedBy:             purchaseOrder.AssignedBy,
			AssignedAt:             timestamppb.New(purchaseOrder.AssignedAt),
			Locked:                 purchaseOrder.Locked,
			LockedBy:               purchaseOrder.LockedBy,
			PurchaseOrderItems:     items,
			Receiving:              grs,
		},
	}
	return
}

func (h *BridgeGrpcHandler) CreatePurchaseOrder(ctx context.Context, req *bridgeService.CreatePurchaseOrderRequest) (res *bridgeService.GetPurchaseOrderDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreatePurchaseOrder")
	defer span.End()

	var (
		purchaseOrder dto.PurchaseOrderResponse
		items         []*dto.CreatePurchaseOrderItemRequest
		itemsRes      []*bridgeService.PurchaseOrderItem
	)

	for _, item := range req.Items {
		items = append(items, &dto.CreatePurchaseOrderItemRequest{
			ItemID:        item.ItemId,
			OrderQty:      item.OrderQty,
			UnitPrice:     item.UnitPrice,
			Note:          item.Note,
			PurchaseQty:   item.PurchaseQty,
			IncludeTax:    int8(item.IncludeTax),
			TaxPercentage: item.TaxPercentage,
		})
	}

	purchaseOrder, err = h.ServicesPurchaseOrder.Create(ctx, &dto.CreatePurchaseOrderRequest{
		VendorID:           req.VendorId,
		SiteID:             req.SiteId,
		OrderDate:          req.OrderDate,
		StrEtaDate:         req.StrEtaDate,
		EtaTime:            req.EtaTime,
		DeliveryFee:        req.DeliveryFee,
		Note:               req.Note,
		TaxPct:             req.TaxPct,
		Latitude:           req.Latitude,
		Longitude:          req.Longitude,
		PurchaseOrderItems: items,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	for _, item := range purchaseOrder.PurchaseOrderItems {
		itemsRes = append(itemsRes, &bridgeService.PurchaseOrderItem{
			Id:                 item.ID,
			PurchaseOrderId:    item.PurchaseOrderID,
			PurchasePlanItemId: item.PurchasePlanItemID,
			ItemId:             item.ItemID,
			OrderQty:           item.OrderQty,
			UnitPrice:          item.UnitPrice,
			TaxableItem:        item.TaxableItem,
			IncludeTax:         item.IncludeTax,
			TaxPercentage:      item.TaxPercentage,
			TaxAmount:          item.TaxAmount,
			UnitPriceTax:       item.UnitPriceTax,
			Subtotal:           item.Subtotal,
			Weight:             item.Weight,
			Note:               item.Note,
			PurchaseQty:        item.PurchaseQty,
		})
	}

	res = &bridgeService.GetPurchaseOrderDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.PurchaseOrder{
			Id:                     purchaseOrder.ID,
			Code:                   purchaseOrder.Code,
			VendorId:               purchaseOrder.VendorID,
			SiteId:                 purchaseOrder.SiteID,
			TermPaymentPurId:       purchaseOrder.TermPaymentPurID,
			VendorClassificationId: purchaseOrder.VendorClassificationID,
			PurchasePlanId:         purchaseOrder.PurchasePlanID,
			ConsolidatedShipmentId: purchaseOrder.ConsolidatedShipmentID,
			Status:                 purchaseOrder.Status,
			RecognitionDate:        timestamppb.New(purchaseOrder.RecognitionDate),
			EtaDate:                timestamppb.New(purchaseOrder.EtaDate),
			SiteAddress:            purchaseOrder.SiteAddress,
			EtaTime:                purchaseOrder.EtaTime,
			TaxPct:                 purchaseOrder.TaxPct,
			DeliveryFee:            purchaseOrder.DeliveryFee,
			TotalPrice:             purchaseOrder.TotalPrice,
			TaxAmount:              purchaseOrder.TaxAmount,
			TotalCharge:            purchaseOrder.TotalCharge,
			TotalInvoice:           purchaseOrder.TotalInvoice,
			TotalWeight:            purchaseOrder.TotalWeight,
			Note:                   purchaseOrder.Note,
			DeltaPrint:             purchaseOrder.DeltaPrint,
			Latitude:               purchaseOrder.Latitude,
			Longitude:              purchaseOrder.Longitude,
			UpdatedAt:              timestamppb.New(purchaseOrder.UpdatedAt),
			UpdatedBy:              purchaseOrder.UpdatedBy,
			CreatedFrom:            purchaseOrder.CreatedFrom,
			HasFinishedGr:          int32(purchaseOrder.HasFinishedGr),
			CreatedAt:              timestamppb.New(purchaseOrder.CreatedAt),
			CreatedBy:              purchaseOrder.CreatedBy,
			CommittedAt:            timestamppb.New(purchaseOrder.CommittedAt),
			CommittedBy:            purchaseOrder.CommittedBy,
			AssignedTo:             purchaseOrder.AssignedTo,
			AssignedBy:             purchaseOrder.AssignedBy,
			AssignedAt:             timestamppb.New(purchaseOrder.AssignedAt),
			Locked:                 purchaseOrder.Locked,
			LockedBy:               purchaseOrder.LockedBy,
			PurchaseOrderItems:     itemsRes,
		},
	}
	return
}

func (h *BridgeGrpcHandler) UpdatePurchaseOrder(ctx context.Context, req *bridgeService.UpdatePurchaseOrderRequest) (res *bridgeService.GetPurchaseOrderDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.UpdatePurchaseOrder")
	defer span.End()

	var (
		purchaseOrder dto.PurchaseOrderResponse
		items         []*dto.UpdatePurchaseOrderItemRequest
		itemsRes      []*bridgeService.PurchaseOrderItem
	)

	for _, item := range req.Items {
		items = append(items, &dto.UpdatePurchaseOrderItemRequest{
			ItemID:        item.ItemId,
			OrderQty:      item.OrderQty,
			UnitPrice:     item.UnitPrice,
			Note:          item.Note,
			PurchaseQty:   item.PurchaseQty,
			IncludeTax:    int8(item.IncludeTax),
			TaxPercentage: item.TaxPercentage,
		})
	}

	purchaseOrder, err = h.ServicesPurchaseOrder.Update(ctx, &dto.UpdatePurchaseOrderRequest{
		VendorID:           req.VendorId,
		SiteID:             req.SiteId,
		OrderDate:          req.OrderDate,
		StrEtaDate:         req.StrEtaDate,
		EtaTime:            req.EtaTime,
		DeliveryFee:        req.DeliveryFee,
		Note:               req.Note,
		TaxPct:             req.TaxPct,
		Latitude:           req.Latitude,
		Longitude:          req.Longitude,
		PurchaseOrderItems: items,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	for _, item := range purchaseOrder.PurchaseOrderItems {
		itemsRes = append(itemsRes, &bridgeService.PurchaseOrderItem{
			Id:                 item.ID,
			PurchaseOrderId:    item.PurchaseOrderID,
			PurchasePlanItemId: item.PurchasePlanItemID,
			ItemId:             item.ItemID,
			OrderQty:           item.OrderQty,
			UnitPrice:          item.UnitPrice,
			TaxableItem:        item.TaxableItem,
			IncludeTax:         item.IncludeTax,
			TaxPercentage:      item.TaxPercentage,
			TaxAmount:          item.TaxAmount,
			UnitPriceTax:       item.UnitPriceTax,
			Subtotal:           item.Subtotal,
			Weight:             item.Weight,
			Note:               item.Note,
			PurchaseQty:        item.PurchaseQty,
		})
	}

	res = &bridgeService.GetPurchaseOrderDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.PurchaseOrder{
			Id:                     purchaseOrder.ID,
			Code:                   purchaseOrder.Code,
			VendorId:               purchaseOrder.VendorID,
			SiteId:                 purchaseOrder.SiteID,
			TermPaymentPurId:       purchaseOrder.TermPaymentPurID,
			VendorClassificationId: purchaseOrder.VendorClassificationID,
			PurchasePlanId:         purchaseOrder.PurchasePlanID,
			ConsolidatedShipmentId: purchaseOrder.ConsolidatedShipmentID,
			Status:                 purchaseOrder.Status,
			RecognitionDate:        timestamppb.New(purchaseOrder.RecognitionDate),
			EtaDate:                timestamppb.New(purchaseOrder.EtaDate),
			SiteAddress:            purchaseOrder.SiteAddress,
			EtaTime:                purchaseOrder.EtaTime,
			TaxPct:                 purchaseOrder.TaxPct,
			DeliveryFee:            purchaseOrder.DeliveryFee,
			TotalPrice:             purchaseOrder.TotalPrice,
			TaxAmount:              purchaseOrder.TaxAmount,
			TotalCharge:            purchaseOrder.TotalCharge,
			TotalInvoice:           purchaseOrder.TotalInvoice,
			TotalWeight:            purchaseOrder.TotalWeight,
			Note:                   purchaseOrder.Note,
			DeltaPrint:             purchaseOrder.DeltaPrint,
			Latitude:               purchaseOrder.Latitude,
			Longitude:              purchaseOrder.Longitude,
			UpdatedAt:              timestamppb.New(purchaseOrder.UpdatedAt),
			UpdatedBy:              purchaseOrder.UpdatedBy,
			CreatedFrom:            purchaseOrder.CreatedFrom,
			HasFinishedGr:          int32(purchaseOrder.HasFinishedGr),
			CreatedAt:              timestamppb.New(purchaseOrder.CreatedAt),
			CreatedBy:              purchaseOrder.CreatedBy,
			CommittedAt:            timestamppb.New(purchaseOrder.CommittedAt),
			CommittedBy:            purchaseOrder.CommittedBy,
			AssignedTo:             purchaseOrder.AssignedTo,
			AssignedBy:             purchaseOrder.AssignedBy,
			AssignedAt:             timestamppb.New(purchaseOrder.AssignedAt),
			Locked:                 purchaseOrder.Locked,
			LockedBy:               purchaseOrder.LockedBy,
			PurchaseOrderItems:     itemsRes,
		},
	}
	return
}

func (h *BridgeGrpcHandler) UpdateProductPurchaseOrder(ctx context.Context, req *bridgeService.UpdateProductPurchaseOrderRequest) (res *bridgeService.GetPurchaseOrderDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.UpdatePurchaseOrder")
	defer span.End()

	var (
		purchaseOrder dto.PurchaseOrderResponse
		items         []*dto.UpdatePurchaseOrderItemRequest
		itemsRes      []*bridgeService.PurchaseOrderItem
	)

	for _, item := range req.Items {
		items = append(items, &dto.UpdatePurchaseOrderItemRequest{
			ItemID:        item.ItemId,
			OrderQty:      item.OrderQty,
			UnitPrice:     item.UnitPrice,
			Note:          item.Note,
			PurchaseQty:   item.PurchaseQty,
			IncludeTax:    int8(item.IncludeTax),
			TaxPercentage: item.TaxPercentage,
		})
	}

	purchaseOrder, err = h.ServicesPurchaseOrder.UpdateProduct(ctx, &dto.UpdateProductPurchaseOrderRequest{
		DeliveryFee:        req.DeliveryFee,
		TaxPct:             req.TaxPct,
		PurchaseOrderItems: items,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	for _, item := range purchaseOrder.PurchaseOrderItems {
		itemsRes = append(itemsRes, &bridgeService.PurchaseOrderItem{
			Id:                 item.ID,
			PurchaseOrderId:    item.PurchaseOrderID,
			PurchasePlanItemId: item.PurchasePlanItemID,
			ItemId:             item.ItemID,
			OrderQty:           item.OrderQty,
			UnitPrice:          item.UnitPrice,
			TaxableItem:        item.TaxableItem,
			IncludeTax:         item.IncludeTax,
			TaxPercentage:      item.TaxPercentage,
			TaxAmount:          item.TaxAmount,
			UnitPriceTax:       item.UnitPriceTax,
			Subtotal:           item.Subtotal,
			Weight:             item.Weight,
			Note:               item.Note,
			PurchaseQty:        item.PurchaseQty,
		})
	}

	res = &bridgeService.GetPurchaseOrderDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.PurchaseOrder{
			Id:                     purchaseOrder.ID,
			Code:                   purchaseOrder.Code,
			VendorId:               purchaseOrder.VendorID,
			SiteId:                 purchaseOrder.SiteID,
			TermPaymentPurId:       purchaseOrder.TermPaymentPurID,
			VendorClassificationId: purchaseOrder.VendorClassificationID,
			PurchasePlanId:         purchaseOrder.PurchasePlanID,
			ConsolidatedShipmentId: purchaseOrder.ConsolidatedShipmentID,
			Status:                 purchaseOrder.Status,
			RecognitionDate:        timestamppb.New(purchaseOrder.RecognitionDate),
			EtaDate:                timestamppb.New(purchaseOrder.EtaDate),
			SiteAddress:            purchaseOrder.SiteAddress,
			EtaTime:                purchaseOrder.EtaTime,
			TaxPct:                 purchaseOrder.TaxPct,
			DeliveryFee:            purchaseOrder.DeliveryFee,
			TotalPrice:             purchaseOrder.TotalPrice,
			TaxAmount:              purchaseOrder.TaxAmount,
			TotalCharge:            purchaseOrder.TotalCharge,
			TotalInvoice:           purchaseOrder.TotalInvoice,
			TotalWeight:            purchaseOrder.TotalWeight,
			Note:                   purchaseOrder.Note,
			DeltaPrint:             purchaseOrder.DeltaPrint,
			Latitude:               purchaseOrder.Latitude,
			Longitude:              purchaseOrder.Longitude,
			UpdatedAt:              timestamppb.New(purchaseOrder.UpdatedAt),
			UpdatedBy:              purchaseOrder.UpdatedBy,
			CreatedFrom:            purchaseOrder.CreatedFrom,
			HasFinishedGr:          int32(purchaseOrder.HasFinishedGr),
			CreatedAt:              timestamppb.New(purchaseOrder.CreatedAt),
			CreatedBy:              purchaseOrder.CreatedBy,
			CommittedAt:            timestamppb.New(purchaseOrder.CommittedAt),
			CommittedBy:            purchaseOrder.CommittedBy,
			AssignedTo:             purchaseOrder.AssignedTo,
			AssignedBy:             purchaseOrder.AssignedBy,
			AssignedAt:             timestamppb.New(purchaseOrder.AssignedAt),
			Locked:                 purchaseOrder.Locked,
			LockedBy:               purchaseOrder.LockedBy,
			PurchaseOrderItems:     itemsRes,
		},
	}
	return
}

func (h *BridgeGrpcHandler) CommitPurchaseOrder(ctx context.Context, req *bridgeService.CommitPurchaseOrderRequest) (res *bridgeService.GetPurchaseOrderDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPurchaseOrderDetail")
	defer span.End()

	err = h.ServicesPurchaseOrder.Commit(ctx, req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetPurchaseOrderDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    nil,
	}
	return
}

func (h *BridgeGrpcHandler) CancelPurchaseOrder(ctx context.Context, req *bridgeService.CancelPurchaseOrderRequest) (res *bridgeService.GetPurchaseOrderDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPurchaseOrderDetail")
	defer span.End()

	err = h.ServicesPurchaseOrder.Cancel(ctx, req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetPurchaseOrderDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    nil,
	}
	return
}

func (h *BridgeGrpcHandler) GetPurchaseOrderGPList(ctx context.Context, req *bridgeService.GetPurchaseOrderGPListRequest) (res *bridgeService.GetPurchaseOrderGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPurchaseOrderGPList")
	defer span.End()

	res, err = h.ServicesPurchaseOrder.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) GetPurchaseOrderGPDetail(ctx context.Context, req *bridgeService.GetPurchaseOrderGPDetailRequest) (res *bridgeService.GetPurchaseOrderGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPurchaseOrderGPDetail")
	defer span.End()

	res, err = h.ServicesPurchaseOrder.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) CreatePurchaseOrderGP(ctx context.Context, req *bridgeService.CreatePurchaseOrderGPRequest) (res *bridgeService.CreatePurchaseOrderGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreatePurchaseOrderGP")
	defer span.End()

	var poDetail []dto.PODTL
	for _, podtl := range req.Detail {
		poDetail = append(poDetail, dto.PODTL{
			Ord:      podtl.Ord,
			Itemnmbr: podtl.Itemnmbr,
			Uofm:     podtl.Uofm,
			Qtyorder: podtl.Qtyorder,
			Qtycance: podtl.Qtycance,
			Unitcost: podtl.Unitcost,
			Notetext: podtl.Notetext,
		})
	}

	var resAPi dto.CommonPurchaseOrderGPResponse
	resAPi, err = h.ServicesPurchaseOrder.CreateGP(ctx, &dto.CreatePurchaseOrderGPRequest{
		Interid:           req.Interid,
		Potype:            int64(req.Potype),
		Ponumber:          req.Ponumber,
		Docdate:           req.Docdate,
		Buyerid:           req.Buyerid,
		Vendorid:          req.Vendorid,
		Curncyid:          req.Curncyid,
		Deprtmnt:          req.Deprtmnt,
		Locncode:          req.Locncode,
		Taxschid:          req.Taxschid,
		Subtotal:          req.Subtotal,
		Trdisamt:          req.Trdisamt,
		Frtamnt:           req.Frtamnt,
		Miscamnt:          req.Miscamnt,
		Taxamnt:           req.Taxamnt,
		PrpPurchaseplanNo: req.PrpPurchaseplanNo,
		// CsReference:             *req.CsReference,
		PrpPaymentMethod:        req.PrpPaymentMethod,
		PrpRegion:               req.PrpRegion,
		PrpEstimatedarrivalDate: req.PrpEstimatedarrivalDate,
		Notetext:                req.Notetext,
		Detail:                  poDetail,
		PrpPaymentTerm:          req.Pymtrmid,
		DueDate:                 req.Duedate,
		PRStatus:                req.PrStatus,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.CreatePurchaseOrderGPResponse{
		Code:     int32(codes.OK),
		Message:  codes.OK.String(),
		Ponumber: resAPi.Ponumber,
	}
	return
}

func (h *BridgeGrpcHandler) UpdatePurchaseOrderGP(ctx context.Context, req *bridgeService.UpdatePurchaseOrderGPRequest) (res *bridgeService.CreatePurchaseOrderGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.UpdatePurchaseOrderGP")
	defer span.End()

	var poDetail []dto.PODTL
	for _, podtl := range req.Detail {
		poDetail = append(poDetail, dto.PODTL{
			Ord:      podtl.Ord,
			Itemnmbr: podtl.Itemnmbr,
			Uofm:     podtl.Uofm,
			Qtyorder: podtl.Qtyorder,
			Qtycance: podtl.Qtycance,
			Unitcost: podtl.Unitcost,
			Notetext: podtl.Notetext,
		})
	}

	var resAPi dto.CommonPurchaseOrderGPResponse
	resAPi, err = h.ServicesPurchaseOrder.UpdateGP(ctx, &dto.UpdatePurchaseOrderGPRequest{
		Interid:           req.Interid,
		Potype:            int64(req.Potype),
		Ponumber:          req.Ponumber,
		Docdate:           req.Docdate,
		Buyerid:           req.Buyerid,
		Vendorid:          req.Vendorid,
		Curncyid:          req.Curncyid,
		Deprtmnt:          req.Deprtmnt,
		Locncode:          req.Locncode,
		Taxschid:          req.Taxschid,
		Subtotal:          req.Subtotal,
		Trdisamt:          req.Trdisamt,
		Frtamnt:           req.Frtamnt,
		Miscamnt:          req.Miscamnt,
		Taxamnt:           req.Taxamnt,
		PrpPurchaseplanNo: req.PrpPurchaseplanNo,
		CsReference: &dto.ConsolidatedShipmentReference{
			PRPCSNo:      req.CsReference.PrpCsNo,
			PRVehicleNo:  req.CsReference.PrpVehicleNo,
			PRDriverName: req.CsReference.PrpDriverName,
			PhonName:     req.CsReference.Phonname,
		},
		Pymtrmid:                req.Pymtrmid,
		Duedate:                 req.Duedate,
		PrpPaymentMethod:        req.PrpPaymentMethod,
		PrpRegion:               req.PrpRegion,
		PrpEstimatedarrivalDate: req.PrpEstimatedarrivalDate,
		Notetext:                req.Notetext,
		Detail:                  poDetail,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.CreatePurchaseOrderGPResponse{
		Code:     int32(codes.OK),
		Message:  codes.OK.String(),
		Ponumber: resAPi.Ponumber,
	}
	return
}

func (h *BridgeGrpcHandler) CommitPurchaseOrderGP(ctx context.Context, req *bridgeService.CommitPurchaseOrderGPRequest) (res *bridgeService.CreateTransferRequestGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateTransferRequestGP")
	defer span.End()

	res, err = h.ServicePurchaseOrder.CommitPurchaseOrderGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) CancelPurchaseOrderGP(ctx context.Context, req *bridgeService.CancelPurchaseOrderGPRequest) (res *bridgeService.CancelPurchaseOrderGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CancelPurchaseOrder")
	defer span.End()

	_, err = h.ServicePurchaseOrder.CancelPurchaseOrderGP(ctx, &dto.CancelPurchaseOrderGPRequest{
		PurchaseOrderNumber: req.PoNumber,
		UserId:              req.UserId,
	})

	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	h.Option.Common.Logger.AddMessage(log.InfoLevel, "CancelPurchaseOrder").Print()

	res = &bridgeService.CancelPurchaseOrderGPResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}

	return
}

func (h *BridgeGrpcHandler) CreateConsolidatedShipmentGP(ctx context.Context, req *bridgeService.CreateConsolidatedShipmentGPRequest) (res *bridgeService.CreateConsolidatedShipmentGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateConsolidatedShipmentGP")
	defer span.End()

	var poList []dto.POLST
	for _, polst := range req.PurchaseOrder {
		poList = append(poList, dto.POLST{
			Ponumber: polst.Ponumber,
		})
	}

	_, err = h.ServicesPurchaseOrder.CreateConsolidatedShipmentGP(ctx, &dto.CreateConsolidatedShipmentGPRequest{
		Interid:         req.Interid,
		PRPCSNo:         req.PrpCsNo,
		PrpDriverName:   req.PrpDriverName,
		PrVehicleNumber: req.PrVehicleNumber,
		PhoneName:       req.Phonname,
		PurchaseOrders:  poList,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	h.Option.Common.Logger.AddMessage(log.InfoLevel, "CreateConsolidatedShipmentGP").Print()

	res = &bridgeService.CreateConsolidatedShipmentGPResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}

func (h *BridgeGrpcHandler) UpdateConsolidatedShipmentGP(ctx context.Context, req *bridgeService.UpdateConsolidatedShipmentGPRequest) (res *bridgeService.UpdateConsolidatedShipmentGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.UpdateConsolidatedShipmentGP")
	defer span.End()

	var poList []dto.POLST
	for _, polst := range req.PurchaseOrder {
		poList = append(poList, dto.POLST{
			Ponumber: polst.Ponumber,
		})
	}

	_, err = h.ServicesPurchaseOrder.UpdateConsolidatedShipmentGP(ctx, &dto.UpdateConsolidatedShipmentGPRequest{
		Interid:         req.Interid,
		PRPCSNo:         req.PrpCsNo,
		PrpDriverName:   req.PrpDriverName,
		PrVehicleNumber: req.PrVehicleNumber,
		PhoneName:       req.Phonname,
		PurchaseOrders:  poList,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	h.Option.Common.Logger.AddMessage(log.InfoLevel, "UpdateConsolidatedShipmentGP").Print()

	res = &bridgeService.UpdateConsolidatedShipmentGPResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}
