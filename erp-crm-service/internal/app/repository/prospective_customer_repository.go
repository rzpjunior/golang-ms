package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/model"
)

type IProspectiveCustomerRepository interface {
	Get(ctx context.Context, req *dto.ProspectiveCustomerGetRequest) (prospectiveCustomers []*model.ProspectiveCustomer, count int64, err error)
	GetDetail(ctx context.Context, req *dto.ProspectiveCustomerGetDetailRequest) (prospectiveCustomer *model.ProspectiveCustomer, err error)
	Update(ctx context.Context, prospectiveCustomer *model.ProspectiveCustomer, columns ...string) (err error)
	Create(ctx context.Context, prospectiveCustomer *model.ProspectiveCustomer) (propCust *model.ProspectiveCustomer, err error)
}

type ProspectiveCustomerRepository struct {
	opt opt.Options
}

func NewProspectiveCustomerRepository() IProspectiveCustomerRepository {
	return &ProspectiveCustomerRepository{
		opt: global.Setup.Common,
	}
}

func (r *ProspectiveCustomerRepository) Get(ctx context.Context, req *dto.ProspectiveCustomerGetRequest) (prospectiveCustomers []*model.ProspectiveCustomer, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ProspectiveCustomerRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.ProspectiveCustomer))

	cond := orm.NewCondition()

	if req.Search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("pic_order_contact__icontains", req.Search).Or("business_name__icontains", req.Search)
		cond = cond.AndCond(condGroup)
	}

	if req.Status != 0 {
		cond = cond.And("reg_status", req.Status)
	}

	if req.ArchetypeID != "" {
		cond = cond.And("archetype_id_gp", req.ArchetypeID)
	}

	if req.CustomerTypeID != "" {
		cond = cond.And("customer_type_id_gp", req.CustomerTypeID)
	}

	if req.RegionID != "" {
		cond = cond.And("region_id_gp", req.RegionID)
	}

	if req.CustomerID != "" {
		cond = cond.And("customer_id_gp", req.CustomerID)
	}

	if req.RequestBy == "customer" {
		cond = cond.And("salesperson_id_gp__isnull", true)
	} else if req.RequestBy == "salesperson" {
		cond = cond.And("salesperson_id_gp__isnull", false)
		if req.SalesPersonID != "" {
			cond = cond.And("salesperson_id_gp", req.SalesPersonID)
		}
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	_, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &prospectiveCustomers)
	if err != nil {
		span.RecordError(err)
		return
	}

	count, err = qs.Count()

	return
}

func (r *ProspectiveCustomerRepository) GetDetail(ctx context.Context, req *dto.ProspectiveCustomerGetDetailRequest) (prospectiveCustomer *model.ProspectiveCustomer, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ProspectiveCustomerRepository.GetDetail")
	defer span.End()

	prospectiveCustomer = &model.ProspectiveCustomer{}

	var cols []string

	if req.ID != 0 {
		cols = append(cols, "id")
		prospectiveCustomer.ID = req.ID
	}

	if req.Code != "" {
		cols = append(cols, "code")
		prospectiveCustomer.Code = req.Code
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, prospectiveCustomer, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r ProspectiveCustomerRepository) Update(ctx context.Context, prospectiveCustomer *model.ProspectiveCustomer, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ProspectiveCustomerRepository.Decline")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, prospectiveCustomer, columns...)

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

func (r ProspectiveCustomerRepository) Create(ctx context.Context, prospectiveCustomer *model.ProspectiveCustomer) (propCust *model.ProspectiveCustomer, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "Customer.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	id, err := tx.InsertWithCtx(ctx, prospectiveCustomer)

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

	prospectiveCustomer.ID = id
	propCust = prospectiveCustomer
	return
}
