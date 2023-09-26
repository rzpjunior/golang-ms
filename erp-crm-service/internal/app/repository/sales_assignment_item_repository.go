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

type ISalesAssignmentItemRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, salespersonID string, submitDateFrom, submitDateTo, startDateFrom, startDateTo, endDateFrom, endDateTo time.Time, task int, outOfRoute int, customerType int32) (salesAssignmentItems []*model.SalesAssignmentItem, count int64, err error)
	GetSubmissions(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, salespersonID string, submitDateFrom time.Time, submitDateTo time.Time, task int, outOfRoute int) (salesAssignmentItems []*model.SalesAssignmentItem, count int64, err error)
	GetGroupBySalespersonID(ctx context.Context, territoryID string, salespersonID string, startDateFrom time.Time, startDateTo time.Time, task []int) (salesAssignmentItems []*model.SalesAssignmentItem, count int64, err error)
	GetByTask(ctx context.Context, status []int, territoryID string, salespersonID string, startDateFrom time.Time, startDateTo time.Time, task int) (salesAssignmentItems []*model.SalesAssignmentItem, count int64, err error)
	GetBySalesAssignmentItemID(ctx context.Context, salesAssignmentID int64, status int, taskType int, finishDateFrom time.Time, finishDateTo time.Time) (salesAssignmentItems []*model.SalesAssignmentItem, err error)
	GetByID(ctx context.Context, id int64) (salesAssignmentItems *model.SalesAssignmentItem, err error)
	Create(ctx context.Context, salesAssignmentItem *model.SalesAssignmentItem) (err error)
	Cancel(ctx context.Context, id int64) (err error)
	GetSingleActiveTask(ctx context.Context, salesPersonID string) (exist bool, err error)
	Update(ctx context.Context, salesAssignmentItem *model.SalesAssignmentItem, columns ...string) (err error)
	GetMultiTaskActive(ctx context.Context, salesPersonID string) (res []*model.SalesAssignmentItem, total int64, err error)
}

type SalesAssignmentItemRepository struct {
	opt opt.Options
}

func NewSalesAssignmentItemRepository() ISalesAssignmentItemRepository {
	return &SalesAssignmentItemRepository{
		opt: global.Setup.Common,
	}
}

