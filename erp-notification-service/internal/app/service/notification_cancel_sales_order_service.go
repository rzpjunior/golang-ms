package service

import (
	"context"
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

type INotificationCancelSalesOrderService interface {
	SendNotification(ctx context.Context, req *pb.SendNotificationCancelSalesOrderRequest) (res *pb.SendNotificationCancelSalesOrderResponse, err error)
}

type NotificationCancelSalesOrderService struct {
	opt                    opt.Options
	RepositoryNotification repository.INotificationRepository
	RepositoryTransaction  repository.INotificationTransactionRepository
}

func NewNotificationCancelSalesOrderService() INotificationCancelSalesOrderService {
	return &NotificationCancelSalesOrderService{
		opt:                    global.Setup.Common,
		RepositoryNotification: repository.NewNotificationRepository(),
		RepositoryTransaction:  repository.NewNotificationTransactionRepository(),
	}
}

func (s *NotificationCancelSalesOrderService) SendNotification(ctx context.Context, req *pb.SendNotificationCancelSalesOrderRequest) (res *pb.SendNotificationCancelSalesOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "NotificationCancelSalesOrderService.SendNotification")
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

	for _, dataSo := range req.Data {
		notification.Message = strings.ReplaceAll(notification.Message, "#sales_order_code#", dataSo.SalesOrderCode)

		NP.Title = notification.Title
		NP.Body = notification.Message
		NP.Sound = "default"
		SP.Priority = "high"
		serverKey = s.opt.Env.GetString("firebase.cma_server_key")

		data := map[string]string{
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
			Title:      notification.Title,
			Message:    notification.Message,
			CreatedAt:  time.Now(),
		}

		span.AddEvent("creating new notification cancel sales order")
		err = s.RepositoryTransaction.Send(ctx, notificationCancelSalesOrder)

		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	res = &pb.SendNotificationCancelSalesOrderResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Success: true,
	}

	return
}
