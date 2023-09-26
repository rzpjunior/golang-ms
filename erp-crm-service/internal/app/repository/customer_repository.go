package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/crm_service"
)

type ICustomerRepository interface {
	GetDetail(ctx context.Context, req *dto.CustomerRequestGetDetail) (customer *model.Customer, err error)
	Update(ctx context.Context, customer *model.Customer, columns ...string) (err error)
	Create(ctx context.Context, customer *model.Customer) (cust *model.Customer, err error)
	SyncGP(ctx context.Context, customer *model.Customer) (err error)
	GetCustomerID(ctx context.Context, req *crm_service.GetCustomerIDRequest) (customer []*model.Customer, err error)
}

type CustomerRepository struct {
	opt opt.Options
}

func NewCustomerRepository() ICustomerRepository {
	return &CustomerRepository{
		opt: global.Setup.Common,
	}
}

func (r *CustomerRepository) GetDetail(ctx context.Context, req *dto.CustomerRequestGetDetail) (customer *model.Customer, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerRepository.GetDetail")
	defer span.End()

	customer = &model.Customer{}
	var cols []string

	if req.ID != 0 {
		customer.ID = req.ID
		cols = append(cols, "id")
	}

	if req.CustomerIDGP != "" {
		customer.CustomerIDGP = req.CustomerIDGP
		cols = append(cols, "customer_id_gp")
	}

	if req.Email != "" {
		customer.Email = req.Email
		cols = append(cols, "Email")
	}

	if req.ReferrerCode != "" {
		customer.ReferralCode = req.ReferrerCode
		cols = append(cols, "referral_code")
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, customer, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CustomerRepository) Update(ctx context.Context, customer *model.Customer, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "Customer.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, customer, columns...)

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

func (r *CustomerRepository) SyncGP(ctx context.Context, customer *model.Customer) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemRepository.SyncGP")
	defer span.End()

	db := r.opt.Database.Write

	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, _, err = tx.ReadOrCreateWithCtx(ctx, customer, "customer_id_gp", []string{}...)
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

func (r *CustomerRepository) Create(ctx context.Context, customer *model.Customer) (cust *model.Customer, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "Customer.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	id, err := tx.InsertWithCtx(ctx, customer)

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
	customer.ID = id
	cust = customer
	return
}

func (r *CustomerRepository) GetCustomerID(ctx context.Context, req *crm_service.GetCustomerIDRequest) (customer []*model.Customer, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserCustomerRepository.GetFirebaseToken")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.Customer))

	cond := orm.NewCondition()
	if len(req.CustomerIdGp) != 0 {
		cond = cond.And("customer_id_gp__in", req.CustomerIdGp)
	}

	qs = qs.SetCond(cond)

	qs = qs.OrderBy("-id")
	_, err = qs.AllWithCtx(ctx, &customer)

	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
