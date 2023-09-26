package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-notification-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/internal/app/repository"
	"github.com/NaySoftware/go-fcm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type INotificationTransactionService interface {
	SendNotification(ctx context.Context, req *dto.SendNotificationTransactionRequest) (err error)
	Get(ctx context.Context, req *dto.GetNotificationTransactionRequest) (res []*dto.NotificationTransactionResponse, total int64, err error)
	UpdateRead(ctx context.Context, req *dto.UpdateReadNotificationTransactionRequest) (err error)
	CountUnread(ctx context.Context, req *dto.CountUnreadNotificationTransactionRequest) (count int64, err error)
}

type NotificationTransactionService struct {
	opt                               opt.Options
	RepositoryNotificationTransaction repository.INotificationTransactionRepository
}

func NewNotificationTransactionService() INotificationTransactionService {
	return &NotificationTransactionService{
		opt:                               global.Setup.Common,
		RepositoryNotificationTransaction: repository.NewNotificationTransactionRepository(),
	}
}

func (s *NotificationTransactionService) SendNotification(ctx context.Context, req *dto.SendNotificationTransactionRequest) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "NotificationTransactionService.SendNotification")
	defer span.End()

	var (
		NP        fcm.NotificationPayload
		SP        fcm.FcmMsg
		serverKey string
	)

	notification, err := s.RepositoryNotificationTransaction.GetMessageTemplate(ctx, 0, req.NotifCode)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	notification.Message = strings.ReplaceAll(notification.Message, "#sales_order_code#", req.RefCode)

	NP.Title = notification.Title
	NP.Body = notification.Message
	NP.Sound = "default"
	SP.Priority = "high"
	serverKey = s.opt.Env.GetString("firebase.cma_server_key")

	data := map[string]string{
		"id":   req.RefID,
		"type": req.Type,
	}
	c := fcm.NewFcmClient(serverKey)
	c.NewFcmMsgTo(req.SendTo, data)
	c.SetNotificationPayload(&NP)
	status, err := c.Send()

	if err == nil {
		status.PrintResults()
	} else {
		fmt.Println(err)
	}

	notificationTransaction := &model.NotificationTransaction{
		ID:         primitive.NewObjectID(),
		CustomerID: req.CustomerID,
		RefID:      req.RefID,
		Type:       req.Type,
		Title:      notification.Title,
		Message:    notification.Message,
		Read:       2,
		CreatedAt:  time.Now(),
		Status:     req.Status,
	}

	span.AddEvent("creating new notification transaction")
	err = s.RepositoryNotificationTransaction.Send(ctx, notificationTransaction)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *NotificationTransactionService) Get(ctx context.Context, req *dto.GetNotificationTransactionRequest) (res []*dto.NotificationTransactionResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "NotificationTransactionService.Get")
	defer span.End()

	var notificationTransactions []*model.NotificationTransaction
	notificationTransactions, total, err = s.RepositoryNotificationTransaction.Get(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, notificationTransaction := range notificationTransactions {
		mongoId := notificationTransaction.ID.Hex()
		res = append(res, &dto.NotificationTransactionResponse{
			ID:         mongoId,
			CustomerID: notificationTransaction.CustomerID,
			RefID:      notificationTransaction.RefID,
			Type:       notificationTransaction.Type,
			Title:      notificationTransaction.Title,
			Message:    notificationTransaction.Message,
			Read:       notificationTransaction.Read,
			CreatedAt:  notificationTransaction.CreatedAt,
		})
	}

	return
}

func (s *NotificationTransactionService) UpdateRead(ctx context.Context, req *dto.UpdateReadNotificationTransactionRequest) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "NotificationTransactionService.UpdateRead")
	defer span.End()

	filter := &model.NotificationTransaction{
		RefID:      req.RefID,
		CustomerID: req.CustomerID,
	}

	span.AddEvent("Update read notification transaction")
	err = s.RepositoryNotificationTransaction.UpdateRead(ctx, filter)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *NotificationTransactionService) CountUnread(ctx context.Context, req *dto.CountUnreadNotificationTransactionRequest) (count int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "NotificationTransactionService.CountUnread")
	defer span.End()

	filter := &model.NotificationTransaction{
		Read:       2,
		CustomerID: req.CustomerID,
	}

	span.AddEvent("Count unread notification transaction")
	count, err = s.RepositoryNotificationTransaction.CountUnread(ctx, filter)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
