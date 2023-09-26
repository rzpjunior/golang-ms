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

type IMerchantDeliveryLogRepository interface {
	Create(ctx context.Context, merchantDeliveryLog *model.MerchantDeliveryLog) (err error)
	GetFirst(ctx context.Context, deliveryRunSheetItemID int64) (merchantDeliveryLogs *model.MerchantDeliveryLog, err error)
}

type MerchantDeliveryLogRepository struct {
	opt opt.Options
}

func NewMerchantDeliveryLogRepository() IMerchantDeliveryLogRepository {
	return &MerchantDeliveryLogRepository{
		opt: global.Setup.Common,
	}
}

func (r *MerchantDeliveryLogRepository) Create(ctx context.Context, merchantDeliveryLog *model.MerchantDeliveryLog) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "MerchantDeliveryLogRepository.Create")
	defer span.End()

	db := r.opt.Mongox
	_, err = db.Insert(ctx, "merchant_delivery_log", merchantDeliveryLog)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *MerchantDeliveryLogRepository) GetFirst(ctx context.Context, deliveryRunSheetItemID int64) (merchantDeliveryLogs *model.MerchantDeliveryLog, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "MerchantDeliveryLogRepository.Get")
	defer span.End()

	db := r.opt.Mongox

	opts := options.FindOneOptions{}

	var ret []byte
	ret, err = db.FindByFilter(ctx, "merchant_delivery_log", bson.M{"delivery_run_sheet_item_id": deliveryRunSheetItemID}, &opts)
	if err != nil {
		span.RecordError(err)
		return
	}

	err = json.Unmarshal(ret, &merchantDeliveryLogs)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
