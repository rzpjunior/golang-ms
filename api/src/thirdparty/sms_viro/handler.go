package sms_viro

import (
	"git.edenfarm.id/cuxs/cuxs"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.POST("", h.updateStatusSMSViro)
}

//create : function to create new data based on input
func (h *Handler) updateStatusSMSViro(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest

	if e = ctx.Bind(&r); e == nil {
		ctx.Data(updateOTPSmsViro(r))
	}
	return ctx.Serve(e)
}
