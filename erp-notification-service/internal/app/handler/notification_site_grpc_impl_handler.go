package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/notification_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *NotificationGrpcHandler) SendNotificationHelper(ctx context.Context, req *pb.SendNotificationHelperRequest) (res *pb.SuccessResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "NotificationGrpcHandler.SendNotificationHelper")
	defer span.End()

	res = &pb.SuccessResponse{}
	res, err = h.ServicesNotificationSite.SendNotification(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}
