package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"
)

type IUserCustomerRepository interface {
	GetDetail(ctx context.Context, reqCustomer *model.UserCustomer) (userCustomer *model.UserCustomer, err error)
	Update(ctx context.Context, UserCustomer *model.UserCustomer, columns ...string) (err error)
	Create(ctx context.Context, UserCustomer *model.UserCustomer) (err error)
}

type UserCustomerRepository struct {
	opt opt.Options
}

func NewUserCustomerRepository() IUserCustomerRepository {
	return &UserCustomerRepository{
		opt: global.Setup.Common,
	}
}

func (r *UserCustomerRepository) GetDetail(ctx context.Context, reqCustomer *model.UserCustomer) (userCustomer *model.UserCustomer, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserCustomerRepository.GetDetail")
	defer span.End()

	var cols []string
	userCustomer = &model.UserCustomer{}
	if reqCustomer.ID != 0 {
		cols = append(cols, "id")
		userCustomer.ID = reqCustomer.ID
	}

	if reqCustomer.Code != "" {
		cols = append(cols, "code")
		userCustomer.Code = reqCustomer.Code
	}

	if reqCustomer.CustomerID != 0 {
		cols = append(cols, "customer_id")
		userCustomer.CustomerID = reqCustomer.CustomerID
	}

	if reqCustomer.Status != 0 {
		cols = append(cols, "status")
		userCustomer.Status = reqCustomer.Status
	}

	if reqCustomer.LoginToken != "" {
		cols = append(cols, "login_token")
		userCustomer.LoginToken = reqCustomer.LoginToken
	}
	// if reqCustomer.CustomerIDGP != "" {
	// 	cols = append(cols, "customer_id_gp")
	// 	userCustomer.CustomerIDGP = reqCustomer.CustomerIDGP
	// }

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, userCustomer, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	//var merchant *bridgeService.

	return
}

func (r *UserCustomerRepository) Update(ctx context.Context, UserCustomer *model.UserCustomer, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserCustomerRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, UserCustomer, columns...)
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

func (r *UserCustomerRepository) Create(ctx context.Context, UserCustomer *model.UserCustomer) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserCustomerRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, UserCustomer)
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
