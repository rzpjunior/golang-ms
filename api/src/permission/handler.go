// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package permission

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"github.com/labstack/echo/v4"
)

// Handler collection handler for auth.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.permission)
	r.GET("/privilege", h.privilege, auth.Authorized("filter_rdl"))
	r.GET("/user", h.getUserPermission, auth.Authorized("filter_rdl"))
}

// permission endpoint untuk get sesion data yang lagi login.
func (h *Handler) permission(c echo.Context) (e error) {
	var ps []*model.Permission
	ctx := c.(*cuxs.Context)
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	orSelect.Raw("SELECT * FROM permission WHERE parent_id IS NULL AND status = 1").QueryRows(&ps)
	for _, parent := range ps {
		if _, e = orSelect.Raw("SELECT * FROM permission WHERE parent_id =? AND status = 1", parent.ID).QueryRows(&parent.Child); e == nil {
			for _, child := range parent.Child {
				orSelect.Raw("SELECT * FROM permission WHERE parent_id =? AND status = 1", child.ID).QueryRows(&child.GrandChild)
				for _, gc := range child.GrandChild {
					gc.ID = common.Encrypt(gc.ID)
				}
				child.ID = common.Encrypt(child.ID)

			}
		}
	}

	ctx.ResponseData = ps

	return ctx.Serve(e)
}

// privilege endpoint untuk get sesion data yang lagi login.
func (h *Handler) privilege(c echo.Context) (e error) {
	var up []*model.UserPermission
	var priv []string
	var s *auth.SessionData
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	ctx := c.(*cuxs.Context)

	s, e = auth.UserSession(ctx)
	orSelect.Raw("SELECT * FROM user_permission WHERE user_id = ? ", s.Staff.User.ID).QueryRows(&up)
	for _, p := range up {
		priv = append(priv, p.PermissionValue)
	}
	ctx.ResponseData = priv

	return ctx.Serve(e)
}

func (h *Handler) getUserPermission(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	data, total, e := repository.GetUserPermissions(ctx.RequestQuery())
	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}
