package koli_increment

import (
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("pco_prt"))
	r.GET("/id", h.detail, auth.Authorized("pco_prt"))
	r.POST("/print", h.receiverPrint, auth.Authorized("pco_prt"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.DeliveryKoliIncrement
	var total int64

	if data, total, e = repository.GetDeliveryKoliIncrements(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	if ctx.ResponseData, e = repository.GetDeliveryKoliIncrement("id", id); e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}

// receiverPrint : function to print label selected
func (h *Handler) receiverPrint(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r printRequest
	req := make(map[string]interface{})

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)

	}
	req["plis"] = r.DeliveryKoliIncrement
	file := util.SendPrint(req, "read/label_reprint")

	ctx.Files(file)

	return ctx.Serve(e)

}
