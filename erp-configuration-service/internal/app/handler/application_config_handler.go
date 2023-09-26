package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type ApplicationConfigHandler struct {
	Option                   global.HandlerOptions
	ServiceApplicationConfig service.IApplicationConfigService
}

// URLMapping implements router.RouteHandlers
func (h *ApplicationConfigHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceApplicationConfig = service.NewApplicationConfigService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/mobile-config", h.Get)
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
	r.PUT("/:id", h.Update, cMiddleware.Authorized("app_upd"))

}

func (h ApplicationConfigHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	status := ctx.GetParamInt("status")
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	application := ctx.GetParamInt("application")
	attribute := ctx.GetParamString("attribute")

	param := &dto.ApplicationConfigRequestGet{
		Status:      int8(status),
		Search:      search,
		OrderBy:     orderBy,
		Application: int8(application),
		Attribute:   attribute,
		Offset:      int32(page.Start),
		Limit:       int32(page.Limit),
	}

	var menus []dto.ApplicationConfigResponse
	var total int64
	menus, total, err = h.ServiceApplicationConfig.Get(ctx.Request().Context(), param)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(menus, total, page)

	return ctx.Serve(err)
}

func (h ApplicationConfigHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var applicationConfig dto.ApplicationConfigResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	applicationConfig, err = h.ServiceApplicationConfig.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = applicationConfig

	return ctx.Serve(err)
}

func (h ApplicationConfigHandler) Update(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.ApplicationConfigRequestUpdate

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServiceApplicationConfig.Update(ctx.Request().Context(), req, id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}
