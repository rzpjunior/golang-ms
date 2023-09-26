package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type GlossaryHandler struct {
	Option           global.HandlerOptions
	ServicesGlossary service.IGlossaryService
}

// URLMapping implements router.RouteHandlers
func (h *GlossaryHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesGlossary = service.NewServiceGlossary()

	cMiddleware := middleware.NewMiddleware()

	r.GET("", h.GetGlossary, cMiddleware.AuthorizedHelperMobile())
}

func (h GlossaryHandler) GetGlossary(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.GetGlossaryRequest

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	req.Table = ctx.GetParamString("table")
	req.Attribute = ctx.GetParamString("attribute")
	req.ValueInt = ctx.GetParamInt("value_int")
	req.ValueName = ctx.GetParamString("value_name")

	glossaries, total, err := h.ServicesGlossary.GetGlossary(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.DataList(glossaries, total, page)

	return ctx.Serve(err)
}
