package migration

import (
	"git.edenfarm.id/cuxs/cuxs"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("/eden-point", h.migrationEdenPoint)
}

// migrationEdenPoint : function to trigger migration eden point of merchant
func (h *Handler) migrationEdenPoint(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	limit := ctx.QueryParam("limit")
	offset := ctx.QueryParam("offset")
	update := ctx.QueryParam("update")

	if e = migrationEdenPoint(limit, offset, update); e != nil {
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}
