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

type INotificationTransactionRepository interface {
	Get(ctx context.Context, req *dto.GetNotificationTransactionRequest) (notificationTransactions []*model.NotificationTransaction, count int64, err error)
	Send(ctx context.Context, notificationTransaction *model.NotificationTransaction) (err error)
	UpdateRead(ctx context.Context, filter *model.NotificationTransaction) (err error)
	CountUnread(ctx context.Context, filter *model.NotificationTransaction) (count int64, err error)
	GetMessageTemplate(ctx context.Context, id int64, code string) (notification *model.Notification, err error)
}

type NotificationTransactionRepository struct {
	opt opt.Options
}

func NewNotificationTransactionRepository() INotificationTransactionRepository {
	return &NotificationTransactionRepository{
		opt: global.Setup.Common,
	}
}

func (r *NotificationTransactionRepository) Get(ctx context.Context, req *dto.GetNotificationTransactionRequest) (notificationTransaction []*model.NotificationTransaction, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "NotificationTransactionRepository.Get")
	defer span.End()

	db := r.opt.Mongox

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(req.Limit).SetSkip(req.Offset)

	var ret []byte
	ret, err = db.GetByFilter(ctx, "notification_transaction", &model.NotificationTransaction{
		CustomerID: req.CustomerID,
	}, opts)

	if err != nil {
		span.RecordError(err)
		return
	}

	err = json.Unmarshal(ret, &notificationTransaction)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *NotificationTransactionRepository) Send(ctx context.Context, notificationTransaction *model.NotificationTransaction) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "NotificationTransactionRepository.Send")
	defer span.End()

	db := r.opt.Mongox

	err = db.CreateIndex(ctx, "notification_transaction", "ref_id", false)
	err = db.CreateIndex(ctx, "notification_transaction", "customer_id", false)
	_, err = db.Insert(ctx, "notification_transaction", notificationTransaction)

	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *NotificationTransactionRepository) UpdateRead(ctx context.Context, filter *model.NotificationTransaction) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "NotificationTransactionRepository.UpdateRead")
	defer span.End()

	db := r.opt.Mongox

	update := &model.NotificationTransaction{
		Read: 1,
	}

	err = db.UpdateBulk(ctx, "notification_transaction", filter, update)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *NotificationTransactionRepository) CountUnread(ctx context.Context, filter *model.NotificationTransaction) (count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "NotificationTransactionRepository.CountUnread")
	defer span.End()

	db := r.opt.Mongox

	count, err = db.GetCountByFilter(ctx, "notification_transaction", filter, nil)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *NotificationTransactionRepository) GetMessageTemplate(ctx context.Context, id int64, code string) (notification *model.Notification, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "GetMessageTemplate.GetDetail")
	defer span.End()

	var cols []string
	notification = &model.Notification{}
	if id != 0 {
		cols = append(cols, "id")
		notification.ID = id
	}

	if code != "" {
		cols = append(cols, "code")
		notification.Code = code
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, notification, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
