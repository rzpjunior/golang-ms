package repository

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/model"
)

type IAddressCoordinateLogRepository interface {
	Get(ctx context.Context, req dto.AddressCoordinateLogGetRequest) (addressCoordinateLogs []*model.AddressCoordinateLog, count int64, err error)
	GetByID(ctx context.Context, id int64) (addressCoordinateLog *model.AddressCoordinateLog, err error)
	Create(ctx context.Context, model *model.AddressCoordinateLog) (err error)
	GetMostTrusted(ctx context.Context, addressID string) (addressCoordinateLogs *model.AddressCoordinateLog, err error)
}

type AddressCoordinateLogRepository struct {
	opt opt.Options
}

func NewAddressCoordinateLogRepository() IAddressCoordinateLogRepository {
	return &AddressCoordinateLogRepository{
		opt: global.Setup.Common,
	}
}

func (r *AddressCoordinateLogRepository) Get(ctx context.Context, req dto.AddressCoordinateLogGetRequest) (addressCoordinateLogs []*model.AddressCoordinateLog, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "AddressCoordinateLogRepository.Get")
	defer span.End()

	// RETURN DUMMY
	sampleLatitude := -6.1879324
	sampleLongitude := 106.7376164
	dummies := []*model.AddressCoordinateLog{
		{
			ID:             1,
			Latitude:       sampleLatitude,
			Longitude:      sampleLongitude,
			LogChannelID:   6,
			MainCoordinate: 0,
			CreatedAt:      time.Now(),
			CreatedBy:      999,
			AddressID:      "dummy address",
			SalesOrderID:   "SO0001",
		},
		{
			ID:             2,
			Latitude:       sampleLatitude,
			Longitude:      sampleLongitude,
			LogChannelID:   6,
			MainCoordinate: 0,
			CreatedAt:      time.Now(),
			CreatedBy:      999,
			AddressID:      "dummy address",
			SalesOrderID:   "SO0001",
		},
	}
	return dummies, 2, nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.DeliveryRunSheetItem))

	cond := orm.NewCondition()

	if len(req.ArrAddressIDs) > 0 {
		cond = cond.And("address_id__in", req.ArrAddressIDs)
	}

	if len(req.ArrSalesOrderIDs) > 0 {
		cond = cond.And("sales_order_id__in", req.ArrSalesOrderIDs)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	if req.GroupBy != "" {
		qs = qs.GroupBy(req.GroupBy)
	}

	count, err = qs.AllWithCtx(ctx, &addressCoordinateLogs)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *AddressCoordinateLogRepository) GetByID(ctx context.Context, id int64) (addressCoordinateLog *model.AddressCoordinateLog, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "AddressCoordinateLogRepository.GetByID")
	defer span.End()

	// RETURN DUMMY
	dummy := &model.AddressCoordinateLog{
		ID:             1,
		Latitude:       -6.1879324,
		Longitude:      106.7376164,
		LogChannelID:   6,
		MainCoordinate: 0,
		CreatedAt:      time.Now(),
		CreatedBy:      999,
		AddressID:      "dummy address",
		SalesOrderID:   "SO0001",
	}
	return dummy, nil

	addressCoordinateLog = &model.AddressCoordinateLog{
		ID: id,
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, addressCoordinateLog, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *AddressCoordinateLogRepository) Create(ctx context.Context, model *model.AddressCoordinateLog) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "AddressCoordinateLogRepository.Create")
	defer span.End()

	// RETURN DUMMY
	return nil

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, model)
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

func (r *AddressCoordinateLogRepository) GetMostTrusted(ctx context.Context, addressID string) (addressCoordinateLogs *model.AddressCoordinateLog, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "AddressCoordinateLogRepository.GetOneMostTrusted")
	defer span.End()

	// RETURN DUMMY
	sampleLatitude := -6.0
	sampleLongitude := 105.0
	dummy := &model.AddressCoordinateLog{
		ID:             1,
		Latitude:       sampleLatitude,
		Longitude:      sampleLongitude,
		LogChannelID:   6,
		MainCoordinate: 0,
		CreatedAt:      time.Now(),
		CreatedBy:      999,
		AddressID:      "dummy address",
		SalesOrderID:   "SO0001",
	}
	return dummy, nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.DeliveryRunSheetItem))

	cond := orm.NewCondition()

	cond = cond.And("address_id", addressID)

	qs = qs.SetCond(cond)

	qs = qs.OrderBy("log_channel_id")

	_, err = qs.Limit(1).AllWithCtx(ctx, &addressCoordinateLogs)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
