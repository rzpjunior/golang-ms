package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"

	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/model"
)

func RepositoryRole() IRoleRepository {
	m := new(RoleRepository)
	m.opt = global.Setup.Common
	return m
}

type IRoleRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, divisionID int64) (roles []*model.Role, count int64, err error)
	GetByID(ctx context.Context, id int64) (role *model.Role, err error)
	GetByName(ctx context.Context, name string) (role *model.Role, err error)
	Create(ctx context.Context, role *model.Role) (err error)
	Update(ctx context.Context, role *model.Role, columns ...string) (err error)
}

type RoleRepository struct {
	opt opt.Options
}

func (r *RoleRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, divisionID int64) (roles []*model.Role, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "RoleRepository.Get")
	defer span.End()

	db := r.opt.Database.Read
	qs := db.QueryTable(new(model.Role))

	cond := orm.NewCondition()

	if search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("name__icontains", search).Or("code__icontains", search)
		cond = cond.AndCond(condGroup)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	if divisionID != 0 {
		cond = cond.And("division_id", divisionID)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &roles)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *RoleRepository) GetByID(ctx context.Context, id int64) (role *model.Role, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "RoleRepository.GetByID")
	defer span.End()

	role = &model.Role{
		ID: id,
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, role, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *RoleRepository) GetByName(ctx context.Context, name string) (role *model.Role, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "RoleRepository.GetByID")
	defer span.End()

	role = &model.Role{
		Name: name,
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, role, "name")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *RoleRepository) Create(ctx context.Context, role *model.Role) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "RoleRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, role)
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

func (r *RoleRepository) Update(ctx context.Context, role *model.Role, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "RoleRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, role, columns...)
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
