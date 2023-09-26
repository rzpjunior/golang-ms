package service

import (
	"context"
	"strconv"
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

type INotificationPaymentReminderService interface {
	SendNotification(ctx context.Context, req *pb.SendNotificationPaymentReminderRequest) (res *pb.SendNotificationPaymentReminderResponse, err error)
}

type NotificationPaymentReminderService struct {
	opt                    opt.Options
	RepositoryNotification repository.INotificationRepository
	RepositoryTransaction  repository.INotificationTransactionRepository
}

func NewNotificationPaymentReminderService() INotificationPaymentReminderService {
	return &NotificationPaymentReminderService{
		opt:                    global.Setup.Common,
		RepositoryNotification: repository.NewNotificationRepository(),
		RepositoryTransaction:  repository.NewNotificationTransactionRepository(),
	}
}

func (s *NotificationPaymentReminderService) SendNotification(ctx context.Context, req *pb.SendNotificationPaymentReminderRequest) (res *pb.SendNotificationPaymentReminderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "NotificationPaymentReminderService.SendNotification")
	defer span.End()

	var (
		NP            fcm.NotificationPayload
		SP            fcm.FcmMsg
		serverKey     string
		templateNotif *model.Notification
	)

	notification, err := s.RepositoryNotification.GetMessageTemplate(ctx, 0, req.NotifCode)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	location, err := time.LoadLocation(req.TimezoneLocation)

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	year, month, day := time.Now().In(location).Date()
	var bulan string

	switch month.String() {
	case "January":
		bulan = "Jan"
	case "February":
		bulan = "Feb"
	case "March":
		bulan = "Mar"
	case "April":
		bulan = "Apr"
	case "May":
		bulan = "Mei"
	case "June":
		bulan = "Jun"
	case "July":
		bulan = "Jul"
	case "August":
		bulan = "Ags"
	case "September":
		bulan = "Sep"
	case "October":
		bulan = "Okt"
	case "November":
		bulan = "Nov"
	case "December":
		bulan = "Des"
	}

	for _, dataSo := range req.Data {
		templateNotif = new(model.Notification)
		templateNotif.Message = notification.Message

		templateNotif.Message = strings.ReplaceAll(templateNotif.Message, "#sales_order_code#", dataSo.SalesOrderCode)
		templateNotif.Message = strings.ReplaceAll(templateNotif.Message, "#current_date#", strconv.Itoa(day)+" "+bulan+" "+strconv.Itoa(year))
		templateNotif.Message = strings.ReplaceAll(templateNotif.Message, "#time_limit#", req.OrderTimeLimit)

		NP.Title = templateNotif.Title
		NP.Body = templateNotif.Message
		NP.Sound = "default"
		SP.Priority = "high"
		serverKey = s.opt.Env.GetString("firebase.cma_server_key")

		data := map[string]string{
			"id":   dataSo.RefId, // ref_id = sales_order_id
			"type": req.Type,
			"code": dataSo.SalesOrderCode,
		}
		c := fcm.NewFcmClient(serverKey)
		c.NewFcmMsgTo(dataSo.SendTo, data)
		c.SetNotificationPayload(&NP)
		status, er := c.Send()

		if er == nil {
			status.PrintResults()
		} else {
			span.RecordError(er)
			s.opt.Logger.AddMessage(log.ErrorLevel, er)
			return
		}
		notificationCancelSalesOrder := &model.NotificationTransaction{
			ID:         primitive.NewObjectID(),
			RefID:      dataSo.RefId,
			CustomerID: dataSo.CustomerId,
			Type:       req.Type,
			Title:      templateNotif.Title,
			Message:    templateNotif.Message,
			CreatedAt:  time.Now(),
		}

		span.AddEvent("creating new notification payment reminder sales order")
		err = s.RepositoryTransaction.Send(ctx, notificationCancelSalesOrder)

		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	res = &pb.SendNotificationPaymentReminderResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Success: true,
	}

	return
}
