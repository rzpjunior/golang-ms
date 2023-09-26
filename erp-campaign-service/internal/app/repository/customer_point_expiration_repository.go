package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/model"
)

type ICustomerPointExpirationRepository interface {
	GetDetail(ctx context.Context, id, customerID int64) (CustomerPointExpiration *model.CustomerPointExpiration, err error)
	Create(ctx context.Context, CustomerPointExpiration *model.CustomerPointExpiration) (err error)
	Update(ctx context.Context, CustomerPointExpiration *model.CustomerPointExpiration, columns ...string) (err error)
}

type CustomerPointExpirationRepository struct {
	opt opt.Options
}

func NewCustomerPointExpirationRepository() ICustomerPointExpirationRepository {
	return &CustomerPointExpirationRepository{
		opt: global.Setup.Common,
	}
}

func (r *CustomerPointExpirationRepository) GetDetail(ctx context.Context, id, customerID int64) (CustomerPointExpiration *model.CustomerPointExpiration, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerPointExpirationRepository.GetByID")
	defer span.End()

	CustomerPointExpiration = &model.CustomerPointExpiration{}

	var cols []string

	if id != 0 {
		CustomerPointExpiration.ID = id
		cols = append(cols, "id")
	}

	if customerID != 0 {
		CustomerPointExpiration.CustomerID = customerID
		cols = append(cols, "CustomerID")
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, CustomerPointExpiration, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CustomerPointExpirationRepository) Create(ctx context.Context, CustomerPointExpiration *model.CustomerPointExpiration) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerPointExpirationRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	CustomerPointExpiration.ID, err = tx.InsertWithCtx(ctx, CustomerPointExpiration)
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

func (r *CustomerPointExpirationRepository) Update(ctx context.Context, CustomerPointExpiration *model.CustomerPointExpiration, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerPointExpirationRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, CustomerPointExpiration, columns...)

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
