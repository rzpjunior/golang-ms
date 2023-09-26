package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-notification-service/internal/app/dto"
	notificationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/notification_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *NotificationGrpcHandler) SendNotificationTransaction(ctx context.Context, req *notificationService.SendNotificationTransactionRequest) (res *notificationService.SendNotificationTransactionResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "NotificationGrpcHandler.GetNotificationUpdate")
	defer span.End()

	param := &dto.SendNotificationTransactionRequest{
		CustomerID: req.CustomerId,
		RefID:      req.RefId,
		Type:       req.Type,
		SendTo:     req.SendTo,
		NotifCode:  req.NotifCode,
		RefCode:    req.RefCode,
	}

	err = h.ServicesNotificationTransaction.SendNotification(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &notificationService.SendNotificationTransactionResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}

func (h *NotificationGrpcHandler) GetNotificationTransactionList(ctx context.Context, req *notificationService.GetNotificationTransactionListRequest) (res *notificationService.GetNotificationTransactionListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "NotificationGrpcHandler.GetNotificationList")
	defer span.End()

	param := &dto.GetNotificationTransactionRequest{
		CustomerID: req.CustomerId,
		Limit:      req.Limit,
		Offset:     req.Offset,
	}

	var notificatiosTransactions []*dto.NotificationTransactionResponse
	notificatiosTransactions, _, err = h.ServicesNotificationTransaction.Get(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*notificationService.NotificationTransaction
	for _, v := range notificatiosTransactions {
		data = append(data, &notificationService.NotificationTransaction{
			Id:         v.ID,
			CustomerId: v.CustomerID,
			RefId:      v.RefID,
			Type:       v.Type,
			Title:      v.Title,
			Message:    v.Message,
			Read:       v.Read,
			CreatedAt:  timestamppb.New(v.CreatedAt),
		})
	}

	res = &notificationService.GetNotificationTransactionListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *NotificationGrpcHandler) UpdateReadNotificationTransaction(ctx context.Context, req *notificationService.UpdateReadNotificationTransactionRequest) (res *notificationService.UpdateReadNotificationTransactionResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "NotificationGrpcHandler.GetNotificationUpdate")
	defer span.End()

	param := &dto.UpdateReadNotificationTransactionRequest{
		RefID:      req.RefId,
		CustomerID: req.CustomerId,
	}

	err = h.ServicesNotificationTransaction.UpdateRead(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &notificationService.UpdateReadNotificationTransactionResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}

func (h *NotificationGrpcHandler) CountUnreadNotificationTransaction(ctx context.Context, req *notificationService.CountUnreadNotificationTransactionRequest) (res *notificationService.CountUnreadNotificationTransactionResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "NotificationGrpcHandler.GetNotificationUpdate")
	defer span.End()

	param := &dto.CountUnreadNotificationTransactionRequest{
		CustomerID: req.CustomerId,
	}

	var count int64
	count, err = h.ServicesNotificationTransaction.CountUnread(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &notificationService.CountUnreadNotificationTransactionResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    count,
	}
	return
}
