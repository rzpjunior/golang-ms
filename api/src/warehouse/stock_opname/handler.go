package stock_opname

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
	r.GET("", h.read, auth.Authorized("st_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("st_rdd"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.POST("", h.create, auth.Authorized("st_crt"))
	r.GET("/export/form", h.exportForm, auth.Authorized("st_exp_form"))
	r.PUT("/cancel/:id", h.cancel, auth.Authorized("st_can"))
	r.PUT("/confirm/:id", h.confirm, auth.Authorized("st_cnf"))
	r.GET("/download/form/:id", h.download, auth.Authorized("st_dl"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var data []*model.StockOpname
	var total int64

	if data, total, e = repository.GetStockOpnames(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.StockOpname
	var total int64

	if data, total, e = repository.GetFilterStockOpnames(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetStockOpname("id", id); e != nil {
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

func (h *Handler) exportForm(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var backdate time.Time

	isExport := ctx.QueryParam("export") == "1"
	warehouseID, _ := common.Decrypt(ctx.QueryParam("warehouse_id"))
	categoryID, _ := common.Decrypt(ctx.QueryParam("category_id"))
	stockTypeID := ctx.QueryParam("stock_type")
	classification := ctx.QueryParam("classification")

	warehouse, _ := repository.GetWarehouse("id", warehouseID)
	category, _ := repository.GetCategory("id", categoryID)
	stockType, _ := repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_int", stockTypeID)

	filter := map[string]interface{}{"table": "stock_opname", "attribute": "opname_reason"}
	exclude := map[string]interface{}{}
	opnameReason, _, _ := repository.GetGlossariesByFilter(filter, exclude)

	backdate = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time
	data, total, e := repository.GetExportForm(rq, warehouse.ID, category.ID, classification)
	if e != nil {
		return ctx.Serve(e)
	}
	if isExport {
		var file string
		if stockType.ValueName == "good stock" {
			if file, e = ExportGoodStockFormXls(backdate, data, warehouse, opnameReason); e != nil {
				return ctx.Serve(e)
			}
		} else {
			if file, e = ExportWasteStockFormXls(backdate, data, warehouse, opnameReason); e != nil {
				return ctx.Serve(e)
			}
		}
		ctx.Files(file)
	} else {
		ctx.Data(data, total)
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

// print : function to print delivery order
func (h *Handler) download(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	var stockOpname *model.StockOpname
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if stockOpname, e = repository.GetStockOpname("id", id); e == nil {
			ctx.ResponseData, e = Download(stockOpname)
		}
	}
	return ctx.Serve(e)
}
