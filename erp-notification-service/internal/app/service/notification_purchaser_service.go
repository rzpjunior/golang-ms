package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/internal/app/repository"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/notification_service"
	"github.com/NaySoftware/go-fcm"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
)

type INotificationPurchaserService interface {
	SendNotification(ctx context.Context, req *pb.SendNotificationPurchaserRequest) (res *pb.SendNotificationPurchaserResponse, err error)
}

type NotificationPurchaserService struct {
	opt                             opt.Options
	RepositoryNotification          repository.INotificationRepository
	RepositoryNotificationPurchaser repository.INotificationPurchaserRepository
}

func NewNotificationPurchaserService() INotificationPurchaserService {
	return &NotificationPurchaserService{
		opt:                             global.Setup.Common,
		RepositoryNotification:          repository.NewNotificationRepository(),
		RepositoryNotificationPurchaser: repository.NewNotificationPurchaserRepository(),
	}
}

func (s *NotificationPurchaserService) SendNotification(ctx context.Context, req *pb.SendNotificationPurchaserRequest) (res *pb.SendNotificationPurchaserResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "NotificationPurchaserService.SendNotification")
	defer span.End()

	var (
		NP        fcm.NotificationPayload
		SP        fcm.FcmMsg
		serverKey string
	)

	notification, err := s.RepositoryNotification.GetMessageTemplate(ctx, 0, req.NotifCode)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	notification.Message = strings.ReplaceAll(notification.Message, "#purchase_plan_code#", req.RefId)
	notification.Message = strings.ReplaceAll(notification.Message, "#purchasing_manager_name#", req.StaffId)

	NP.Title = notification.Title
	NP.Body = notification.Message
	NP.Sound = "default"
	SP.Priority = "high"
	serverKey = s.opt.Env.GetString("firebase.purchaser_server_key")

	data := map[string]string{
		"id":   req.RefId,
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

	notificationPurchaser := &model.NotificationPurchaser{
		ID:                primitive.NewObjectID(),
		FieldPurcahaserID: req.StaffId,
		RefID:             req.RefId,
		Type:              req.Type,
		Title:             notification.Title,
		Message:           notification.Message,
		CreatedAt:         time.Now(),
	}

	span.AddEvent("creating new notification purchaser")
	err = s.RepositoryNotificationPurchaser.Insert(ctx, notificationPurchaser)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &pb.SendNotificationPurchaserResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Success: true,
	}

	return
}
