package repository

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/model"
)

type ISalesAssignmentRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, startDateFrom time.Time, startDateTo time.Time, endDateFrom time.Time, endDateTo time.Time) (salesAssignments []*model.SalesAssignment, count int64, err error)
	GetByID(ctx context.Context, id int64) (salesAssignment *model.SalesAssignment, err error)
	Create(ctx context.Context, salesAssignment *model.SalesAssignment) (err error)
	Cancel(ctx context.Context, id int64) (err error)
}

type SalesAssignmentRepository struct {
	opt opt.Options
}

func NewSalesAssignmentRepository() ISalesAssignmentRepository {
	return &SalesAssignmentRepository{
		opt: global.Setup.Common,
	}
}

func (r *SalesAssignmentRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, startDateFrom time.Time, startDateTo time.Time, endDateFrom time.Time, endDateTo time.Time) (salesAssignments []*model.SalesAssignment, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.SalesAssignment))

	cond := orm.NewCondition()

	if territoryID != "" {
		cond = cond.And("territory_id_gp", territoryID)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	if timex.IsValid(startDateFrom) && !timex.IsValid(startDateTo) {
		cond = cond.And("start_date", startDateFrom)
	} else {
		if timex.IsValid(startDateFrom) {
			cond = cond.And("start_date__gte", timex.ToStartTime(startDateFrom))
		}

		if timex.IsValid(startDateTo) {
			cond = cond.And("start_date__lte", timex.ToLastTime(startDateTo))
		}
	}

	if timex.IsValid(endDateFrom) && !timex.IsValid(endDateTo) {
		cond = cond.And("end_date", endDateFrom)
	} else {
		if timex.IsValid(endDateFrom) {
			cond = cond.And("end_date__gte", timex.ToStartTime(endDateFrom))
		}

		if timex.IsValid(endDateTo) {
			cond = cond.And("end_date__lte", timex.ToLastTime(endDateTo))
		}
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &salesAssignments)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesAssignmentRepository) GetByID(ctx context.Context, id int64) (salesAssignment *model.SalesAssignment, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentRepository.GetByID")
	defer span.End()

	salesAssignment = &model.SalesAssignment{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, salesAssignment, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesAssignmentRepository) Create(ctx context.Context, salesAssignment *model.SalesAssignment) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, salesAssignment)
	if err != nil {
		span.RecordError(err)
		return tx.Rollback()
	}

	err = tx.Commit()
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesAssignmentRepository) Cancel(ctx context.Context, id int64) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentRepository.Archive")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	salesAssignment := &model.SalesAssignment{
		ID:     id,
		Status: statusx.ConvertStatusName(statusx.Cancelled),
	}

	_, err = tx.UpdateWithCtx(ctx, salesAssignment, "Status")
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return
	}
	err = tx.Commit()

	return
}
