package repository

import (
	"context"
	"encoding/json"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/dto"
	logisticService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ICourierLogRepository interface {
	Create(ctx context.Context, model *dto.CourierLog) (err error)
	GetLastCourierLog(ctx context.Context, req *logisticService.GetLastCourierLogRequest) (courierLog *dto.CourierLog, err error)
	GetAllCourierLocation(ctx context.Context, arrCourierID []string) (courierLogs []*dto.CourierLog, err error)
}

type CourierLogRepository struct {
	opt opt.Options
}

func NewCourierLogRepository() ICourierLogRepository {
	return &CourierLogRepository{
		opt: global.Setup.Common,
	}
}

func (r *CourierLogRepository) Create(ctx context.Context, model *dto.CourierLog) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CourierLogRepository.Create")
	defer span.End()

	db := r.opt.Mongox
	db.CreateIndex(ctx, "courier_log", "courier_id", false)
	db.CreateIndex(ctx, "courier_log", "created_at", false)

	_, err = db.Insert(ctx, "courier_log", model)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CourierLogRepository) GetLastCourierLog(ctx context.Context, req *logisticService.GetLastCourierLogRequest) (courierLog *dto.CourierLog, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CourierLogRepository.Get")
	defer span.End()

	db := r.opt.Mongox

	opts := options.FindOneOptions{
		Sort: bson.D{{"created_at", -1}},
	}

	var ret []byte
	ret, err = db.FindByFilter(ctx, "courier_log", bson.M{"courier_id": req.CourierId}, &opts)
	if err != nil {
		span.RecordError(err)
		return
	}

	err = json.Unmarshal(ret, &courierLog)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CourierLogRepository) GetAllCourierLocation(ctx context.Context, arrCourierID []string) (courierLogs []*dto.CourierLog, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CourierLogRepository.GetAllCourierLocation")
	defer span.End()

	db := r.opt.Mongox

	twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"courier_id": bson.M{
					"$in": arrCourierID, // Filter by courier IDs
				},
				"created_at": bson.M{
					"$gte": twentyFourHoursAgo, // Filter documents created within the last 24 hours
				},
			},
		},
		{
			"$sort": bson.M{
				"created_at": -1, // Sort by created_at in descending order
			},
		},
		{
			"$group": bson.M{
				"_id":      "$courier_id", // Group by courier_id
				"_idField": bson.M{"$first": "$_id"},
				"latest_document": bson.M{
					"$first": "$$ROOT", // Retrieve the first document within each group
				},
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$latest_document", // Replace the root document with the latest_document
			},
		},
	}

	// Execute the aggregation using FindAggregate function
	result, err := db.FindAggregate(ctx, "courier_log", pipeline)
	if err != nil {
		span.RecordError(err)
		return
	}
	err = json.Unmarshal(result, &courierLogs)
	if err != nil {
		span.RecordError(err)
		return
	}
	return
}
