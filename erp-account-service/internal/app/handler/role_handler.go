package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type RoleHandler struct {
	Option global.HandlerOptions
}

// URLMapping declare endpoint with handler function.
func (h *RoleHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	r.GET("", h.Get, middleware.NewMiddleware().Authorized())
	r.GET("/:id", h.Detail, middleware.NewMiddleware().Authorized())
	r.POST("", h.Create, middleware.NewMiddleware().Authorized("rol_crt"))
	r.PUT("/:id", h.Update, middleware.NewMiddleware().Authorized("rol_upd"))
	r.PUT("/archive/:id", h.Archive, middleware.NewMiddleware().Authorized("rol_arc"))
	r.PUT("/unarchive/:id", h.UnArchive, middleware.NewMiddleware().Authorized("rol_urc"))
}

func (h RoleHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sRole := service.ServiceRole()

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
	divisionID := ctx.GetParamInt("division_id")

	var roles []*dto.RoleResponse
	var total int64
	roles, total, err = sRole.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, int64(divisionID))
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(roles, total, page)

	return ctx.Serve(err)
}

func (h RoleHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sRole := service.ServiceRole()

	var role dto.RoleResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	role, err = sRole.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = role

	return ctx.Serve(err)
}

func (h RoleHandler) Create(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sRole := service.ServiceRole()

	var req dto.RoleRequestCreate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = sRole.Create(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h RoleHandler) Update(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sRole := service.ServiceRole()

	var req dto.RoleRequestUpdate

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

	ctx.ResponseData, err = sRole.Update(ctx.Request().Context(), req, id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h RoleHandler) Archive(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sRole := service.ServiceRole()

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = sRole.Archive(ctx.Request().Context(), id)

	return ctx.Serve(err)
}

func (h RoleHandler) UnArchive(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sRole := service.ServiceRole()

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = sRole.UnArchive(ctx.Request().Context(), id)

	return ctx.Serve(err)
}
