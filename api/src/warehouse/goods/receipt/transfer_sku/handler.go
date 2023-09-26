package transfer_sku

import (
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/src/auth"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.POST("", h.create, auth.Authorized("tfs_crt"))
}

// create : function to create new transfer sku based on input
func (h *Handler) create(c echo.Context) (err error) {
	ctx := c.(*cuxs.Context)
	var r createRequest

	if r.Session, err = auth.UserSession(ctx); err != nil {
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&r); err != nil {
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = Save(r)

	return ctx.Serve(err)
}
