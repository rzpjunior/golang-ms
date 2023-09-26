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

type MembershipHandler struct {
	Option            global.HandlerOptions
	ServiceMembership service.IMembershipService
	ServiceWRT        service.IWRTService
}

// URLMapping implements router.RouteHandlers
func (h *MembershipHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceMembership = service.NewMembershipService()
	h.ServiceWRT = service.NewWRTService()
	cMiddleware := middleware.NewMiddleware()

	r.POST("", h.GetMembershipList, cMiddleware.Authorized("private"))
	r.POST("/reward_list", h.GetMembershipRewardList, cMiddleware.Authorized("private"))
	r.POST("/reward_detail", h.GetMembershipRewardDetail, cMiddleware.Authorized("private"))
}

// getMembershipList : to get data of membership
func (h MembershipHandler) GetMembershipList(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.RequestGetMembershipList

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = h.ServiceMembership.Get(ctx.Request().Context(), req)
	return ctx.Serve(e)
}

// getMembershipList : to get data of membership
func (h MembershipHandler) GetMembershipRewardList(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.RequestGetRewardList

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = h.ServiceMembership.GetReward(ctx.Request().Context(), req)
	return ctx.Serve(e)
}

// getMembershipDetail : to get data of membership
func (h MembershipHandler) GetMembershipRewardDetail(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.RequestGetRewardList

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = h.ServiceMembership.GetRewardDetail(ctx.Request().Context(), req)
	return ctx.Serve(e)
}
