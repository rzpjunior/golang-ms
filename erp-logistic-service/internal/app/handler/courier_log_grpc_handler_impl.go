package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	logisticService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *LogisticGrpcHandler) CreateCourierLog(ctx context.Context, req *logisticService.CreateCourierLogRequest) (res *logisticService.CreateCourierLogResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateCourierLog")
	defer span.End()

	courierLog, err := h.ServicesCourierLog.Create(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.CreateCourierLogResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.CourierLog{
			Id:           courierLog.ID,
			CourierId:    courierLog.CourierID,
			SalesOrderId: courierLog.SalesOrderID,
			Latitude:     courierLog.Latitude,
			Longitude:    courierLog.Longitude,
			CreatedAt:    timestamppb.New(courierLog.CreatedAt),
		},
	}

	return
}

func (h *LogisticGrpcHandler) GetLastCourierLog(ctx context.Context, req *logisticService.GetLastCourierLogRequest) (res *logisticService.GetLastCourierLogResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetLastCourierLog")
	defer span.End()

	res, err = h.ServicesCourierLog.GetLastCourierLog(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	res.Code = int32(codes.OK)
	res.Message = codes.OK.String()

	return
}
