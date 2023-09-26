package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type GlossaryHandler struct {
	Option          global.HandlerOptions
	ServiceGlossary service.IGlossaryService
}

// URLMapping implements router.RouteHandlers
func (h *GlossaryHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceGlossary = service.NewGlossaryService()

	r.GET("", h.Get)
}

func (h GlossaryHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	table := ctx.GetParamString("table")
	attribute := ctx.GetParamString("attribute")
	valueInt := ctx.GetParamInt("value_int")
	valueName := ctx.GetParamString("value_name")

	var Glossarys []dto.GlossaryResponse
	var total int64
	Glossarys, total, err = h.ServiceGlossary.Get(ctx.Request().Context(), page.Start, page.Limit, table, attribute, valueInt, valueName)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(Glossarys, total, page)

	return ctx.Serve(err)
}
