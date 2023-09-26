package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/notification_service"
)

type INotificationCampaignService interface {
	GetHistoryCampaign(ctx context.Context, req *dto.NotificationCampaignRequestGet) (res []*dto.NotificationCampaignResponse, err error)
	UpdateRead(ctx context.Context, req *dto.NotificationCampaignRequestUpdateRead) (err error)
	CountUnread(ctx context.Context, req *dto.NotificationCampaignRequestCountUnread) (res *dto.NotificationCampaignCountUnreadResponse, err error)
}

type NotificationCampaignService struct {
	opt opt.Options
}

func NewNotificationCampaignService() INotificationCampaignService {
	return &NotificationCampaignService{
		opt: global.Setup.Common,
	}
}

func (s *NotificationCampaignService) GetHistoryCampaign(ctx context.Context, req *dto.NotificationCampaignRequestGet) (res []*dto.NotificationCampaignResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "NotificationCampaignService.Get")
	defer span.End()

	notificationCampaigns, err := s.opt.Client.NotificationServiceGrpc.GetNotificationCampaignList(ctx, &notification_service.GetNotificationCampaignListRequest{
		CustomerId: req.Session.Customer.ID,
		Limit:      req.Limit,
		Offset:     req.Offset * req.Limit,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, notif := range notificationCampaigns.Data {
		var notificationDetail *campaign_service.GetPushNotificationDetailResponse
		// Default
		if notif.NotificationCampaignId == "" {
			notif.NotificationCampaignId = "1"
		}
		notificationDetail, err = s.opt.Client.CampaignServiceGrpc.GetPushNotificationDetail(ctx, &campaign_service.GetPushNotificationDetailRequest{
			Id: utils.ToInt64(notif.NotificationCampaignId),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		res = append(res, &dto.NotificationCampaignResponse{
			NotificationCampaignID:   notif.NotificationCampaignId,
			NotificationCampaignName: notificationDetail.Data.CampaignName,
			Title:                    notificationDetail.Data.Title,
			Message:                  notificationDetail.Data.Message,
			RedirectTo:               notif.RedirectTo,
			RedirectToName:           notif.RedirectToName,
			RedirectValue:            notif.RedirectValue,
			RedirectValueName:        notif.RedirectValueName,
			Sent:                     notif.Sent,
			Opened:                   notif.Opened,
			CreatedAt:                notif.CreatedAt.AsTime(),
		})
	}

	return
}

func (s *NotificationCampaignService) UpdateRead(ctx context.Context, req *dto.NotificationCampaignRequestUpdateRead) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "NotificationCampaignService.Get")
	defer span.End()

	_, err = s.opt.Client.NotificationServiceGrpc.UpdateReadNotificationCampaign(ctx, &notification_service.UpdateReadNotificationCampaignRequest{
		CustomerId:             req.Session.Customer.ID,
		NotificationCampaignId: req.Data.NotificationCampaignID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	notificationDetail, err := s.opt.Client.CampaignServiceGrpc.GetPushNotificationDetail(ctx, &campaign_service.GetPushNotificationDetailRequest{
		Id: utils.ToInt64(req.Data.NotificationCampaignID),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	_, err = s.opt.Client.CampaignServiceGrpc.UpdatePushNotification(ctx, &campaign_service.UpdatePushNotificationRequest{
		Id:     utils.ToInt64(req.Data.NotificationCampaignID),
		Opened: notificationDetail.Data.Opened + 1,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *NotificationCampaignService) CountUnread(ctx context.Context, req *dto.NotificationCampaignRequestCountUnread) (res *dto.NotificationCampaignCountUnreadResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "NotificationCampaignService.Get")
	defer span.End()

	var countUnreadNotif *notification_service.CountUnreadNotificationCampaignResponse
	countUnreadNotif, err = s.opt.Client.NotificationServiceGrpc.CountUnreadNotificationCampaign(ctx, &notification_service.CountUnreadNotificationCampaignRequest{
		CustomerId: req.Session.Customer.ID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.NotificationCampaignCountUnreadResponse{
		Unread: countUnreadNotif.Data,
	}
	return
}
