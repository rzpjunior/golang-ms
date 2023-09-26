package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type SMSViroHandler struct {
	Option         global.HandlerOptions
	ServiceSMSViro service.ISMSViroService
}

// URLMapping declare endpoint with handler function.
func (h *SMSViroHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	cMiddleware := middleware.NewMiddleware()
	h.ServiceSMSViro = service.NewSMSViroService()

	r.POST("", h.updateStatusSMSViro, cMiddleware.Authorized("public"))
}

func (h SMSViroHandler) updateStatusSMSViro(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.UpdateRequestSMSViro

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = h.ServiceSMSViro.Update(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	return ctx.Serve(e)
}
