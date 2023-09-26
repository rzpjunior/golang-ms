package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetItemTransferItemDetail(ctx context.Context, req *bridgeService.GetItemTransferItemDetailRequest) (res *bridgeService.GetItemTransferItemDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemTransferItemDetail")
	defer span.End()

	var (
		itemTransferItem dto.ItemTransferItemResponse
	)
	itemTransferItem, err = h.ServicesItemTransferItem.GetDetail(ctx, req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetItemTransferItemDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.ItemTransferItem{
			Id:          itemTransferItem.ID,
			DeliverQty:  itemTransferItem.DeliverQty,
			ReceiveQty:  itemTransferItem.ReceiveQty,
			RequestQty:  itemTransferItem.RequestQty,
			ReceiveNote: itemTransferItem.ReceiveNote,
			UnitCost:    itemTransferItem.UnitCost,
			Subtotal:    itemTransferItem.Subtotal,
			Weight:      itemTransferItem.Weight,
			Note:        itemTransferItem.Note,
		},
	}
	return
}
