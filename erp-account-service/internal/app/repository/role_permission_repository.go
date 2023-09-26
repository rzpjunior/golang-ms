package repository

import (
	"context"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"

	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/model"
)

func RepositoryRolePermssion() IRolePermissionRepository {
	m := new(RolePermissionRepository)
	m.opt = global.Setup.Common
	return m
}

type IRolePermissionRepository interface {
	Get(ctx context.Context, offset int, limit int, search string) (rolePermissions []*model.RolePermission, count int64, err error)
	GetByID(ctx context.Context, id int64) (rolePermission *model.RolePermission, err error)
	GetByRoleID(ctx context.Context, roleID int64) (rolePermission []*model.RolePermission, err error)
	Create(ctx context.Context, rolePermission *model.RolePermission) (err error)
	Update(ctx context.Context, rolePermission *model.RolePermission, columns ...string) (err error)
	Delete(ctx context.Context, id int64) (err error)
}

type RolePermissionRepository struct {
	opt opt.Options
}

func (r *RolePermissionRepository) Get(ctx context.Context, offset int, limit int, search string) (rolePermissions []*model.RolePermission, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "RolePermissionRepository.Get")
	defer span.End()

	db := r.opt.Database.Read
	count, err = db.QueryTable(new(model.RolePermission)).Offset(offset).Limit(limit).AllWithCtx(ctx, &rolePermissions)

	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *RolePermissionRepository) GetByID(ctx context.Context, id int64) (rolePermission *model.RolePermission, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "RolePermissionRepository.GetByID")
	defer span.End()

	rolePermission = &model.RolePermission{
		ID: id,
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, rolePermission, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *RolePermissionRepository) GetByRoleID(ctx context.Context, roleID int64) (rolePermissions []*model.RolePermission, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "RolePermissionRepository.GetByRoleID")
	defer span.End()

	db := r.opt.Database.Read
	_, err = db.QueryTable(new(model.RolePermission)).Filter("role_id", roleID).Filter("status", 1).AllWithCtx(ctx, &rolePermissions)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *RolePermissionRepository) Create(ctx context.Context, rolePermission *model.RolePermission) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "RolePermissionRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, rolePermission)
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

func (r *RolePermissionRepository) Update(ctx context.Context, rolePermission *model.RolePermission, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "RolePermissionRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, rolePermission, columns...)
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

func (r *RolePermissionRepository) Delete(ctx context.Context, id int64) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "RolePermissionRepository.GetByRoleID")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.DeleteWithCtx(ctx, &model.RolePermission{ID: id})
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
