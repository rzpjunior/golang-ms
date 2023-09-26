package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/internal/app/model"
)

type INotificationPurchaserRepository interface {
	Insert(ctx context.Context, notificationPurchaser *model.NotificationPurchaser) (err error)
}

type NotificationPurchaserRepository struct {
	opt opt.Options
}

func NewNotificationPurchaserRepository() INotificationPurchaserRepository {
	return &NotificationPurchaserRepository{
		opt: global.Setup.Common,
	}
}

func (r *NotificationPurchaserRepository) Insert(ctx context.Context, notificationPurchaser *model.NotificationPurchaser) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "NotificationPurchaserRepository.Insert")
	defer span.End()

	db := r.opt.Mongox
	_, err = db.Insert(ctx, "notification_purchaser", notificationPurchaser)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
