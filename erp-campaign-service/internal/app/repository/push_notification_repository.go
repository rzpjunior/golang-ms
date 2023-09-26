package repository

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/model"
)

type IPushNotificationRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, regionID string, scheduledAtFrom time.Time, scheduledAtTo time.Time) (pushNotifications []*model.PushNotification, count int64, err error)
	GetByID(ctx context.Context, id int64) (pushNotification *model.PushNotification, err error)
	Create(ctx context.Context, pushNotification *model.PushNotification) (err error)
	Update(ctx context.Context, pushNotification *model.PushNotification, columns ...string) (err error)
}

type PushNotificationRepository struct {
	opt opt.Options
}

func NewPushNotificationRepository() IPushNotificationRepository {
	return &PushNotificationRepository{
		opt: global.Setup.Common,
	}
}

func (r *PushNotificationRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, regionID string, scheduledAtFrom time.Time, scheduledAtTo time.Time) (pushNotifications []*model.PushNotification, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PushNotificationRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.PushNotification))

	cond := orm.NewCondition()

	if search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("campaign_name__icontains", search).Or("code__icontains", search)
		cond = cond.AndCond(condGroup)
	}

	if regionID != "" {
		regionIDStr := regionID
		condGroupRegion := orm.NewCondition()
		condGroupRegion = condGroupRegion.And("regions__icontains", ","+regionIDStr+",").Or("regions__istartswith", regionIDStr+",").Or("regions__iendswith", ","+regionIDStr).Or("regions", regionIDStr)

		cond = cond.AndCond(condGroupRegion)
	}

	if timex.IsValid(scheduledAtFrom) || timex.IsValid(scheduledAtTo) {
		condGroupScheduledAt := orm.NewCondition()
		if timex.IsValid(scheduledAtFrom) {
			condGroupScheduledAt = condGroupScheduledAt.And("scheduled_at__gte", timex.ToStartTime(scheduledAtFrom))
		}
		if timex.IsValid(scheduledAtTo) {
			condGroupScheduledAt = condGroupScheduledAt.And("scheduled_at__lte", timex.ToLastTime(scheduledAtTo))
		}
		cond = cond.AndCond(condGroupScheduledAt)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	_, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &pushNotifications)
	if err != nil {
		span.RecordError(err)
		return
	}

	count, err = qs.Count()

	return
}

func (r *PushNotificationRepository) GetByID(ctx context.Context, id int64) (pushNotification *model.PushNotification, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PushNotificationRepository.GetByID")
	defer span.End()

	pushNotification = &model.PushNotification{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, pushNotification, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PushNotificationRepository) Create(ctx context.Context, pushNotification *model.PushNotification) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PushNotificationRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	pushNotification.ID, err = tx.InsertWithCtx(ctx, pushNotification)
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PushNotificationRepository) Update(ctx context.Context, pushNotification *model.PushNotification, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PushNotificationRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, pushNotification, columns...)

	if err == nil {
		err = tx.Commit()
		if err != nil {
			span.RecordError(err)
			return
		}
	} else {
		span.RecordError(err)
		err = tx.Rollback()
		if err != nil {
			span.RecordError(err)
			return
		}
	}

	return
}
