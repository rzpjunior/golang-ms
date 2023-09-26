package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"

	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/model"
)

func RepositoryMenu() IMenuRepository {
	m := new(MenuRepository)
	m.opt = global.Setup.Common
	return m
}

type IMenuRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (menus []*model.Menu, count int64, err error)
	GetByUserID(ctx context.Context, userID int64) (menus []*model.Menu, err error)
}

type MenuRepository struct {
	opt opt.Options
}

func (r *MenuRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (menus []*model.Menu, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "MenuRepository.Get")
	defer span.End()

	db := r.opt.Database.Read
	qs := db.QueryTable(new(model.Menu))

	cond := orm.NewCondition()

	if search != "" {
		cond = cond.And("title__icontains", search)
	}

	if status != 0 {
		cond = cond.And("status", status)
	}

	qs = qs.SetCond(cond)

	if orderBy != "" {
		qs = qs.OrderBy(orderBy)
	}

	count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &menus)
	if err != nil {
		span.RecordError(err)
		return
	}

	for _, menu := range menus {
		if _, err = db.QueryTable(new(model.Menu)).RelatedSel(1).Filter("status", 1).Filter("parent_id", menu.ID).OrderBy("order").All(&menu.Child); err != nil {
			span.RecordError(err)
			return
		}
	}

	return
}

func (r *MenuRepository) GetByUserID(ctx context.Context, userID int64) (menus []*model.Menu, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "MenuRepository.GetByID")
	defer span.End()

	db := r.opt.Database.Read

	where := "where m.status = 1 and p.status = 1 and m.parent_id is null "
	q := "select m.* " +
		"from menu m " +
		"join permission p on m.permission_id = p.id "

	if userID > 0 {
		join := "join role_permission rp on p.id = rp.permission_id " +
			"join user_role ur on ur.role_id = rp.role_id "
		where = where + "and ur.user_id = ? group by m.id order by m.order asc"
		q = q + join
	}
	q = q + where

	if _, err = db.Raw(q, userID).QueryRows(&menus); err != nil {
		span.RecordError(err)
		return
	}

	for _, menu := range menus {
		where := "where m.status = ? and p.status = ? and m.parent_id = ? "
		q := "select m.* " +
			"from menu m " +
			"join permission p on m.permission_id = p.id "

		if userID > 0 {
			join := "join role_permission rp on p.id = rp.permission_id " +
				"join user_role ur on ur.role_id = rp.role_id "
			where = where + "and ur.user_id = ? group by m.id order by m.order asc"
			q = q + join
		}
		q = q + where

		if _, err = db.Raw(q, 1, 1, menu.ID, userID).QueryRows(&menu.Child); err != nil {
			span.RecordError(err)
			return
		}
	}
	return
}
