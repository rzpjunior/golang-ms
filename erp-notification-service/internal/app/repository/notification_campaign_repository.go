package repository

import (
	"context"
	"encoding/json"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/internal/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type INotificationCampaignRepository interface {
	Get(ctx context.Context, req *dto.GetNotificationCampaignRequest) (notificationCampaigns []*model.NotificationCampaign, count int64, err error)
	Insert(ctx context.Context, notificationCampaign *model.NotificationCampaign) (err error)
	UpdateRead(ctx context.Context, filter *model.NotificationCampaign) (err error)
	CountUnread(ctx context.Context, filter *model.NotificationCampaign) (count int64, err error)
}

type NotificationCampaignRepository struct {
	opt opt.Options
}

func NewNotificationCampaignRepository() INotificationCampaignRepository {
	return &NotificationCampaignRepository{
		opt: global.Setup.Common,
	}
}

func (r *NotificationCampaignRepository) Get(ctx context.Context, req *dto.GetNotificationCampaignRequest) (notificationCampaigns []*model.NotificationCampaign, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "NotificationCampaignRepository.Get")
	defer span.End()

	db := r.opt.Mongox

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(req.Limit).SetSkip(req.Offset)

	var ret []byte
	ret, err = db.GetByFilter(ctx, "notification_campaign", &model.NotificationCampaign{
		CustomerID: req.CustomerID,
	}, opts)

	if err != nil {
		span.RecordError(err)
		return
	}

	err = json.Unmarshal(ret, &notificationCampaigns)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *NotificationCampaignRepository) Insert(ctx context.Context, notificationCampaign *model.NotificationCampaign) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "NotificationCampaignRepository.Insert")
	defer span.End()

	db := r.opt.Mongox
	_, err = db.Insert(ctx, "notification_campaign", notificationCampaign)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *NotificationCampaignRepository) UpdateRead(ctx context.Context, filter *model.NotificationCampaign) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "NotificationCampaignRepository.UpdateRead")
	defer span.End()

	db := r.opt.Mongox

	update := &model.NotificationCampaign{
		Opened: 1,
	}

	err = db.UpdateBulk(ctx, "notification_campaign", filter, update)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *NotificationCampaignRepository) CountUnread(ctx context.Context, filter *model.NotificationCampaign) (count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "NotificationCampaignRepository.CountUnread")
	defer span.End()

	db := r.opt.Mongox

	count, err = db.GetCountByFilter(ctx, "notification_campaign", filter, nil)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
