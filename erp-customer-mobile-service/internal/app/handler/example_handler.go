package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/middleware"
	"github.com/labstack/echo/v4"
)

type ExampleHandler struct {
	Option global.HandlerOptions
}

// URLMapping declare endpoint with handler function.
func (h *ExampleHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	cMiddleware := middleware.NewMiddleware()

	r.POST("/verifyOtpRegist", h.verifyOtpRegist, cMiddleware.Authorized("public"))
}

func (h ExampleHandler) verifyOtpRegist(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	return ctx.Serve(e)
}
