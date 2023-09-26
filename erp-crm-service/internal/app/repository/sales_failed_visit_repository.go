package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/model"
)

type ISalesFailedVisitRepository interface {
	Get(ctx context.Context, offset int, limit int, salesAssignmentItemId int64, failedStatus int32, orderBy string) (salesFailedVisits []*model.SalesFailedVisit, count int64, err error)
	GetByID(ctx context.Context, id int64) (salesFailedVisit *model.SalesFailedVisit, err error)
	Create(ctx context.Context, salesFailedVisit *model.SalesFailedVisit) (id int64, err error)
	Update(ctx context.Context, salesFailedVisit *model.SalesFailedVisit, columns ...string) (err error)
}

type SalesFailedVisitRepository struct {
	opt opt.Options
}

func NewSalesFailedVisitRepository() ISalesFailedVisitRepository {
	return &SalesFailedVisitRepository{
		opt: global.Setup.Common,
	}
}

func (r *SalesFailedVisitRepository) Get(ctx context.Context, offset int, limit int, salesAssignmentItemId int64, failedStatus int32, orderBy string) (salesFailedVisits []*model.SalesFailedVisit, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesFailedVisitRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.SalesFailedVisit))

	cond := orm.NewCondition()

	if salesAssignmentItemId != 0 {
		cond = cond.And("sales_assignment_item_id", salesAssignmentItemId)
	}

	if failedStatus != 0 {
		cond = cond.And("failed_status", salesAssignmentItemId)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &salesFailedVisits)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesFailedVisitRepository) GetByID(ctx context.Context, id int64) (salesFailedVisit *model.SalesFailedVisit, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesFailedVisitRepository.GetByID")
	defer span.End()

	salesFailedVisit = &model.SalesFailedVisit{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, salesFailedVisit, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesFailedVisitRepository) Create(ctx context.Context, salesFailedVisit *model.SalesFailedVisit) (id int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesFailedVisitRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	id, err = tx.InsertWithCtx(ctx, salesFailedVisit)
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

func (r *SalesFailedVisitRepository) Update(ctx context.Context, salesFailedVisit *model.SalesFailedVisit, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesFailedVisitRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, salesFailedVisit, columns...)

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
