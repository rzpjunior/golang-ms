package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/dto"
	logisticService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *LogisticGrpcHandler) GetDeliveryRunReturnItemList(ctx context.Context, req *logisticService.GetDeliveryRunReturnItemListRequest) (res *logisticService.GetDeliveryRunReturnItemListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetDeliveryRunReturnItemList")
	defer span.End()

	deliveryRunReturnItems, _, err := h.ServicesDeliveryRunReturnItem.Get(ctx, dto.DeliveryRunReturnItemGetRequest{
		Offset:                  int(req.Offset),
		Limit:                   int(req.Limit),
		OrderBy:                 req.OrderBy,
		ArrDeliveryRunReturnIDs: req.DeliveryRunReturnId,
		ArrDeliveryOrderItemIDs: req.DeliveryOrderItemId,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*logisticService.DeliveryRunReturnItem
	for _, drri := range deliveryRunReturnItems {
		data = append(data, &logisticService.DeliveryRunReturnItem{
			Id:                  drri.ID,
			DeliveryRunReturnId: drri.DeliveryRunReturnID,
			DeliveryOrderItemId: drri.DeliveryOrderItemID,
			ReceiveQty:          drri.ReceiveQty,
			ReturnReason:        int32(drri.ReturnReason),
			ReturnEvidence:      drri.ReturnEvidence,
			Subtotal:            drri.Subtotal,
		})
	}

	res = &logisticService.GetDeliveryRunReturnItemListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}

	return
}

func (h *LogisticGrpcHandler) GetDeliveryRunReturnItemDetail(ctx context.Context, req *logisticService.GetDeliveryRunReturnItemDetailRequest) (res *logisticService.GetDeliveryRunReturnItemDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetDeliveryRunReturnItemDetail")
	defer span.End()

	deliveryRunReturnItem, err := h.ServicesDeliveryRunReturnItem.GetDetail(ctx, req.Id, req.DeliveryOrderItemId)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.GetDeliveryRunReturnItemDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunReturnItem{
			Id:                  deliveryRunReturnItem.ID,
			DeliveryRunReturnId: deliveryRunReturnItem.DeliveryRunReturnID,
			DeliveryOrderItemId: deliveryRunReturnItem.DeliveryOrderItemID,
			ReceiveQty:          deliveryRunReturnItem.ReceiveQty,
			ReturnReason:        int32(deliveryRunReturnItem.ReturnReason),
			ReturnEvidence:      deliveryRunReturnItem.ReturnEvidence,
			Subtotal:            deliveryRunReturnItem.Subtotal,
		},
	}

	return
}

func (h *LogisticGrpcHandler) CreateDeliveryRunReturnItem(ctx context.Context, req *logisticService.CreateDeliveryRunReturnItemRequest) (res *logisticService.CreateDeliveryRunReturnItemResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateDeliveryRunReturnItem")
	defer span.End()

	deliveryRunReturnItem, err := h.ServicesDeliveryRunReturnItem.Create(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.CreateDeliveryRunReturnItemResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunReturnItem{
			Id:                  deliveryRunReturnItem.ID,
			DeliveryRunReturnId: deliveryRunReturnItem.DeliveryRunReturnID,
			DeliveryOrderItemId: deliveryRunReturnItem.DeliveryOrderItemID,
			ReceiveQty:          deliveryRunReturnItem.ReceiveQty,
			ReturnReason:        int32(deliveryRunReturnItem.ReturnReason),
			ReturnEvidence:      deliveryRunReturnItem.ReturnEvidence,
			Subtotal:            deliveryRunReturnItem.Subtotal,
		},
	}

	return
}

func (h *LogisticGrpcHandler) UpdateDeliveryRunReturnItem(ctx context.Context, req *logisticService.UpdateDeliveryRunReturnItemRequest) (res *logisticService.UpdateDeliveryRunReturnItemResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.UpdateDeliveryRunReturnItem")
	defer span.End()

	deliveryRunReturnItem, err := h.ServicesDeliveryRunReturnItem.Update(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.UpdateDeliveryRunReturnItemResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunReturnItem{
			Id:                  deliveryRunReturnItem.ID,
			DeliveryRunReturnId: deliveryRunReturnItem.DeliveryRunReturnID,
			DeliveryOrderItemId: deliveryRunReturnItem.DeliveryOrderItemID,
			ReceiveQty:          deliveryRunReturnItem.ReceiveQty,
			ReturnReason:        int32(deliveryRunReturnItem.ReturnReason),
			ReturnEvidence:      deliveryRunReturnItem.ReturnEvidence,
			Subtotal:            deliveryRunReturnItem.Subtotal,
		},
	}

	return
}

func (h *LogisticGrpcHandler) DeleteDeliveryRunReturnItem(ctx context.Context, req *logisticService.DeleteDeliveryRunReturnItemRequest) (res *logisticService.DeleteDeliveryRunReturnItemResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.DeleteDeliveryRunReturnItem")
	defer span.End()

	deliveryRunReturnItem, err := h.ServicesDeliveryRunReturnItem.Delete(ctx, req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.DeleteDeliveryRunReturnItemResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunReturnItem{
			Id:                  deliveryRunReturnItem.ID,
			DeliveryRunReturnId: deliveryRunReturnItem.DeliveryRunReturnID,
			DeliveryOrderItemId: deliveryRunReturnItem.DeliveryOrderItemID,
			ReceiveQty:          deliveryRunReturnItem.ReceiveQty,
			ReturnReason:        int32(deliveryRunReturnItem.ReturnReason),
			ReturnEvidence:      deliveryRunReturnItem.ReturnEvidence,
			Subtotal:            deliveryRunReturnItem.Subtotal,
		},
	}

	return
}
