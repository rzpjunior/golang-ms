package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	logisticService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *LogisticGrpcHandler) CreateMerchantDeliveryLog(ctx context.Context, req *logisticService.CreateMerchantDeliveryLogRequest) (res *logisticService.CreateMerchantDeliveryLogResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateMerchantDeliveryLog")
	defer span.End()

	merchantDeliveryLog, err := h.ServicesMerchantDeliveryLog.Create(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.CreateMerchantDeliveryLogResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.MerchantDeliveryLog{
			Id:                     merchantDeliveryLog.Id,
			DeliveryRunSheetItemId: merchantDeliveryLog.DeliveryRunSheetItemId,
			Latitude:               merchantDeliveryLog.Latitude,
			Longitude:              merchantDeliveryLog.Longitude,
			CreatedAt:              timestamppb.New(merchantDeliveryLog.CreatedAt),
		},
	}

	return
}

func (h *LogisticGrpcHandler) GetFirstMerchantDeliveryLog(ctx context.Context, req *logisticService.GetFirstMerchantDeliveryLogRequest) (res *logisticService.GetFirstMerchantDeliveryLogResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetFirstMerchantDeliveryLog")
	defer span.End()

	merchantDeliveryLog, err := h.ServicesMerchantDeliveryLog.GetFirst(ctx, req.DeliveryRunSheetItemId)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.GetFirstMerchantDeliveryLogResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.MerchantDeliveryLog{
			Id:                     merchantDeliveryLog.Id,
			DeliveryRunSheetItemId: merchantDeliveryLog.DeliveryRunSheetItemId,
			Latitude:               merchantDeliveryLog.Latitude,
			Longitude:              merchantDeliveryLog.Longitude,
			CreatedAt:              timestamppb.New(merchantDeliveryLog.CreatedAt),
		},
	}

	return
}
