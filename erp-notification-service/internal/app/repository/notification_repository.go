package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/internal/app/model"
)

type INotificationRepository interface {
	GetMessageTemplate(ctx context.Context, id int64, code string) (notification *model.Notification, err error)
}

type NotificationRepository struct {
	opt opt.Options
}

func NewNotificationRepository() INotificationRepository {
	return &NotificationRepository{
		opt: global.Setup.Common,
	}
}

func (r *NotificationRepository) GetMessageTemplate(ctx context.Context, id int64, code string) (notification *model.Notification, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "NotificationRepository.GetMessageTemplate")
	defer span.End()

	var cols []string
	notification = &model.Notification{}
	if id != 0 {
		notification.ID = id
		cols = append(cols, "id")
	}

	if code != "" {
		notification.Code = code
		cols = append(cols, "code")
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, notification, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
