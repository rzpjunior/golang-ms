package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	logisticService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *LogisticGrpcHandler) CreatePostponeDeliveryLog(ctx context.Context, req *logisticService.CreatePostponeDeliveryLogRequest) (res *logisticService.CreatePostponeDeliveryLogResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreatePostponeDeliveryLog")
	defer span.End()

	postponeDeliveryLog, err := h.ServicesPostponeDeliveryLog.Create(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.CreatePostponeDeliveryLogResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &logisticService.PostponeDeliveryLog{
			Id:                     postponeDeliveryLog.ID,
			DeliveryRunSheetItemId: postponeDeliveryLog.DeliveryRunSheetItemID,
			PostponeReason:         postponeDeliveryLog.PostponeReason,
			StartedAtUnix:          postponeDeliveryLog.StartedAt.Unix(),
			PostponedAtUnix:        postponeDeliveryLog.PostponedAt.Unix(),
			PostponeEvidence:       postponeDeliveryLog.PostponeEvidence,
		},
	}

	return
}
