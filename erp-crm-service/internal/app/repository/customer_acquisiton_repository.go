package repository

import (
	"context"
	"errors"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/model"
)

type ICustomerAcquisitionRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, salespersonID string, submitDateFrom time.Time, submitDateTo time.Time) (customerAcquisitions []*model.CustomerAcquisition, count int64, err error)
	GetSubmissions(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, salespersonID string, submitDateFrom time.Time, submitDateTo time.Time) (customerAcquisitions []*model.CustomerAcquisition, count int64, err error)
	GetPerformances(ctx context.Context, territoryID string, salespersonID string, submitDateFrom time.Time, submitDateTo time.Time) (customerAcquisitions []*model.CustomerAcquisition, count int64, err error)
	GetByID(ctx context.Context, id int64) (customerAcquisition *model.CustomerAcquisition, err error)
	GetByCode(ctx context.Context, code string) (customerAcquisition *model.CustomerAcquisition, err error)
	GetByTerritoryID(ctx context.Context, territoryID string) (customerAcquisitions []*model.CustomerAcquisition, count int64, err error)
	GetSingleActiveTask(ctx context.Context, salesPersonID string) (exist bool, err error)
	Update(ctx context.Context, ca *model.CustomerAcquisition, columns ...string) (err error)
	GetMultiTaskActive(ctx context.Context, salesPersonID string) (res []*model.CustomerAcquisition, total int64, err error)
	Create(ctx context.Context, customerAcq *model.CustomerAcquisition) (err error)
	CreateWithItem(ctx context.Context, customerAcq *model.CustomerAcquisition, item []*model.CustomerAcquisitionItem) (id int64, err error)
	GetWithExcludedIds(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, salespersonID string, submitDateFrom time.Time, submitDateTo time.Time, excludedIds []int64) (customerAcquisitions []*model.CustomerAcquisition, count int64, err error)
	CountCustomerAcquisition(ctx context.Context, salespersonID int64, submitDateFrom time.Time, submitDateTo time.Time) (count int64, err error)
}

type CustomerAcquisitionRepository struct {
	opt opt.Options
}

func NewCustomerAcquisitionRepository() ICustomerAcquisitionRepository {
	return &CustomerAcquisitionRepository{
		opt: global.Setup.Common,
	}
}

