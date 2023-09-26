package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"

	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/model"
)

func RepositoryUser() IUserRepository {
	m := new(UserRepository)
	m.opt = global.Setup.Common
	return m
}

type IUserRepository interface {
	Get(ctx context.Context, req *dto.GetUserRequest) (users []*model.User, count int64, err error)
	GetByID(ctx context.Context, id int64) (user *model.User, err error)
	GetByEmail(ctx context.Context, email string) (user *model.User, err error)
	GetByEmployeeCode(ctx context.Context, employeeCode string) (user *model.User, err error)
	GetByDivisionID(ctx context.Context, divisionID int64) (users []*model.User, err error)
	Create(ctx context.Context, user *model.User) (err error)
	Update(ctx context.Context, user *model.User, columns ...string) (err error)
	GetBySalesAppLoginToken(ctx context.Context, token string) (user *model.User, err error)
	GetByEdnAppLoginToken(ctx context.Context, token string) (user *model.User, err error)
}

type UserRepository struct {
	opt opt.Options
}

func (r *UserRepository) Get(ctx context.Context, req *dto.GetUserRequest) (users []*model.User, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.User))

	cond := orm.NewCondition()

	if req.Search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("name__icontains", req.Search).Or("nickname__icontains", req.Search).Or("email__icontains", req.Search).Or("employee_code__icontains", req.Search)
		cond = cond.AndCond(condGroup)
	}

	if req.Status != 0 {
		cond = cond.And("status", req.Status)
	}

	if req.SiteID != "" {
		cond = cond.And("site_id_gp", req.SiteID)
	}

	if req.RegionID != "" {
		cond = cond.And("region_id_gp", req.RegionID)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	if req.Apps == "mob-purchaser" {
		var idPurchasers []int64
		db.Raw("select u.id from `user` u inner join user_role ur on u.id =ur.user_id inner join `role` r on r.id = ur.role_id where r.code = 'ROL0050'").QueryRows(&idPurchasers)

		count, err = qs.Filter("id__in", idPurchasers).Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &users)
	} else {
		count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &users)
	}

	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (user *model.User, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserRepository.GetByID")
	defer span.End()

	user = &model.User{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, user, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (user *model.User, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserRepository.GetByEmail")
	defer span.End()

	user = &model.User{
		Email: email,
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, user, "email")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *UserRepository) GetByEmployeeCode(ctx context.Context, employeeCode string) (user *model.User, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserRepository.GetByEmployeeCode")
	defer span.End()

	user = &model.User{
		EmployeeCode: employeeCode,
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, user, "employee_code")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *UserRepository) GetByDivisionID(ctx context.Context, divisionID int64) (users []*model.User, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserRepository.GetByDivisionID")
	defer span.End()

	db := r.opt.Database.Read

	q := "SELECT u.* FROM user u JOIN user_role ur ON ur.user_id = u.id JOIN role r ON r.id = ur.role_id WHERE u.status = 1 AND r.division_id = ? "

	db.RawWithCtx(ctx, q, divisionID).QueryRows(&users)

	return
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, user)

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

func (r *UserRepository) Update(ctx context.Context, user *model.User, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, user, columns...)

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

func (r *UserRepository) GetBySalesAppLoginToken(ctx context.Context, token string) (user *model.User, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserRepository.GetBySalesAppLoginToken")
	defer span.End()

	user = &model.User{
		SalesAppLoginToken: token,
		Status:             1,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, user, "salesapp_login_token", "status")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *UserRepository) GetByEdnAppLoginToken(ctx context.Context, token string) (user *model.User, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "UserRepository.GetByEdnAppLoginToken")
	defer span.End()

	user = &model.User{
		EdnAppLoginToken: token,
		Status:           1,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, user, "edn_app_login_token", "status")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
