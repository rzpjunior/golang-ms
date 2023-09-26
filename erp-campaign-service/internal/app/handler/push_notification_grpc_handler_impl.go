package handler

import (
	context "context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	dto "git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *CampaignGrpcHandler) GetPushNotificationList(ctx context.Context, req *pb.GetPushNotificationListRequest) (res *pb.GetPushNotificationListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPushNotificationList")
	defer span.End()

	var pushNotifications []*dto.PushNotificationResponse

	pushNotifications, _, err = h.ServicePushNotification.Get(ctx, int(req.Offset), int(req.Limit), 0, "", req.OrderBy, "", time.Time{}, time.Time{})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	var data []*pb.PushNotification
	for _, pushNotification := range pushNotifications {
		data = append(data, &pb.PushNotification{
			ID:             pushNotification.ID,
			Code:           pushNotification.Code,
			CampaignName:   pushNotification.CampaginName,
			Regions:        pushNotification.Regions,
			RegionNames:    pushNotification.RegionNames,
			Archetypes:     pushNotification.Archetypes,
			ArchetypeNames: pushNotification.ArchetypeNames,
			RedirectTo:     int32(pushNotification.RedirectTo),
			RedirectValue:  pushNotification.RedirectValue,
			Title:          pushNotification.Title,
			Message:        pushNotification.Message,
			PushNow:        int32(pushNotification.PushNow),
			ScheduledAt:    timestamppb.New(pushNotification.ScheduledAt),
			SuccessSent:    int32(pushNotification.SuccessSent),
			FailedSent:     int32(pushNotification.FailedSent),
			Opened:         int32(pushNotification.Opened),
			CreatedAt:      timestamppb.New(pushNotification.CreatedAt),
			CreatedBy:      pushNotification.CreatedBy,
			UpdatedAt:      timestamppb.New(pushNotification.UpdatedAt),
			Status:         int32(pushNotification.Status),
			StatusConvert:  statusx.ConvertStatusValue(pushNotification.Status),
		})
	}

	res = &pb.GetPushNotificationListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CampaignGrpcHandler) GetPushNotificationDetail(ctx context.Context, req *pb.GetPushNotificationDetailRequest) (res *pb.GetPushNotificationDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPushNotificationDetail")
	defer span.End()

	var pushNotification *dto.PushNotificationResponse

	pushNotification, err = h.ServicePushNotification.GetDetailMobile(ctx, req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data *pb.PushNotification

	data = &pb.PushNotification{
		ID:            pushNotification.ID,
		Code:          pushNotification.Code,
		CampaignName:  pushNotification.CampaginName,
		Regions:       pushNotification.Regions,
		Archetypes:    pushNotification.Archetypes,
		RedirectTo:    int32(pushNotification.RedirectTo),
		RedirectValue: pushNotification.RedirectValue,
		Title:         pushNotification.Title,
		Message:       pushNotification.Message,
		PushNow:       int32(pushNotification.PushNow),
		ScheduledAt:   timestamppb.New(pushNotification.ScheduledAt),
		SuccessSent:   int32(pushNotification.SuccessSent),
		FailedSent:    int32(pushNotification.FailedSent),
		Opened:        int32(pushNotification.Opened),
		CreatedAt:     timestamppb.New(pushNotification.CreatedAt),
		CreatedBy:     pushNotification.CreatedBy,
		UpdatedAt:     timestamppb.New(pushNotification.UpdatedAt),
		Status:        int32(pushNotification.Status),
		StatusConvert: statusx.ConvertStatusValue(pushNotification.Status),
	}

	res = &pb.GetPushNotificationDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CampaignGrpcHandler) UpdatePushNotification(ctx context.Context, req *pb.UpdatePushNotificationRequest) (res *pb.UpdatePushNotificationResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPushNotificationDetail")
	defer span.End()

	reqUpdateOpened := &dto.PushNotificationRequestUpdateOpened{
		ID:     req.Id,
		Opened: int(req.Opened),
	}
	err = h.ServicePushNotification.UpdateOpened(ctx, reqUpdateOpened)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &pb.UpdatePushNotificationResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}
