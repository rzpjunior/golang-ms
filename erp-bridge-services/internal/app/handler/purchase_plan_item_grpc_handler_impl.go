package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetPurchasePlanItemList(ctx context.Context, req *bridgeService.GetPurchasePlanItemListRequest) (res *bridgeService.GetPurchasePlanItemListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPurchasePlanItemList")
	defer span.End()

	var purchasePlanItems []dto.PurchasePlanItemResponse
	purchasePlanItems, _, err = h.ServicesPurchasePlanItem.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.PurchasePlanItem
	for _, purchasePlanItem := range purchasePlanItems {
		data = append(data, &bridgeService.PurchasePlanItem{
			Id:              purchasePlanItem.ID,
			PurchasePlanId:  purchasePlanItem.PurchasePlanID,
			ItemId:          purchasePlanItem.ItemID,
			PurchasePlanQty: purchasePlanItem.PurchasePlanQty,
			PurchaseQty:     purchasePlanItem.PurchaseQty,
			UnitPrice:       purchasePlanItem.UnitPrice,
			Subtotal:        purchasePlanItem.Subtotal,
			Weight:          purchasePlanItem.Weight,
		})
	}

	res = &bridgeService.GetPurchasePlanItemListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetPurchasePlanItemDetail(ctx context.Context, req *bridgeService.GetPurchasePlanItemDetailRequest) (res *bridgeService.GetPurchasePlanItemDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPurchasePlanItemDetail")
	defer span.End()

	var purchasePlanItem dto.PurchasePlanItemResponse
	purchasePlanItem, err = h.ServicesPurchasePlanItem.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetPurchasePlanItemDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.PurchasePlanItem{
			Id:              purchasePlanItem.ID,
			PurchasePlanId:  purchasePlanItem.PurchasePlanID,
			ItemId:          purchasePlanItem.ItemID,
			PurchasePlanQty: purchasePlanItem.PurchasePlanQty,
			PurchaseQty:     purchasePlanItem.PurchaseQty,
			UnitPrice:       purchasePlanItem.UnitPrice,
			Subtotal:        purchasePlanItem.Subtotal,
			Weight:          purchasePlanItem.Weight,
		},
	}
	return
}
