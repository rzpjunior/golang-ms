package paymentedn

import (
	"strconv"

	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.PUT("/cancel/:id", h.cancel, auth.Authorized("sp_edn_can"))
	r.GET("/print/:id", h.receivePrint, auth.Authorized("sp_edn_prt"))
}

// cancel : function to unarchive requested data based on parameters
func (h *Handler) cancel(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Cancel(r)

	return ctx.Serve(e)
}

// receivePrint : function to print data based on Sales Payment id with order type EDN Sales
func (h *Handler) receivePrint(c echo.Context) (e error) {

	ctx := c.(*cuxs.Context)
	var (
		sp               *model.SalesPayment
		remainingInvoice float64
		id               int64
		format           int
		file             string
	)
	req := make(map[string]interface{})

	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	format, e = strconv.Atoi(ctx.QueryParam("format"))

	if sp, e = repository.GetSalesPayment("id", id); e != nil {
		return ctx.Serve(e)
	}

	if e = sp.SalesInvoice.Read("ID"); e != nil {
		return ctx.Serve(e)
	}

	if e = sp.SalesInvoice.SalesOrder.Read("ID"); e != nil {
		return ctx.Serve(e)
	}

	if e = sp.SalesInvoice.SalesOrder.Warehouse.Read("ID"); e != nil {
		return ctx.Serve(e)
	}

	if e = sp.SalesInvoice.SalesOrder.Branch.Read("ID"); e != nil {
		return ctx.Serve(e)
	}

	if remainingInvoice, e = repository.CheckRemainingSalesInvoiceAmount(sp.SalesInvoice.ID); e != nil {
		return ctx.Serve(e)
	}

	if remainingInvoice > 0 && sp.SalesInvoice.Status != 2 {
		req["sisa_hutang"] = remainingInvoice
	} else {
		req["sisa_hutang"] = "LUNAS"
	}

	req["company"] = "UD OKI"
	req["sp"] = sp
	req["branch"] = sp.SalesInvoice.SalesOrder.Branch
	req["warehouse"] = sp.SalesInvoice.SalesOrder.Warehouse

	if format == 1 {
		file = util.SendPrint(req, "read/kwitansi_edn")
	} else {
		file = util.SendPrint(req, "read/kwitansi_edn_termal")
	}

	ctx.Files(file)

	return ctx.Serve(e)
}
