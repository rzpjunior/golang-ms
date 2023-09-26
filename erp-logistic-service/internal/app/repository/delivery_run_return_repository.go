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

type IDeliveryRunReturnRepository interface {
	Get(ctx context.Context, req dto.DeliveryRunReturnGetRequest) (deliveryRunReturns []*model.DeliveryRunReturn, count int64, err error)
	GetByID(ctx context.Context, id int64, code string, deliveryRunSheetItemID int64) (deliveryRunReturn *model.DeliveryRunReturn, err error)
	Create(ctx context.Context, model *model.DeliveryRunReturn) (err error)
	Update(ctx context.Context, model *model.DeliveryRunReturn, columns ...string) (err error)
	Delete(ctx context.Context, id int64) (err error)
}

type DeliveryRunReturnRepository struct {
	opt opt.Options
}

func NewDeliveryRunReturnRepository() IDeliveryRunReturnRepository {
	return &DeliveryRunReturnRepository{
		opt: global.Setup.Common,
	}
}

func (r *DeliveryRunReturnRepository) Get(ctx context.Context, req dto.DeliveryRunReturnGetRequest) (deliveryRunReturns []*model.DeliveryRunReturn, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryRunReturnRepository.Get")
	defer span.End()

	// RETURN DUMMY
	dummies := []*model.DeliveryRunReturn{
		{
			ID:                     1,
			Code:                   "dummy code get all DRR",
			TotalPrice:             99999,
			TotalCharge:            100000,
			CreatedAt:              time.Now(),
			DeliveryRunSheetItemID: 1,
		},
	}
	return dummies, 1, nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.DeliveryRunSheet))

	cond := orm.NewCondition()

	if len(req.ArrDeliveryRunSheetItemIDs) > 0 {
		cond = cond.And("delivery_run_sheet_item_id__in", req.ArrDeliveryRunSheetItemIDs)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &deliveryRunReturns)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *DeliveryRunReturnRepository) GetByID(ctx context.Context, id int64, code string, deliveryRunSheetItemID int64) (deliveryRunReturn *model.DeliveryRunReturn, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryRunReturnRepository.GetByID")
	defer span.End()

	deliveryRunReturn = new(model.DeliveryRunReturn)

	var cols []string

	if id != 0 {
		deliveryRunReturn.ID = id
		cols = append(cols, "id")
	}

	if code != "" {
		deliveryRunReturn.Code = code
		cols = append(cols, "code")
	}

	if deliveryRunSheetItemID != 0 {
		deliveryRunReturn.DeliveryRunSheetItemID = deliveryRunSheetItemID
		cols = append(cols, "delivery_run_sheet_item_id")
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, deliveryRunReturn, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *DeliveryRunReturnRepository) Create(ctx context.Context, model *model.DeliveryRunReturn) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryRunReturnRepository.Create")
	defer span.End()

	// RETURN DUMMY
	// return nil

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

func (r *DeliveryRunReturnRepository) Update(ctx context.Context, model *model.DeliveryRunReturn, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryRunReturnRepository.Update")
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
		tx.Rollback()
		return
	}

	return
}

func (r *DeliveryRunReturnRepository) Delete(ctx context.Context, id int64) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DeliveryRunReturnRepository.Delete")
	defer span.End()

	// RETURN DUMMY
	// return nil

	db := r.opt.Database.Write

	qs := db.QueryTable(new(model.DeliveryRunReturn))

	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = qs.Filter("id", id).Delete()
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
