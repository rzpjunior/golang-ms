package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetSalesInvoiceItemList(ctx context.Context, req *bridgeService.GetSalesInvoiceItemListRequest) (res *bridgeService.GetSalesInvoiceItemListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesInvoiceItemList")
	defer span.End()

	var salesInvoiceItems []dto.SalesInvoiceItemResponse
	salesInvoiceItems, _, err = h.ServiceSalesInvoiceItem.Get(ctx, req.SalesInvoiceId)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var salesInvoiceItemsGrpc []*bridgeService.SalesInvoiceItem
	for _, salesInvoiceItem := range salesInvoiceItems {
		salesInvoiceItemsGrpc = append(salesInvoiceItemsGrpc, &bridgeService.SalesInvoiceItem{
			Id:               salesInvoiceItem.ID,
			SalesOrderItemId: salesInvoiceItem.SalesOrderItemID,
			SalesInvoiceId:   salesInvoiceItem.SalesInvoiceID,
			ItemId:           salesInvoiceItem.ItemID,
			Note:             salesInvoiceItem.Note,
			InvoiceQty:       salesInvoiceItem.InvoiceQty,
			UnitPrice:        salesInvoiceItem.UnitPrice,
			Subtotal:         salesInvoiceItem.Subtotal,
			TaxableItem:      int32(salesInvoiceItem.TaxableItem),
			TaxPercentage:    salesInvoiceItem.TaxPercentage,
			SkuDiscAmount:    salesInvoiceItem.SkuDiscAmount,
		})
	}

	res = &bridgeService.GetSalesInvoiceItemListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    salesInvoiceItemsGrpc,
	}
	return
}
