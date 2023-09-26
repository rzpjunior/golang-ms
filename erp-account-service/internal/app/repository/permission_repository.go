package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"

	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/model"
)

func RepositoryPermission() IPermissionRepository {
	m := new(PermissionRepository)
	m.opt = global.Setup.Common
	return m
}

type IPermissionRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (permissions []*model.Permission, count int64, err error)
	GetByID(ctx context.Context, id int64) (permission *model.Permission, err error)
	GetTree(ctx context.Context) (permissions []*model.Permission, err error)
	GetPrivilege(ctx context.Context, userID int64) (permissions []string, err error)
	GetTreeByRoleID(ctx context.Context, roleID int64) (permissions []*model.Permission, err error)
}

type PermissionRepository struct {
	opt opt.Options
}

func (r *PermissionRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (permissions []*model.Permission, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PermissionRepository.Get")
	defer span.End()

	db := r.opt.Database.Read
	qs := db.QueryTable(new(model.Permission))

	cond := orm.NewCondition()

	if search != "" {
		cond = cond.And("name__icontains", search)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &permissions)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PermissionRepository) GetByID(ctx context.Context, id int64) (permission *model.Permission, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PermissionRepository.GetByID")
	defer span.End()

	permission = &model.Permission{
		ID: id,
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, permission, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *PermissionRepository) GetTree(ctx context.Context) (permissions []*model.Permission, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PermissionRepository.GetTree")
	defer span.End()

	db := r.opt.Database.Read

	_, err = db.QueryTable(new(model.Permission)).Filter("parent_id__isnull", true).AllWithCtx(ctx, &permissions)
	if err != nil {
		span.RecordError(err)
		return
	}

	for _, permision := range permissions {
		if _, err = db.RawWithCtx(ctx, "SELECT id, parent_id, name, value, status, created_at, updated_at FROM permission WHERE parent_id = ? AND status = 1", permision.ID).QueryRows(&permision.Child); err != nil {
			span.RecordError(err)
			return
		}

		for _, child := range permision.Child {
			if _, err = db.RawWithCtx(ctx, "SELECT id, parent_id, name, value, status, created_at, updated_at FROM permission WHERE parent_id = ? AND status = 1", child.ID).QueryRows(&child.GrandChild); err != nil {
				span.RecordError(err)
				return
			}
		}
	}

	return
}

func (r *PermissionRepository) GetTreeByRoleID(ctx context.Context, roleID int64) (permissions []*model.Permission, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PermissionRepository.GetTreeByRoleID")
	defer span.End()

	db := r.opt.Database.Read

	q := "SELECT p.id, p.parent_id, p.name, p.value, p.status, p.created_at, p.updated_at FROM permission p "
	q += "INNER JOIN role_permission rp ON rp.permission_id = p.id WHERE rp.role_id = ?"

	_, err = db.RawWithCtx(ctx, q+" AND p.parent_id IS NULL", roleID).QueryRows(&permissions)
	if err != nil {
		span.RecordError(err)
		return
	}

	for _, permision := range permissions {
		if _, err = db.RawWithCtx(ctx, q+" AND p.parent_id = ? AND p.status = 1", roleID, permision.ID).QueryRows(&permision.Child); err != nil {
			span.RecordError(err)
			return
		}

		for _, child := range permision.Child {
			if _, err = db.RawWithCtx(ctx, q+" AND p.parent_id = ? AND p.status = 1", roleID, child.ID).QueryRows(&child.GrandChild); err != nil {
				span.RecordError(err)
				return
			}
		}
	}

	return
}

func (r *PermissionRepository) GetPrivilege(ctx context.Context, userID int64) (permissions []string, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "PermissionRepository.GetPrivilge")
	defer span.End()

	db := r.opt.Database.Read
	db.RawWithCtx(ctx, "select p.value from role_permission rp join `role` r on r.id = rp.role_id  join permission p on rp.permission_id = p.id join user_role ur on ur.role_id = r.id  WHERE ur.user_id = ? GROUP BY p.value", userID).QueryRows(&permissions)

	return
}
