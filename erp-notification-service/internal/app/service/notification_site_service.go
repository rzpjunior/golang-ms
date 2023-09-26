package service

import (
	"context"
	"fmt"
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

type INotificationSiteService interface {
	SendNotification(ctx context.Context, req *pb.SendNotificationHelperRequest) (res *pb.SuccessResponse, err error)
}

type NotificationSiteService struct {
	opt                        opt.Options
	RepositoryNotification     repository.INotificationRepository
	RepositoryNotificationSite repository.INotificationSiteRepository
}

func NewNotificationSiteService() INotificationSiteService {
	return &NotificationSiteService{
		opt:                        global.Setup.Common,
		RepositoryNotification:     repository.NewNotificationRepository(),
		RepositoryNotificationSite: repository.NewNotificationSiteRepository(),
	}
}

func (s *NotificationSiteService) SendNotification(ctx context.Context, req *pb.SendNotificationHelperRequest) (res *pb.SuccessResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "NotificationSiteService.SendNotification")
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

	NP.Title = notification.Title
	NP.Body = notification.Message
	NP.Sound = "default"
	SP.Priority = "high"
	serverKey = s.opt.Env.GetString("firebase.helper_server_key")

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

	notificationSite := &model.NotificationSite{
		ID:        primitive.NewObjectID(),
		HelperID:  req.StaffId,
		RefID:     req.RefId,
		Type:      req.Type,
		Title:     notification.Title,
		Message:   notification.Message,
		CreatedAt: time.Now(),
	}

	span.AddEvent("creating new notification site")
	err = s.RepositoryNotificationSite.Insert(ctx, notificationSite)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &pb.SuccessResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Success: true,
	}

	return
}
