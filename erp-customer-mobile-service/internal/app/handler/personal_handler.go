package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type PersonalHandler struct {
	Option          global.HandlerOptions
	ServicePersonal service.IPersonalService
}

// URLMapping declare endpoint with handler function.
func (h *PersonalHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicePersonal = service.NewPersonalService()
	cMiddleware := middleware.NewMiddleware()
	r.POST("", h.create, cMiddleware.Authorized("public"))

}

func (h PersonalHandler) create(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.SaveRegistrationRequest

	req.Platform = ctx.Request().Header.Get("Platform")
	req.AppVersion = ctx.Request().Header.Get("appVersion")

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = h.ServicePersonal.SaveRegistration(ctx, &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}
