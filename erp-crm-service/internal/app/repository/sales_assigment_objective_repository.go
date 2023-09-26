package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/model"
)

type ISalesAssignmentObjectiveRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, codes []string, orderBy string) (salesAssignmentObjectives []*model.SalesAssignmentObjective, count int64, err error)
	GetByID(ctx context.Context, id int64) (salesAssignment *model.SalesAssignmentObjective, err error)
	GetByCode(ctx context.Context, code string) (salesAssignment *model.SalesAssignmentObjective, err error)
	Create(ctx context.Context, salesAssignmentObjective *model.SalesAssignmentObjective) (err error)
	Update(ctx context.Context, salesAssignmentObjective *model.SalesAssignmentObjective, columns ...string) (err error)
}

type SalesAssignmentObjectiveRepository struct {
	opt opt.Options
}

func NewSalesAssignmentObjectiveRepository() ISalesAssignmentObjectiveRepository {
	return &SalesAssignmentObjectiveRepository{
		opt: global.Setup.Common,
	}
}

func (r *SalesAssignmentObjectiveRepository) Get(ctx context.Context, offset int, limit int, status int, search string, codes []string, orderBy string) (salesAssignmentObjectives []*model.SalesAssignmentObjective, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentObjectiveRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.SalesAssignmentObjective))

	cond := orm.NewCondition()

	if search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("code__icontains", search).Or("name__icontains", search)
		cond = cond.AndCond(condGroup)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	if codes != nil {
		condCodes := orm.NewCondition()
		condCodes = condCodes.And("code__in", codes)
		cond = cond.AndCond(condCodes)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &salesAssignmentObjectives)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesAssignmentObjectiveRepository) GetByID(ctx context.Context, id int64) (salesAssignmentObjective *model.SalesAssignmentObjective, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentObjectiveRepository.GetByID")
	defer span.End()

	db := r.opt.Database.Read

	salesAssignmentObjective = &model.SalesAssignmentObjective{
		ID: id,
	}
	err = db.ReadWithCtx(ctx, salesAssignmentObjective, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesAssignmentObjectiveRepository) GetByCode(ctx context.Context, code string) (salesAssignmentObjective *model.SalesAssignmentObjective, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentObjectiveRepository.GetByID")
	defer span.End()

	db := r.opt.Database.Read

	salesAssignmentObjective = &model.SalesAssignmentObjective{
		Code: code,
	}
	err = db.ReadWithCtx(ctx, salesAssignmentObjective, "code")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesAssignmentObjectiveRepository) Create(ctx context.Context, salesAssignmentObjective *model.SalesAssignmentObjective) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentObjectiveRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, salesAssignmentObjective)
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

func (r *SalesAssignmentObjectiveRepository) Update(ctx context.Context, salesAssignmentObjective *model.SalesAssignmentObjective, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentObjectiveRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, salesAssignmentObjective, columns...)

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