func (r *SalesAssignmentItemRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, salespersonID string, submitDateFrom, submitDateTo, startDateFrom, startDateTo, endDateFrom, endDateTo time.Time, task int, outOfRoute int, customerType int32) (salesAssignmentItems []*model.SalesAssignmentItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentItemRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.SalesAssignmentItem))

	cond := orm.NewCondition()

	if salespersonID != "" {
		cond = cond.And("salesperson_id_gp", salespersonID)
	}

	if task != 0 {
		cond = cond.And("task", task)
	}

	if outOfRoute != 0 {
		cond = cond.And("out_of_route", outOfRoute)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	if customerType > 0 {
		cond = cond.And("customer_type", customerType)
	}

	if timex.IsValid(startDateFrom) {
		cond = cond.And("start_date__gte", timex.ToStartTime(startDateFrom))
	}
	if timex.IsValid(startDateTo) {
		cond = cond.And("start_date__lte", timex.ToLastTime(startDateTo))
	}

	if timex.IsValid(endDateFrom) {
		cond = cond.And("end_date__gte", timex.ToStartTime(endDateFrom))
	}
	if timex.IsValid(endDateTo) {
		cond = cond.And("end_date__lte", timex.ToLastTime(endDateTo))
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

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &salesAssignmentItems)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesAssignmentItemRepository) GetSubmissions(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, salespersonID string, submitDateFrom time.Time, submitDateTo time.Time, task int, outOfRoute int) (salesAssignmentItems []*model.SalesAssignmentItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentItemRepository.GetSubmissions")
	defer span.End()

	db := r.opt.Database.Read

	var where []string
	var whereValue []interface{}

	where = append(where, "sai.status IN (?,?) ")
	whereValue = append(whereValue, statusx.ConvertStatusName(statusx.Finished))
	whereValue = append(whereValue, statusx.ConvertStatusName(statusx.Failed))

	if territoryID != "" {
		where = append(where, "sa.territory_id_gp = ? ")
		whereValue = append(whereValue, territoryID)
	}

	if salespersonID != "" {
		where = append(where, "sai.salesperson_id_gp = ? ")
		whereValue = append(whereValue, salespersonID)
	}

	if task != 0 {
		where = append(where, "sai.task = ? ")
		whereValue = append(whereValue, task)
	}

	if outOfRoute != 0 {
		where = append(where, "sai.out_of_route = ? ")
		whereValue = append(whereValue, outOfRoute)
	}

	if status != 0 {
		where = append(where, "sai.status = ? ")
		whereValue = append(whereValue, status)
	}

	if timex.IsValid(submitDateFrom) {
		where = append(where, "sai.submit_date >= ? ")
		whereValue = append(whereValue, timex.ToStartTime(submitDateFrom))
	}
	if timex.IsValid(submitDateTo) {
		where = append(where, "sai.submit_date <= ? ")
		whereValue = append(whereValue, timex.ToLastTime(submitDateTo))
	}

	var whereStr string
	for i, cond := range where {
		if i == 0 {
			whereStr += cond
		} else {
			whereStr += " AND " + cond
		}
	}

	_, err = db.RawWithCtx(ctx, "SELECT sai.* FROM sales_assignment_item sai INNER JOIN sales_assignment sa ON sa.id = sai.sales_assignment_id WHERE "+whereStr+" ORDER BY id DESC", whereValue...).QueryRows(&salesAssignmentItems)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesAssignmentItemRepository) GetGroupBySalespersonID(ctx context.Context, territoryID string, salespersonID string, startDateFrom time.Time, startDateTo time.Time, task []int) (salesAssignmentItems []*model.SalesAssignmentItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentItemRepository.GetGroupBySalespersonID")
	defer span.End()

	db := r.opt.Database.Read

	var where []string
	var whereValue []interface{}

	if territoryID != "" {
		where = append(where, "sa.territory_id_gp = ? ")
		whereValue = append(whereValue, territoryID)
	}

	if salespersonID != "" {
		where = append(where, "sai.salesperson_id_gp = ? ")
		whereValue = append(whereValue, salespersonID)
	}

	if len(task) != 0 {
		var taskWhereIn string
		for i, t := range task {
			if i == 0 {
				taskWhereIn += "?"
			} else {
				taskWhereIn += ", ?"
			}
			whereValue = append(whereValue, t)
		}
		where = append(where, "sai.task IN ("+taskWhereIn+")")
	}

	if timex.IsValid(startDateFrom) {
		where = append(where, "sai.start_date >= ? ")
		whereValue = append(whereValue, timex.ToStartTime(startDateFrom))
	}
	if timex.IsValid(startDateTo) {
		where = append(where, "sai.start_date <= ? ")
		whereValue = append(whereValue, timex.ToLastTime(startDateTo))
	}

	var whereStr string
	for i, cond := range where {
		if i == 0 {
			whereStr += cond
		} else {
			whereStr += " AND " + cond
		}
	}

	_, err = db.RawWithCtx(ctx, "SELECT sai.salesperson_id FROM sales_assignment_item sai INNER JOIN sales_assignment sa ON sa.id = sai.sales_assignment_id WHERE "+whereStr+" GROUP BY sai.salesperson_id ORDER BY salesperson_id DESC", whereValue...).QueryRows(&salesAssignmentItems)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesAssignmentItemRepository) GetByTask(ctx context.Context, status []int, territoryID string, salespersonID string, startDateFrom time.Time, startDateTo time.Time, task int) (salesAssignmentItems []*model.SalesAssignmentItem, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentItemRepository.GetByTask")
	defer span.End()

	db := r.opt.Database.Read

	var where []string
	var whereValue []interface{}

	if len(status) != 0 {
		var whereInParam string
		for i, s := range status {
			if i == 0 {
				whereInParam += "?"
			} else {
				whereInParam += ", ?"
			}
			whereValue = append(whereValue, s)
		}
		where = append(where, "sai.status IN ("+whereInParam+")")
	}

	if territoryID != "" {
		where = append(where, "sa.territory_id_gp = ? ")
		whereValue = append(whereValue, territoryID)
	}

	if salespersonID != "" {
		where = append(where, "sai.salesperson_id_gp = ? ")
		whereValue = append(whereValue, salespersonID)
	}

	if task != 0 {
		where = append(where, "sai.task = ? ")
		whereValue = append(whereValue, task)
	}

	if timex.IsValid(startDateFrom) {
		where = append(where, "sai.start_date >= ? ")
		whereValue = append(whereValue, timex.ToStartTime(startDateFrom))
	}
	if timex.IsValid(startDateTo) {
		where = append(where, "sai.start_date <= ? ")
		whereValue = append(whereValue, timex.ToLastTime(startDateTo))
	}

	var whereStr string
	for i, cond := range where {
		if i == 0 {
			whereStr += cond
		} else {
			whereStr += " AND " + cond
		}
	}

	_, err = db.RawWithCtx(ctx, "SELECT sai.* FROM sales_assignment_item sai INNER JOIN sales_assignment sa ON sa.id = sai.sales_assignment_id WHERE "+whereStr, whereValue...).QueryRows(&salesAssignmentItems)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesAssignmentItemRepository) GetBySalesAssignmentItemID(ctx context.Context, salesAssignmentID int64, status int, taskType int, finishDateFrom time.Time, finishDateTo time.Time) (salesAssignmentItems []*model.SalesAssignmentItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentItemRepository.GetBySalesAssignmentItemID")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.SalesAssignmentItem))

	cond := orm.NewCondition()

	if salesAssignmentID != 0 {
		cond = cond.And("sales_assignment_id", salesAssignmentID)
	}

	if taskType != 0 {
		cond = cond.And("task", taskType)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	if timex.IsValid(finishDateFrom) {
		cond = cond.And("finish_date__gte", timex.ToStartTime(finishDateFrom))
	}
	if timex.IsValid(finishDateTo) {
		cond = cond.And("finish_date__lte", timex.ToLastTime(finishDateTo))
	}

	qs = qs.SetCond(cond)

	_, err = qs.AllWithCtx(ctx, &salesAssignmentItems)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesAssignmentItemRepository) GetByID(ctx context.Context, id int64) (salesAssignmentItem *model.SalesAssignmentItem, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentRepository.GetByID")
	defer span.End()

	salesAssignmentItem = &model.SalesAssignmentItem{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, salesAssignmentItem, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *SalesAssignmentItemRepository) Create(ctx context.Context, salesAssignmentItem *model.SalesAssignmentItem) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentItemRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, salesAssignmentItem)
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

func (r *SalesAssignmentItemRepository) Cancel(ctx context.Context, id int64) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentRepository.Cancel")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	salesAssignmentItem := &model.SalesAssignmentItem{
		ID:     id,
		Status: statusx.ConvertStatusName(statusx.Cancelled),
	}

	_, err = tx.UpdateWithCtx(ctx, salesAssignmentItem, "Status")
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return
	}
	err = tx.Commit()

	return
}

