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

type CustomerHandler struct {
	Option          global.HandlerOptions
	ServiceCustomer service.ICustomerService
	ServiceWRT      service.IWRTService
}

// URLMapping implements router.RouteHandlers
func (h *CustomerHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceCustomer = service.NewCustomerService()
	h.ServiceWRT = service.NewWRTService()
	cMiddleware := middleware.NewMiddleware()

	r.POST("/membership", h.GetCustomerMembership, cMiddleware.Authorized("private"))
	r.POST("/referral_history", h.GetCustomerReferralHistory, cMiddleware.Authorized("private"))
}

func (h CustomerHandler) GetCustomerMembership(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.RequestGetPostSession

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = h.ServiceCustomer.GetCustomerMembership(ctx.Request().Context(), req)
	return ctx.Serve(e)
}

func (h CustomerHandler) GetCustomerReferralHistory(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.RequestGetPostSession

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = h.ServiceCustomer.GetReferralHistory(ctx.Request().Context(), req)
	return ctx.Serve(e)
}
