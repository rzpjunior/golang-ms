package repository

import (
	"context"
	"encoding/json"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IPostponeDeliveryLogRepository interface {
	Get(ctx context.Context, deliveryRunSheetItemID int64) (postponeDeliveryLogs []*model.PostponeDeliveryLog, err error)
	Create(ctx context.Context, postponeDeliveryLog *model.PostponeDeliveryLog) (err error)
}

type PostponeDeliveryLogRepository struct {
	opt opt.Options
}

func NewPostponeDeliveryLogRepository() IPostponeDeliveryLogRepository {
	return &PostponeDeliveryLogRepository{
		opt: global.Setup.Common,
	}
}

func (r *PostponeDeliveryLogRepository) Get(ctx context.Context, deliveryRunSheetItemID int64) (postponeDeliveryLogs []*model.PostponeDeliveryLog, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PostponeDeliveryLogRepository.Get")
	defer span.End()

	db := r.opt.Mongox

	opts := options.FindOptions{}

	var ret []byte
	ret, err = db.GetByFilter(ctx, "postpone_delivery_log", bson.M{"delivery_run_sheet_item_id": deliveryRunSheetItemID}, &opts)
	if err != nil {
		span.RecordError(err)
		return
	}

	err = json.Unmarshal(ret, &postponeDeliveryLogs)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PostponeDeliveryLogRepository) Create(ctx context.Context, postponeDeliveryLog *model.PostponeDeliveryLog) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PostponeDeliveryLogRepository.Create")
	defer span.End()

	db := r.opt.Mongox
	
	_, err = db.Insert(ctx, "postpone_delivery_log", postponeDeliveryLog)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
