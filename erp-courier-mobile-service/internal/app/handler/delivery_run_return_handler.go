package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-courier-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-courier-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-courier-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-courier-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type DeliveryRunReturnHandler struct {
	Option                    global.HandlerOptions
	ServicesDeliveryRunReturn service.IDeliveryRunReturnService
}

// URLMapping implements router.RouteHandlers
func (h *DeliveryRunReturnHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesDeliveryRunReturn = service.NewDeliveryRunReturnService()

	cMiddleware := middleware.NewMiddleware()

	//// Return
	r.POST("/return/:id", h.CreateReturn, cMiddleware.AuthorizedCourierApp())
	r.PUT("/return/:id", h.EditReturn, cMiddleware.AuthorizedCourierApp())
	r.DELETE("/return/:id", h.DeleteReturn, cMiddleware.AuthorizedCourierApp())

}

func (h DeliveryRunReturnHandler) CreateReturn(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.DeliveryReturnRequest

	req.CourierID = ctx.Request().Context().Value(constants.KeyCourierID).(string)
	if req.ID, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	err = h.ServicesDeliveryRunReturn.Create(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h DeliveryRunReturnHandler) EditReturn(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.DeliveryReturnRequest

	req.CourierID = ctx.Request().Context().Value(constants.KeyCourierID).(string)
	if req.ID, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	err = h.ServicesDeliveryRunReturn.Update(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h DeliveryRunReturnHandler) DeleteReturn(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.DeleteDeliveryReturnRequest

	req.CourierID = ctx.Request().Context().Value(constants.KeyCourierID).(string)
	if req.ID, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	err = h.ServicesDeliveryRunReturn.Delete(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}
