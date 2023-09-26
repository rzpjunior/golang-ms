package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/model"
)

type ICustomerPointLogRepository interface {
	Get(ctx context.Context, req *dto.CustomerPointLogRequestGet) (CustomerPointLogs []*model.CustomerPointLog, count int64, err error)
	GetDetail(ctx context.Context, req *dto.CustomerPointLogRequestGetDetail) (CustomerPointLog *model.CustomerPointLog, err error)
	Create(ctx context.Context, CustomerPointLog *model.CustomerPointLog) (err error)
	Update(ctx context.Context, CustomerPointLog *model.CustomerPointLog, columns ...string) (err error)
	GetDetailHistoryMobile(ctx context.Context, req *dto.CustomerPointLogRequestGetDetail) (CustomerPointLog *model.PointHistoryList, err error)
	GetReferralDataCustomer(ctx context.Context, referrerID int64) (referralHistoryReturn *dto.ReferralHistoryReturn, err error)
}

type CustomerPointLogRepository struct {
	opt opt.Options
}

func NewCustomerPointLogRepository() ICustomerPointLogRepository {
	return &CustomerPointLogRepository{
		opt: global.Setup.Common,
	}
}

func (r *CustomerPointLogRepository) Get(ctx context.Context, req *dto.CustomerPointLogRequestGet) (CustomerPointLogs []*model.CustomerPointLog, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerPointLogRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.CustomerPointLog))

	cond := orm.NewCondition()

	if req.CustomerID != 0 {
		cond = cond.And("customer_id", req.CustomerID)
	}

	if req.SalesOrderID != 0 {
		cond = cond.And("sales_order_id", req.SalesOrderID)
	}

	if req.TransactionType != 0 {
		cond = cond.And("transaction_type", req.TransactionType)
	}

	if req.Status != 0 {
		cond = cond.And("status", req.Status)
	}

	if req.CreatedDate != "" {
		cond = cond.And("created_date", req.CreatedDate)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &CustomerPointLogs)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CustomerPointLogRepository) GetDetail(ctx context.Context, req *dto.CustomerPointLogRequestGetDetail) (CustomerPointLog *model.CustomerPointLog, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerPointLogRepository.GetByID")
	defer span.End()

	CustomerPointLog = &model.CustomerPointLog{}

	var cols []string

	if req.ID != 0 {
		CustomerPointLog.ID = req.ID
		cols = append(cols, "id")
	}

	if req.CustomerID != 0 {
		CustomerPointLog.CustomerID = req.CustomerID
		cols = append(cols, "customer_id")
	}

	if req.SalesOrderID != 0 {
		CustomerPointLog.SalesOrderID = req.SalesOrderID
		cols = append(cols, "sales_order_id")
	}

	if req.Status != 0 {
		CustomerPointLog.Status = req.Status
		cols = append(cols, "status")
	}

	if req.TransactionType != 0 {
		CustomerPointLog.TransactionType = req.TransactionType
		cols = append(cols, "transaction_type")
	}

	if req.CreatedDate != "" {
		CustomerPointLog.CreatedDate = req.CreatedDate
		cols = append(cols, "created_date")
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, CustomerPointLog, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CustomerPointLogRepository) GetDetailHistoryMobile(ctx context.Context, req *dto.CustomerPointLogRequestGetDetail) (CustomerPointLog *model.PointHistoryList, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerPointLogRepository.GetByID")
	defer span.End()

	CustomerPointLog = &model.PointHistoryList{}

	db := r.opt.Database.Read

	err = db.Raw(
		"SELECT "+
			"created_date "+
			", mpl.point_value "+
			",CASE WHEN (mpl.status = 1) THEN"+
			"  'ISSUED' "+
			"ELSE "+
			"'DEDUCTION' "+
			"END "+
			"AS `status_type`"+
			", mpl.status "+
			"from customer_point_log mpl "+
			"where customer_id = ? and mpl.sales_order_id = ? "+
			"and mpl.created_date >= DATE_FORMAT(CURRENT_DATE(), ?) - INTERVAL 3 MONTH "+
			"GROUP BY mpl.id order by mpl.id DESC ", req.CustomerID, req.SalesOrderID, "%Y-%m-%d").QueryRow(&CustomerPointLog)
	if err != nil {
		CustomerPointLog = &model.PointHistoryList{}
	}
	return
}

func (r *CustomerPointLogRepository) Create(ctx context.Context, CustomerPointLog *model.CustomerPointLog) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerPointLogRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	CustomerPointLog.ID, err = tx.InsertWithCtx(ctx, CustomerPointLog)
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

func (r *CustomerPointLogRepository) Update(ctx context.Context, CustomerPointLog *model.CustomerPointLog, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerPointLogRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, CustomerPointLog, columns...)

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

func (r *CustomerPointLogRepository) GetReferralDataCustomer(ctx context.Context, referrerID int64) (referralHistoryReturn *dto.ReferralHistoryReturn, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerPointLogRepository.Update")
	defer span.End()

	referralHistoryReturn = new(dto.ReferralHistoryReturn)

	// q := "select name, date(created_at) created_at from merchant where referrer_id = ? order by id desc"

	// if referralHistoryReturn.Summary.TotalReferral, err = o.Raw(q, merchantID).QueryRows(&referralHistoryReturn.Detail.ReferralList); err != nil {
	// 	return nil, err
	// }

	db := r.opt.Database.Read

	q := "select mpl.sales_order_id, mpl.created_date, mpl.referee_id, mpl.campaign_name, mpl.point_value " +
		"from customer_point_log mpl " +
		"where mpl.transaction_type = 3 and mpl.status = 1 and mpl.customer_id = ? " +
		"order by mpl.id desc"

	if _, err = db.Raw(q, referrerID).QueryRows(&referralHistoryReturn.Detail.ReferralPointList); err != nil {
		return nil, err
	}

	for _, v := range referralHistoryReturn.Detail.ReferralPointList {
		referralHistoryReturn.Summary.TotalPoint += v.PointValue
	}

	// return referralHistoryReturn, nil
	return
}
