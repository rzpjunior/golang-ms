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

func (h *NotificationGrpcHandler) SendNotificationCampaign(ctx context.Context, req *notificationService.SendNotificationCampaignRequest) (res *notificationService.SendNotificationCampaignResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "NotificationGrpcHandler.SendNotificationCampaign")
	defer span.End()

	var userCustomers []*dto.UserCustomer
	for _, v := range req.UserCustomers {
		userCustomers = append(userCustomers, &dto.UserCustomer{
			CustomerID:     v.CustomerId,
			UserCustomerID: v.UserCustomerId,
			FirebaseToken:  v.FirebaseToken,
		})
	}

	param := &dto.SendNotificationCampaignRequest{
		NotificationCampaignID:   req.NotificationCampaignId,
		NotificationCampaignCode: req.NotificationCampaignCode,
		NotificationCampaignName: req.NotificationCampaignName,
		Title:                    req.Title,
		Message:                  req.Message,
		RedirectTo:               req.RedirectTo,
		RedirectToName:           req.RedirectToName,
		RedirectValue:            req.RedirectValue,
		RedirectValueName:        req.RedirectValueName,
		UserCustomer:             userCustomers,
	}
	var notificationStatus *dto.NotificationStatus
	notificationStatus, err = h.ServicesNotificationCampaign.SendNotificationCampaign(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &notificationService.SendNotificationCampaignResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &notificationService.StatusNotificationCampaign{
			SuccessSent: notificationStatus.SuccessSent,
			FailedSent:  notificationStatus.FailedSent,
		},
	}
	return
}

func (h *NotificationGrpcHandler) GetNotificationCampaignList(ctx context.Context, req *notificationService.GetNotificationCampaignListRequest) (res *notificationService.GetNotificationCampaignListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "NotificationGrpcHandler.GetNotificationList")
	defer span.End()

	param := &dto.GetNotificationCampaignRequest{
		CustomerID: req.CustomerId,
		Limit:      req.Limit,
		Offset:     req.Offset,
	}

	var notificatiosCampaigns []*dto.NotificationCampaignResponse
	notificatiosCampaigns, _, err = h.ServicesNotificationCampaign.Get(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*notificationService.NotificationCampaign
	for _, v := range notificatiosCampaigns {
		data = append(data, &notificationService.NotificationCampaign{
			Id:                     v.ID,
			NotificationCampaignId: v.NotificationCampaignID,
			CustomerId:             v.CustomerID,
			UserCustomerId:         v.UserCustomerID,
			FirebaseToken:          v.FirebaseToken,
			RedirectTo:             v.RedirectTo,
			RedirectToName:         v.RedirectToName,
			RedirectValue:          v.RedirectValue,
			RedirectValueName:      v.RedirectValueName,
			Sent:                   v.Sent,
			Opened:                 v.Opened,
			Conversion:             v.Conversion,
			CreatedAt:              timestamppb.New(v.CreatedAt),
			UpdatedAt:              timestamppb.New(v.UpdatedAt),
			RetryCount:             int32(v.RetryCount),
			FcmResultStatus:        v.FcmResultStatus,
		})
	}

	res = &notificationService.GetNotificationCampaignListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *NotificationGrpcHandler) UpdateReadNotificationCampaign(ctx context.Context, req *notificationService.UpdateReadNotificationCampaignRequest) (res *notificationService.UpdateReadNotificationCampaignResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "NotificationGrpcHandler.GetNotificationUpdate")
	defer span.End()

	param := &dto.UpdateReadNotificationCampaignRequest{
		NotificationCampaignID: req.NotificationCampaignId,
		CustomerID:             req.CustomerId,
	}

	err = h.ServicesNotificationCampaign.UpdateRead(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &notificationService.UpdateReadNotificationCampaignResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}

func (h *NotificationGrpcHandler) CountUnreadNotificationCampaign(ctx context.Context, req *notificationService.CountUnreadNotificationCampaignRequest) (res *notificationService.CountUnreadNotificationCampaignResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "NotificationGrpcHandler.GetNotificationUpdate")
	defer span.End()

	param := &dto.CountUnreadNotificationCampaignRequest{
		CustomerID: req.CustomerId,
	}

	var count int64
	count, err = h.ServicesNotificationCampaign.CountUnread(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &notificationService.CountUnreadNotificationCampaignResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    count,
	}
	return
}
