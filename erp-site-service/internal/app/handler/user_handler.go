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

type UserHandler struct {
	Option       global.HandlerOptions
	ServicesUser service.IUserService
}

// URLMapping implements router.RouteHandlers
func (h *UserHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesUser = service.NewServiceUser()

	cMiddleware := middleware.NewMiddleware()

	r.GET("", h.GetUser, cMiddleware.Authorized())
}

func (h UserHandler) GetUser(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	req := dto.GetUserRequest{
		Offset:     page.Offset,
		Limit:      page.Limit,
		Status:     ctx.GetParamInt("status"),
		Search:     ctx.GetParamString("search"),
		OrderBy:    ctx.GetParamString("order_by"),
		SiteId:     int64(ctx.GetParamInt("site_id")),
		DivisionId: int64(ctx.GetParamInt("division_id")),
		RoleId:     int64(ctx.GetParamInt("role_id")),
	}

	users, total, err := h.ServicesUser.GetUser(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.DataList(users, total, page)

	return ctx.Serve(err)
}