func (r *CustomerAcquisitionRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, salespersonID string, submitDateFrom time.Time, submitDateTo time.Time) (customerAcquisitions []*model.CustomerAcquisition, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerAcquisitionRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.CustomerAcquisition))

	cond := orm.NewCondition()

	if search != "" {
		cond = cond.And("name__icontains", search)
	}

	if territoryID != "" {
		cond = cond.And("territory_id_gp", territoryID)
	}

	if salespersonID != "" {
		cond = cond.And("salesperson_id_gp", salespersonID)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	if timex.IsValid(submitDateFrom) {
		cond = cond.And("submit_date__gte", timex.ToStartTime(submitDateFrom))
	}
	if timex.IsValid(submitDateTo) {
		cond = cond.And("submit_date__lte", timex.ToLastTime(submitDateTo))
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &customerAcquisitions)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CustomerAcquisitionRepository) GetSubmissions(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, salespersonID string, submitDateFrom time.Time, submitDateTo time.Time) (customerAcquisitions []*model.CustomerAcquisition, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerAcquisitionRepository.GetSubmissions")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.CustomerAcquisition))

	cond := orm.NewCondition()

	if search != "" {
		cond = cond.And("name__icontains", search)
	}

	if territoryID != "" {
		cond = cond.And("territory_id_gp", territoryID)
	}

	if salespersonID != "" {
		cond = cond.And("salesperson_id_gp", salespersonID)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	if timex.IsValid(submitDateFrom) {
		cond = cond.And("submit_date__gte", timex.ToStartTime(submitDateFrom))
	}
	if timex.IsValid(submitDateTo) {
		cond = cond.And("submit_date__lte", timex.ToLastTime(submitDateTo))
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Filter("status", statusx.ConvertStatusName(statusx.Finished)).Offset(offset).Limit(limit).AllWithCtx(ctx, &customerAcquisitions)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CustomerAcquisitionRepository) GetPerformances(ctx context.Context, territoryID string, salespersonID string, submitDateFrom time.Time, submitDateTo time.Time) (customerAcquisitions []*model.CustomerAcquisition, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerAcquisitionRepository.GetPerformances")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.CustomerAcquisition))

	cond := orm.NewCondition()

	if territoryID != "" {
		cond = cond.And("territory_id_gp", territoryID)
	}

	if salespersonID != "" {
		cond = cond.And("salesperson_id_gp", salespersonID)
	}

	cond = cond.And("status__in", statusx.ConvertStatusName(statusx.Active), statusx.ConvertStatusName(statusx.Finished))

	if timex.IsValid(submitDateFrom) {
		cond = cond.And("submit_date__gte", timex.ToStartTime(submitDateFrom))
	}
	if timex.IsValid(submitDateTo) {
		cond = cond.And("submit_date__lte", timex.ToLastTime(submitDateTo))
	}

	qs = qs.SetCond(cond)

	count, err = qs.Filter("status", statusx.ConvertStatusName(statusx.Finished)).AllWithCtx(ctx, &customerAcquisitions)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CustomerAcquisitionRepository) GetByID(ctx context.Context, id int64) (customerAcquisition *model.CustomerAcquisition, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerAcquisitionRepository.GetByID")
	defer span.End()

	customerAcquisition = &model.CustomerAcquisition{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, customerAcquisition, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CustomerAcquisitionRepository) GetByCode(ctx context.Context, code string) (customerAcquisition *model.CustomerAcquisition, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerAcquisitionRepository.GetByCode")
	defer span.End()

	customerAcquisition = &model.CustomerAcquisition{
		Code: code,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, customerAcquisition, "code")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CustomerAcquisitionRepository) GetByTerritoryID(ctx context.Context, territoryID string) (customerAcquisitions []*model.CustomerAcquisition, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerAcquisitionRepository.GetByTerritoryID")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.CustomerAcquisition))

	cond := orm.NewCondition()

	if territoryID != "" {
		cond = cond.And("territory_id_gp", territoryID)
	}
	qs = qs.SetCond(cond)

	count, err = qs.AllWithCtx(ctx, &customerAcquisitions)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CustomerAcquisitionRepository) GetSingleActiveTask(ctx context.Context, salesPersonID string) (exist bool, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerAcquisitionRepository.GetSingleActiveTask")
	defer span.End()

	var ca *model.CustomerAcquisition
	db := r.opt.Database.Read
	err = db.Raw("SELECT id FROM customer_acquisition WHERE submit_date is NOT NULL AND finish_date IS NULL AND salesperson_id_gp = ?", salesPersonID).QueryRow(&ca)
	if err != nil && !errors.Is(err, orm.ErrNoRows) {
		span.RecordError(err)
		return
	}

	if ca != nil {
		exist = true
	}

	return
}

func (r *CustomerAcquisitionRepository) Update(ctx context.Context, ca *model.CustomerAcquisition, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerAcquisitionRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, ca, columns...)

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

func (r *CustomerAcquisitionRepository) GetMultiTaskActive(ctx context.Context, salesPersonID string) (res []*model.CustomerAcquisition, total int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerAcquisitionRepository.GetSingleActiveTask")
	defer span.End()

	db := r.opt.Database.Read
	total, err = db.Raw("SELECT id FROM customer_acquisition WHERE submit_date is NOT NULL AND finish_date IS NULL AND salesperson_id_gp = ?", salesPersonID).QueryRows(&res)
	if err != nil && !errors.Is(err, orm.ErrNoRows) {
		span.RecordError(err)
		return
	}

	return
}

func (r *CustomerAcquisitionRepository) Create(ctx context.Context, customerAcq *model.CustomerAcquisition) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerAcquisitionRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, customerAcq)
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

func (r *CustomerAcquisitionRepository) CreateWithItem(ctx context.Context, customerAcq *model.CustomerAcquisition, item []*model.CustomerAcquisitionItem) (caId int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerAcquisitionRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	caId, err = tx.InsertWithCtx(ctx, customerAcq)
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return
	}

	for _, cai := range item {
		cai.CustomerAcquisitionID = caId
	}

	if _, err = tx.InsertMultiWithCtx(ctx, 100, &item); err != nil {
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

func (r *CustomerAcquisitionRepository) GetWithExcludedIds(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, salespersonID string, submitDateFrom time.Time, submitDateTo time.Time, excludedIds []int64) (customerAcquisitions []*model.CustomerAcquisition, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerAcquisitionRepository.GetWithExcludedIds")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.CustomerAcquisition))

	cond := orm.NewCondition()

	if search != "" {
		cond = cond.And("name__icontains", search)
	}

	if territoryID != "" {
		cond = cond.And("territory_id_gp", territoryID)
	}

	if salespersonID != "" {
		cond = cond.And("salesperson_id_gp", salespersonID)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	if timex.IsValid(submitDateFrom) {
		cond = cond.And("submit_date__gte", timex.ToStartTime(submitDateFrom))
	}

	if timex.IsValid(submitDateTo) {
		cond = cond.And("submit_date__lte", timex.ToLastTime(submitDateTo))
	}
	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	if excludedIds != nil {
		qs.Exclude("id__in", excludedIds)
	}
	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &customerAcquisitions)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CustomerAcquisitionRepository) CountCustomerAcquisition(ctx context.Context, salespersonID int64, submitDateFrom time.Time, submitDateTo time.Time) (count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerAcquisitionRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.CustomerAcquisition))

	cond := orm.NewCondition()

	if timex.IsValid(submitDateFrom) {
		cond = cond.And("submit_date__gte", timex.ToStartTime(submitDateFrom))
	}
	if timex.IsValid(submitDateTo) {
		cond = cond.And("submit_date__lte", timex.ToLastTime(submitDateTo))
	}

	cond = cond.And("salesperson_id", salespersonID)

	qs = qs.SetCond(cond)

	count, err = qs.Count()
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
