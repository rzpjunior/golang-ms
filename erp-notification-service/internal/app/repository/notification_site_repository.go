package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/internal/app/model"
)

type INotificationSiteRepository interface {
	Insert(ctx context.Context, notificationSite *model.NotificationSite) (err error)
}

type NotificationSiteRepository struct {
	opt opt.Options
}

func NewNotificationSiteRepository() INotificationSiteRepository {
	return &NotificationSiteRepository{
		opt: global.Setup.Common,
	}
}

func (r *NotificationSiteRepository) Insert(ctx context.Context, notificationSite *model.NotificationSite) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "NotificationSiteRepository.Insert")
	defer span.End()

	db := r.opt.Mongox
	_, err = db.Insert(ctx, "notification_site", notificationSite)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
