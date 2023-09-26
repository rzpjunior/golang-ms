package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/notification_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *NotificationGrpcHandler) SendNotificationCancelSalesOrder(ctx context.Context, req *pb.SendNotificationCancelSalesOrderRequest) (res *pb.SendNotificationCancelSalesOrderResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "NotificationGrpcHandler.SendNotificationCancelSalesOrder")
	defer span.End()

	res = &pb.SendNotificationCancelSalesOrderResponse{}
	res, err = h.ServicesNotificationCancelSalesOrder.SendNotification(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}
