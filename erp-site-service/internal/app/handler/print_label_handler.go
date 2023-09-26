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

type PrintLabelHandler struct {
	Option             global.HandlerOptions
	ServicesPrintLabel service.IPrintLabelService
}

// URLMapping implements router.RouteHandlers
func (h *PrintLabelHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesPrintLabel = service.NewPrintLabelService()

	cMiddleware := middleware.NewMiddleware()

	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/delivery", h.GetDO, cMiddleware.Authorized())

	r.GET("/delivery_koli", h.GetDeliveryKoli, cMiddleware.Authorized())
	r.POST("/reprint", h.ReprintLabel, cMiddleware.Authorized())
}

func (h PrintLabelHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	req := &dto.PrintLabelGetRequest{
		TypePrint: ctx.GetParamString("type_print"),
		Condition: ctx.GetParamString("condition"),
	}

	ctx.ResponseData, err = h.ServicesPrintLabel.Get(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h PrintLabelHandler) ReprintLabel(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req *dto.RePrintLabelGetRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesPrintLabel.ReprintLabel(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h PrintLabelHandler) GetDO(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	req := &dto.PrintLabelGetRequest{
		TypePrint: ctx.GetParamString("type_print"),
		Condition: ctx.GetParamString("condition"),
	}

	ctx.ResponseData, err = h.ServicesPrintLabel.GetDO(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h PrintLabelHandler) GetDeliveryKoli(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	req := &dto.PrintLabelGetRequest{
		// TypePrint: ctx.GetParamString("type_print"),
		Condition: ctx.GetParamString("condition"),
	}

	ctx.ResponseData, err = h.ServicesPrintLabel.GetDeliveryKoli(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}
