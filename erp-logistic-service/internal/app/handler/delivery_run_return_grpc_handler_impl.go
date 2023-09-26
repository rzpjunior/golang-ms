package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/dto"
	logisticService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *LogisticGrpcHandler) GetDeliveryRunReturnList(ctx context.Context, req *logisticService.GetDeliveryRunReturnListRequest) (res *logisticService.GetDeliveryRunReturnListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetDeliveryRunReturnList")
	defer span.End()

	deliveryRunReturns, _, err := h.ServicesDeliveryRunReturn.Get(ctx, dto.DeliveryRunReturnGetRequest{
		Offset:                     int(req.Offset),
		Limit:                      int(req.Limit),
		OrderBy:                    req.OrderBy,
		ArrDeliveryRunSheetItemIDs: req.DeliveryRunSheetItemId,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*logisticService.DeliveryRunReturn
	for _, drr := range deliveryRunReturns {
		data = append(data, &logisticService.DeliveryRunReturn{
			Id:                     drr.ID,
			Code:                   drr.Code,
			DeliveryRunSheetItemId: drr.DeliveryRunSheetItemID,
			TotalPrice:             drr.TotalPrice,
			TotalCharge:            drr.TotalCharge,
			CreatedAt:              timestamppb.New(drr.CreatedAt),
		})
	}

	res = &logisticService.GetDeliveryRunReturnListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}

	return
}

func (h *LogisticGrpcHandler) GetDeliveryRunReturnDetail(ctx context.Context, req *logisticService.GetDeliveryRunReturnDetailRequest) (res *logisticService.GetDeliveryRunReturnDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetDeliveryRunReturnDetail")
	defer span.End()

	deliveryRunReturn, err := h.ServicesDeliveryRunReturn.GetDetail(ctx, req.Id, req.Code, req.DeliveryRunSheetItemId)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.GetDeliveryRunReturnDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunReturn{
			Id:                     deliveryRunReturn.ID,
			Code:                   deliveryRunReturn.Code,
			DeliveryRunSheetItemId: deliveryRunReturn.DeliveryRunSheetItemID,
			TotalPrice:             deliveryRunReturn.TotalPrice,
			TotalCharge:            deliveryRunReturn.TotalCharge,
			CreatedAt:              timestamppb.New(deliveryRunReturn.CreatedAt),
		},
	}

	return
}

func (h *LogisticGrpcHandler) CreateDeliveryRunReturn(ctx context.Context, req *logisticService.CreateDeliveryRunReturnRequest) (res *logisticService.CreateDeliveryRunReturnResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateDeliveryRunReturn")
	defer span.End()

	deliveryRunReturn, err := h.ServicesDeliveryRunReturn.Create(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.CreateDeliveryRunReturnResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunReturn{
			Id:                     deliveryRunReturn.ID,
			Code:                   deliveryRunReturn.Code,
			DeliveryRunSheetItemId: deliveryRunReturn.DeliveryRunSheetItemID,
			TotalPrice:             deliveryRunReturn.TotalPrice,
			TotalCharge:            deliveryRunReturn.TotalCharge,
			CreatedAt:              timestamppb.New(deliveryRunReturn.CreatedAt),
		},
	}

	return
}

func (h *LogisticGrpcHandler) UpdateDeliveryRunReturn(ctx context.Context, req *logisticService.UpdateDeliveryRunReturnRequest) (res *logisticService.UpdateDeliveryRunReturnResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.UpdateDeliveryRunReturn")
	defer span.End()

	deliveryRunReturn, err := h.ServicesDeliveryRunReturn.Update(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.UpdateDeliveryRunReturnResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunReturn{
			Id:                     deliveryRunReturn.ID,
			Code:                   deliveryRunReturn.Code,
			DeliveryRunSheetItemId: deliveryRunReturn.DeliveryRunSheetItemID,
			TotalPrice:             deliveryRunReturn.TotalPrice,
			TotalCharge:            deliveryRunReturn.TotalCharge,
			CreatedAt:              timestamppb.New(deliveryRunReturn.CreatedAt),
		},
	}

	return
}

func (h *LogisticGrpcHandler) DeleteDeliveryRunReturn(ctx context.Context, req *logisticService.DeleteDeliveryRunReturnRequest) (res *logisticService.DeleteDeliveryRunReturnResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.DeleteDeliveryRunReturn")
	defer span.End()

	deliveryRunReturn, err := h.ServicesDeliveryRunReturn.Delete(ctx, req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.DeleteDeliveryRunReturnResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.DeliveryRunReturn{
			Id:                     deliveryRunReturn.ID,
			Code:                   deliveryRunReturn.Code,
			DeliveryRunSheetItemId: deliveryRunReturn.DeliveryRunSheetItemID,
			TotalPrice:             deliveryRunReturn.TotalPrice,
			TotalCharge:            deliveryRunReturn.TotalCharge,
			CreatedAt:              timestamppb.New(deliveryRunReturn.CreatedAt),
		},
	}

	return
}
