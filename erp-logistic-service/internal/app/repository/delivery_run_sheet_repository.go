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

type IDeliveryRunSheetRepository interface {
	Get(ctx context.Context, req dto.DeliveryRunSheetGetRequest) (res []*model.DeliveryRunSheet, count int64, err error)
	GetByID(ctx context.Context, id int64, code string, status ...[]int) (res *model.DeliveryRunSheet, err error)
	Create(ctx context.Context, model *model.DeliveryRunSheet) (err error)
	Update(ctx context.Context, model *model.DeliveryRunSheet, columns ...string) (err error)
}

type DeliveryRunSheetRepository struct {
	opt opt.Options
}

func NewDeliveryRunSheetRepository() IDeliveryRunSheetRepository {
	return &DeliveryRunSheetRepository{
		opt: global.Setup.Common,
	}
}

func (r *DeliveryRunSheetRepository) Get(ctx context.Context, req dto.DeliveryRunSheetGetRequest) (res []*model.DeliveryRunSheet, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryRunSheetRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.DeliveryRunSheet))

	cond := orm.NewCondition()

	if len(req.ArrCourierIDs) > 0 {
		cond = cond.And("courier_id__in", req.ArrCourierIDs)
	}

	if len(req.Status) > 0 {
		cond = cond.And("status__in", req.Status)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &res)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *DeliveryRunSheetRepository) GetByID(ctx context.Context, id int64, code string, status ...[]int) (res *model.DeliveryRunSheet, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryRunSheetRepository.GetByID")
	defer span.End()

	db := r.opt.Database.Read
	drs := model.DeliveryRunSheet{}
	qs := db.QueryTable(new(model.DeliveryRunSheet))

	cond := orm.NewCondition()

	if id != 0 {
		cond = cond.And("id", id)
	}

	if code != "" {
		cond = cond.And("code", code)
	}

	if len(status) > 0 && len(status[0]) > 0 {
		cond = cond.And("status__in", status[0])
	}
	qs = qs.SetCond(cond)
	err = qs.OneWithCtx(ctx, &drs)
	if err != nil {
		span.RecordError(err)
		return
	}
	res = &drs
	return
}

func (r *DeliveryRunSheetRepository) Create(ctx context.Context, model *model.DeliveryRunSheet) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryRunSheetRepository.Create")
	defer span.End()

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

func (r *DeliveryRunSheetRepository) Update(ctx context.Context, model *model.DeliveryRunSheet, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryRunSheetRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, model, columns...)
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

func (r *DeliveryRunSheetRepository) MockDatas() (mockDatas []*model.DeliveryRunSheet) {
	sampleLatitude := -6.1879324
	sampleLongitude := 106.7376164
	mockDatas = append(mockDatas,
		&model.DeliveryRunSheet{
			ID:                1,
			Code:              "DUMMY CODE",
			DeliveryDate:      time.Now(),
			StartedAt:         time.Now(),
			FinishedAt:        time.Now(),
			StartingLatitude:  &sampleLatitude,
			StartingLongitude: &sampleLongitude,
			FinishedLatitude:  &sampleLatitude,
			FinishedLongitude: &sampleLongitude,
			Status:            2,
			CourierID:         "COU0001",
		},
		&model.DeliveryRunSheet{
			ID:                2,
			Code:              "dummy code",
			DeliveryDate:      time.Now(),
			StartedAt:         time.Now(),
			FinishedAt:        time.Now(),
			StartingLatitude:  &sampleLatitude,
			StartingLongitude: &sampleLongitude,
			FinishedLatitude:  &sampleLatitude,
			FinishedLongitude: &sampleLongitude,
			Status:            2,
			CourierID:         "COU0001",
		},
	)

	return mockDatas
}
