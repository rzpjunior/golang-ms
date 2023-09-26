// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"
)

// Handler collection handler for auth.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.Authorized("usr_rdl"))
	r.GET("/staff", h.getStaff, auth.Authorized("usr_rdl"))
	r.GET("/filter", h.getFiltered, auth.Authorized("filter_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("usr_rdd"))
	r.GET("/staff/:id", h.detailStaff, auth.Authorized("usr_rdd"))
	r.GET("/supervisor", h.getSupervisor, auth.Authorized("usr_rdl"))
	r.GET("/supervisor/filter", h.getSupervisorFilter, auth.Authorized("filter_rdl"))
	r.POST("", h.create, auth.Authorized("usr_crt"))
	r.PUT("/:id", h.update, auth.Authorized("usr_upd"))
	r.PUT("/staff/warehouse-access/:id", h.updateWarehouseAccess, auth.Authorized("usr_upd_wrh_acc"))
	r.PUT("/reset/:id", h.resetPassword, auth.Authorized("usr_rst_pass"))
	r.PUT("/update/permission/:id", h.updatePermission, auth.Authorized("usr_upd_pms"))
	r.PUT("/archive/:id", h.archive, auth.Authorized("usr_arc"))
	r.PUT("/unarchive/:id", h.unArchive, auth.Authorized("usr_urc"))
	r.PUT("/delete/:id", h.delete, auth.Authorized("usr_del"))

	r.POST("/helper", h.createHelper, auth.Authorized("hlp_crt"))
	r.GET("/helper", h.getHelper, auth.Authorized("hlp_rdl"))
	r.GET("/helper/:id", h.getHelperDetail, auth.Authorized("filter_rdl"))
	r.PUT("/helper/:id", h.updateHelper, auth.Authorized("hlp_upd"))
	r.PUT("/helper/archive/:id", h.archiveHelper, auth.Authorized("hlp_arc"))
	r.PUT("/helper/unarchive/:id", h.unarchiveHelper, auth.Authorized("hlp_urc"))

	r.GET("/field_purchaser", h.getFieldPurchaser, auth.Authorized("filter_rdl"))

	//r.GET("/helper/assigned-item", h.getPackingOrderItemAssigned, auth.Authorized("filter_rdl"))

	//PACKING ORDER
	r.GET("/helper/mobile", h.getHelper, auth.AuthorizedMobileUniversal())
	r.GET("/helper/assigned-item/mobile", h.getPackingOrderItemAssigned, auth.AuthorizedMobile())

}

func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	data, total, e := repository.GetUsers(ctx.RequestQuery())
	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

func (h *Handler) getFiltered(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	data, total, e := repository.GetFilterStaff(ctx.RequestQuery())
	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

func (h *Handler) getStaff(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	data, total, e := repository.GetStaffs(ctx.RequestQuery())
	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData, e = Save(r)
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetUser("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) detailStaff(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetStaff("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest

	if r.ID, e = ctx.Decrypt("id"); e == nil {
		if r.Session, e = auth.UserSession(ctx); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = r.Update()
			}
		}
	}
	return ctx.Serve(e)
}

func (h *Handler) updateWarehouseAccess(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateWarehouseAccessRequest

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = r.UpdateWarehouseAccess()

	return ctx.Serve(e)

}

func (h *Handler) resetPassword(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r resetPasswordRequest

	if r.ID, e = ctx.Decrypt("id"); e == nil {
		if r.Session, e = auth.UserSession(ctx); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = r.Reset()
			}
		}
	}
	return ctx.Serve(e)
}

func (h *Handler) updatePermission(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updatePermissionRequest

	if r.ID, e = ctx.Decrypt("id"); e == nil {
		if r.Session, e = auth.UserSession(ctx); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = r.UpdatePermission()
			}
		}
	}
	return ctx.Serve(e)
}

// archive : function to set status of active data into archive
func (h *Handler) archive(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r archiveRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Archive(r)
			}
		}
	}

	return ctx.Serve(e)
}

// unArchive : function to set status of archived data into active
func (h *Handler) unArchive(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r unarchiveRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = UnArchive(r)
			}
		}
	}

	return ctx.Serve(e)
}

// delete : function to set status of data into deleted
func (h *Handler) delete(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r deleteRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Delete(r)
			}
		}
	}

	return ctx.Serve(e)
}

// getSupervisor : function to get supervisor of staff
func (h *Handler) getSupervisor(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	data, total, e := repository.GetSupervisor(ctx.RequestQuery())
	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// getSupervisorFilter : function to get supervisor of staff
func (h *Handler) getSupervisorFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	data, total, e := repository.GetSupervisorFilter(ctx.RequestQuery())
	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// Helper

func (h *Handler) getHelper(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	data, total, e := repository.GetHelpers(ctx.RequestQuery())
	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

func (h *Handler) getHelperDetail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetStaff("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) createHelper(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createHelperRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData, e = SaveHelper(r)
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) updateHelper(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateHelperRequest

	if r.ID, e = ctx.Decrypt("id"); e == nil {
		if r.Session, e = auth.UserSession(ctx); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = r.UpdateHelper()
			}
		}
	}
	return ctx.Serve(e)
}

func (h *Handler) archiveHelper(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r archiveHelperRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = ArchiveHelper(r)
			}
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) unarchiveHelper(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r unarchiveHelperRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = UnarchiveHelper(r)
			}
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) getPackingOrderItemAssigned(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	data, total, e := repository.GetPackingOrderItemAssignedToPacker(ctx.RequestQuery())
	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

func (h *Handler) getFieldPurchaser(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Staff
	var total int64
	var warehouseID int64

	warehouse := ctx.QueryParam("warehouse")

	if warehouse != "" {
		if warehouseID, e = common.Decrypt(warehouse); e != nil {
			return ctx.Serve(e)
		}
	}

	if data, total, e = repository.GetFieldPurchaser(rq, warehouseID); e != nil {
		ctx.Serve(e)
	}

	ctx.Data(data, total)

	return ctx.Serve(e)
}
