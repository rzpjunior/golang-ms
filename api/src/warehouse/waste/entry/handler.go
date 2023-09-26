package entry

import (
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/common/now"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("we_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("we_rdd"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.POST("", h.create, auth.Authorized("we_crt"))
	r.PUT("/cancel/:id", h.cancel, auth.Authorized("we_can"))
	r.PUT("/confirm/:id", h.confirm, auth.Authorized("we_cnf"))
	r.GET("/export/form", h.exportForm, auth.Authorized("we_exp_form"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var data []*model.WasteEntry
	var total int64

	if data, total, e = repository.GetWasteEntrys(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.WasteEntry
	var total int64

	if data, total, e = repository.GetFilterWasteEntrys(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetWasteEntry("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// create : function to create new data based on input
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

func (h *Handler) cancel(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Cancel(r)
			}
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) exportForm(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var backdate time.Time

	isExport := ctx.QueryParam("export") == "1"
	warehouseID, _ := common.Decrypt(ctx.QueryParam("warehouse_id"))

	warehouse, _ := repository.GetWarehouse("id", warehouseID)

	filter := map[string]interface{}{"table": "all", "attribute": "waste_reason"}
	exclude := map[string]interface{}{}
	wasteReason, _, _ := repository.GetGlossariesByFilter(filter, exclude)

	backdate = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time
	data, total, e := repository.GetExportFormWaste(rq, warehouse.ID)
	if e == nil {
		if isExport {
			var file string
			if file, e = ExportFormXls(backdate, data, warehouse, wasteReason); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data, total)
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) confirm(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r confirmRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Confirm(r)
			}
		}
	}

	return ctx.Serve(e)
}
