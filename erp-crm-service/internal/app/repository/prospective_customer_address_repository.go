package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/model"
)

type IProspectiveCustomerAddressRepository interface {
	Get(ctx context.Context, prospeciveCustomerID int64) (ProspectiveCustomerAddress []*model.ProspectiveCustomerAddress, count int64, err error)
	GetDetail(ctx context.Context, req *dto.ProspectiveCustomerAddressGetDetailRequest) (prospectiveCustomerAddress *model.ProspectiveCustomerAddress, err error)
	Create(ctx context.Context, customer *model.ProspectiveCustomerAddress) (cust *model.ProspectiveCustomerAddress, err error)
	Update(ctx context.Context, customer *model.ProspectiveCustomerAddress, columns ...string) (err error)
}

type ProspectiveCustomerAddressRepository struct {
	opt opt.Options
}

func NewProspectiveCustomerAddressRepository() IProspectiveCustomerAddressRepository {
	return &ProspectiveCustomerAddressRepository{
		opt: global.Setup.Common,
	}
}

func (r *ProspectiveCustomerAddressRepository) Get(ctx context.Context, prospeciveCustomerID int64) (ProspectiveCustomerAddress []*model.ProspectiveCustomerAddress, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ProspectiveCustomerAddressRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.ProspectiveCustomerAddress))

	cond := orm.NewCondition()

	cond = cond.And("prospective_customer_id", prospeciveCustomerID)

	qs = qs.SetCond(cond)

	count, err = qs.AllWithCtx(ctx, &ProspectiveCustomerAddress)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ProspectiveCustomerAddressRepository) GetDetail(ctx context.Context, req *dto.ProspectiveCustomerAddressGetDetailRequest) (prospectiveCustomerAddress *model.ProspectiveCustomerAddress, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ProspectiveCustomerAddressRepository.GetDetail")
	defer span.End()

	prospectiveCustomerAddress = &model.ProspectiveCustomerAddress{}

	var cols []string

	if req.ID != 0 {
		cols = append(cols, "id")
		prospectiveCustomerAddress.ID = req.ID
	}

	if req.ProspectiveCustomerID != 0 {
		cols = append(cols, "prospective_customer_id")
		prospectiveCustomerAddress.ProspectiveCustomerID = req.ProspectiveCustomerID
	}

	if req.AddressType != "" {
		cols = append(cols, "address_type")
		prospectiveCustomerAddress.AddressType = req.AddressType
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, prospectiveCustomerAddress, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *ProspectiveCustomerAddressRepository) Create(ctx context.Context, customer *model.ProspectiveCustomerAddress) (cust *model.ProspectiveCustomerAddress, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ProspectiveCustomerAddress.Create")
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

func (r *ProspectiveCustomerAddressRepository) Update(ctx context.Context, customer *model.ProspectiveCustomerAddress, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ProspectiveCustomerAddress.Update")
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
