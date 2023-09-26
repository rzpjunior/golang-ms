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

type TermConditionHandler struct {
	Option                global.HandlerOptions
	ServicesTermCondition service.ITermConditionService
}

// URLMapping implements router.RouteHandlers
func (h *TermConditionHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesTermCondition = service.NewTermConditionService()

	cMiddleware := middleware.NewMiddleware()
	r.POST("", h.Get, cMiddleware.Authorized("public"))
	r.POST("/:id", h.Detail, cMiddleware.Authorized("public"))
	r.POST("/accept", h.AcceptTNC, cMiddleware.Authorized("private"))
}

func (h TermConditionHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var termConditions []dto.TermConditionResponse
	var total int64
	termConditions, total, err = h.ServicesTermCondition.Get(ctx.Request().Context(), page.Start, page.Limit)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(termConditions, total, page)

	return ctx.Serve(err)
}

func (h TermConditionHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var termCondition dto.TermConditionResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	termCondition, err = h.ServicesTermCondition.GetDetail(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = termCondition

	return ctx.Serve(err)
}

func (h TermConditionHandler) AcceptTNC(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	var req *dto.RequestAcceptTNC

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if req.Session, err = service.CustomerSession(ctx); err == nil {
		if err = ctx.Bind(&req); err != nil {
			h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
			return ctx.Serve(err)
		}
	}

	res, err := h.ServicesTermCondition.AcceptTNC(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.Data(res)

	return ctx.Serve(err)
}
