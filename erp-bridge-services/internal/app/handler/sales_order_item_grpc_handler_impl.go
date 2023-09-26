package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetSalesOrderItemList(ctx context.Context, req *bridgeService.GetSalesOrderItemListRequest) (res *bridgeService.GetSalesOrderItemListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesOrderItemList")
	defer span.End()

	var salesOrderItems []dto.SalesOrderItemResponse
	salesOrderItems, _, err = h.ServicesSalesOrderItem.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy, req.SalesOrderId, req.SalesOrderId)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.SalesOrderItem
	for _, salesOrderItem := range salesOrderItems {
		data = append(data, &bridgeService.SalesOrderItem{
			Id:            salesOrderItem.ID,
			SalesOrderId:  salesOrderItem.SalesOrderId,
			ItemId:        "1-123ITMS",
			OrderQty:      salesOrderItem.OrderQty,
			DefaultPrice:  salesOrderItem.DefaultPrice,
			UnitPrice:     salesOrderItem.UnitPrice,
			TaxableItem:   salesOrderItem.TaxableItem,
			TaxPercentage: salesOrderItem.TaxPercentage,
			ShadowPrice:   salesOrderItem.ShadowPrice,
			Subtotal:      salesOrderItem.Subtotal,
			Weight:        salesOrderItem.Weight,
			Note:          salesOrderItem.Note,
		})
	}

	res = &bridgeService.GetSalesOrderItemListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetSalesOrderItemDetail(ctx context.Context, req *bridgeService.GetSalesOrderItemDetailRequest) (res *bridgeService.GetSalesOrderItemDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesOrderItemDetail")
	defer span.End()

	var salesOrderItem dto.SalesOrderItemResponse
	salesOrderItem, err = h.ServicesSalesOrderItem.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetSalesOrderItemDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.SalesOrderItem{
			Id:            salesOrderItem.ID,
			SalesOrderId:  salesOrderItem.SalesOrderId,
			ItemId:        "1-123ITMS",
			OrderQty:      salesOrderItem.OrderQty,
			DefaultPrice:  salesOrderItem.DefaultPrice,
			UnitPrice:     salesOrderItem.UnitPrice,
			TaxableItem:   salesOrderItem.TaxableItem,
			TaxPercentage: salesOrderItem.TaxPercentage,
			ShadowPrice:   salesOrderItem.ShadowPrice,
			Subtotal:      salesOrderItem.Subtotal,
			Weight:        salesOrderItem.Weight,
			Note:          salesOrderItem.Note,
		},
	}
	return
}
