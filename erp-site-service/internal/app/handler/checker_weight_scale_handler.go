package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type CheckerWeightScaleHandler struct {
	Option                     global.HandlerOptions
	ServicesCheckerWeightScale service.ICheckerWeightScaleService
}

// URLMapping implements router.RouteHandlers
func (h *CheckerWeightScaleHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesCheckerWeightScale = service.NewCheckerWeightScaleService()

	cMiddleware := middleware.NewMiddleware()

	r.GET("/:id", h.Get, cMiddleware.Authorized())
	r.PUT("/:id", h.Update, cMiddleware.Authorized())
}

func (h CheckerWeightScaleHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.CheckerWeightScaleGetRequest
	if req.PickingOrderItemId, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesCheckerWeightScale.Get(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h CheckerWeightScaleHandler) Update(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.CheckerWeightScaleUpdateRequest
	if req.PickingOrderItemId, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesCheckerWeightScale.Update(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}
