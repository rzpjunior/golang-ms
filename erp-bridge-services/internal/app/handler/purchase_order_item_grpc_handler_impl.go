package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetPurchaseOrderItemList(ctx context.Context, req *bridgeService.GetPurchaseOrderItemListRequest) (res *bridgeService.GetPurchaseOrderItemListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPurchaseOrderItemList")
	defer span.End()

	var purchaseOrderItems []dto.PurchaseOrderItemResponse
	purchaseOrderItems, _, err = h.ServicesPurchaseOrderItem.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.PurchaseOrderItem
	for _, purchaseOrderItem := range purchaseOrderItems {
		data = append(data, &bridgeService.PurchaseOrderItem{
			Id:                 purchaseOrderItem.ID,
			PurchaseOrderId:    purchaseOrderItem.PurchaseOrderID,
			PurchasePlanItemId: purchaseOrderItem.PurchasePlanItemID,
			ItemId:             purchaseOrderItem.ItemID,
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

	res = &bridgeService.GetPurchaseOrderItemListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetPurchaseOrderItemDetail(ctx context.Context, req *bridgeService.GetPurchaseOrderItemDetailRequest) (res *bridgeService.GetPurchaseOrderItemDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPurchaseOrderItemDetail")
	defer span.End()

	var purchaseOrderItem dto.PurchaseOrderItemResponse
	purchaseOrderItem, err = h.ServicesPurchaseOrderItem.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetPurchaseOrderItemDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.PurchaseOrderItem{
			Id:                 purchaseOrderItem.ID,
			PurchaseOrderId:    purchaseOrderItem.PurchaseOrderID,
			PurchasePlanItemId: purchaseOrderItem.PurchasePlanItemID,
			ItemId:             purchaseOrderItem.ItemID,
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
		},
	}
	return
}
