package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/model"
)

type ICustomerPointSummaryRepository interface {
	GetDetail(ctx context.Context, req *dto.CustomerPointSummaryRequestGetDetail) (CustomerPointSummary *model.CustomerPointSummary, err error)
	Create(ctx context.Context, CustomerPointSummary *model.CustomerPointSummary) (err error)
	Update(ctx context.Context, CustomerPointSummary *model.CustomerPointSummary, columns ...string) (err error)
}

type CustomerPointSummaryRepository struct {
	opt opt.Options
}

func NewCustomerPointSummaryRepository() ICustomerPointSummaryRepository {
	return &CustomerPointSummaryRepository{
		opt: global.Setup.Common,
	}
}

func (r *CustomerPointSummaryRepository) GetDetail(ctx context.Context, req *dto.CustomerPointSummaryRequestGetDetail) (CustomerPointSummary *model.CustomerPointSummary, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerPointSummaryRepository.GetByID")
	defer span.End()

	CustomerPointSummary = &model.CustomerPointSummary{}

	var cols []string

	if req.ID != 0 {
		CustomerPointSummary.ID = req.ID
		cols = append(cols, "id")
	}

	if req.CustomerID != 0 {
		CustomerPointSummary.CustomerID = req.CustomerID
		cols = append(cols, "customer_id")
	}

	if req.SummaryDate != "" {
		CustomerPointSummary.SummaryDate = req.SummaryDate
		cols = append(cols, "summary_date")
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, CustomerPointSummary, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CustomerPointSummaryRepository) Create(ctx context.Context, CustomerPointSummary *model.CustomerPointSummary) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerPointSummaryRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	CustomerPointSummary.ID, err = tx.InsertWithCtx(ctx, CustomerPointSummary)
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

func (r *CustomerPointSummaryRepository) Update(ctx context.Context, CustomerPointSummary *model.CustomerPointSummary, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerPointSummaryRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, CustomerPointSummary, columns...)

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