func (r *SalesAssignmentItemRepository) GetSingleActiveTask(ctx context.Context, salesPersonID string) (exist bool, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentRepository.GetSingleActiveTask")
	defer span.End()

	var salesAssignmentItem *model.SalesAssignmentItem
	db := r.opt.Database.Read
	err = db.Raw("SELECT id FROM sales_assignment_item WHERE submit_date is NOT NULL AND finish_date IS NULL AND salesperson_id_gp = ? AND status != 14", salesPersonID).QueryRow(&salesAssignmentItem)
	if err != nil && !errors.Is(err, orm.ErrNoRows) {
		span.RecordError(err)
		return
	}

	if salesAssignmentItem != nil {
		exist = true
	}

	return
}

func (r *SalesAssignmentItemRepository) Update(ctx context.Context, salesAssignmentItem *model.SalesAssignmentItem, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentItemRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, salesAssignmentItem, columns...)

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

func (r *SalesAssignmentItemRepository) GetMultiTaskActive(ctx context.Context, salesPersonID string) (res []*model.SalesAssignmentItem, total int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "SalesAssignmentRepository.CheckActiveTask")
	defer span.End()

	db := r.opt.Database.Read
	total, err = db.Raw("SELECT id FROM sales_assignment_item WHERE submit_date is NOT NULL AND finish_date IS NULL AND salesperson_id_gp = ? AND status != 14", salesPersonID).QueryRows(&res)
	if err != nil && !errors.Is(err, orm.ErrNoRows) {
		span.RecordError(err)
		return
	}

	return
}
