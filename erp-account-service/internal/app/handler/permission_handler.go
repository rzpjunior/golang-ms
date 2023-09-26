package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type PermissionHandler struct {
	Option global.HandlerOptions
}

// URLMapping declare endpoint with handler function.
func (h *PermissionHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	r.GET("", h.Get, middleware.NewMiddleware().Authorized())
	r.GET("/tree", h.GetTree, middleware.NewMiddleware().Authorized())
	r.GET("/privilege", h.GetPrivilege, middleware.NewMiddleware().Authorized())
}

func (h PermissionHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sPermission := service.ServicePermission()

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

	var permissions []*dto.PermissionResponse
	var total int64
	permissions, total, err = sPermission.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(permissions, total, page)

	return ctx.Serve(err)
}

func (h PermissionHandler) GetTree(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sPermission := service.ServicePermission()

	var permissions []*dto.PermissionResponse
	permissions, err = sPermission.GetTree(ctx.Request().Context())
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = permissions

	return ctx.Serve(err)
}

func (h PermissionHandler) GetPrivilege(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sPermission := service.ServicePermission()

	userID := ctx.Request().Context().Value(constants.KeyUserID).(int64)

	var permissions []string
	permissions, err = sPermission.GetPrivilege(ctx.Request().Context(), userID)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = permissions

	return ctx.Serve(err)
}
