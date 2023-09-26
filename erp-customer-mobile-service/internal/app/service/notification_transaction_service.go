package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/notification_service"
)

type INotificationTransactionService interface {
	GetHistoryTransaction(ctx context.Context, req *dto.NotificationTransactionRequestGet) (res []*dto.NotificationTransactionResponse, err error)
	UpdateRead(ctx context.Context, req *dto.NotificationTransactionRequestUpdateRead) (err error)
	CountUnread(ctx context.Context, req *dto.NotificationTransactionRequestCountUnread) (res *dto.NotificationTransactionCountUnreadResponse, err error)
}

type NotificationTransactionService struct {
	opt opt.Options
}

func NewNotificationTransactionService() INotificationTransactionService {
	return &NotificationTransactionService{
		opt: global.Setup.Common,
	}
}

func (s *NotificationTransactionService) GetHistoryTransaction(ctx context.Context, req *dto.NotificationTransactionRequestGet) (res []*dto.NotificationTransactionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "NotificationTransactionService.Get")
	defer span.End()

	notificationTransactions, err := s.opt.Client.NotificationServiceGrpc.GetNotificationTransactionList(ctx, &notification_service.GetNotificationTransactionListRequest{
		CustomerId: req.Session.Customer.ID,
		Limit:      req.Limit,
		Offset:     req.Offset * req.Limit,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, notif := range notificationTransactions.Data {
		res = append(res, &dto.NotificationTransactionResponse{
			ID:       notif.CustomerId,
			RefId:    notif.RefId,
			Type:     notif.Type,
			Title:    notif.Title,
			Message:  notif.Message,
			Read:     notif.Read,
			CreateAt: notif.CreatedAt.AsTime().Format("2006-01-02 15:04:05"),
		})
	}

	return
}

func (s *NotificationTransactionService) UpdateRead(ctx context.Context, req *dto.NotificationTransactionRequestUpdateRead) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "NotificationTransactionService.Get")
	defer span.End()

	_, err = s.opt.Client.NotificationServiceGrpc.UpdateReadNotificationTransaction(ctx, &notification_service.UpdateReadNotificationTransactionRequest{
		CustomerId: req.Session.Customer.ID,
		RefId:      req.Data.RefId,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *NotificationTransactionService) CountUnread(ctx context.Context, req *dto.NotificationTransactionRequestCountUnread) (res *dto.NotificationTransactionCountUnreadResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "NotificationTransactionService.Get")
	defer span.End()

	var countUnreadNotif *notification_service.CountUnreadNotificationTransactionResponse
	countUnreadNotif, err = s.opt.Client.NotificationServiceGrpc.CountUnreadNotificationTransaction(ctx, &notification_service.CountUnreadNotificationTransactionRequest{
		CustomerId: req.Session.Customer.ID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.NotificationTransactionCountUnreadResponse{
		Unread: countUnreadNotif.Data,
	}

	return
}
