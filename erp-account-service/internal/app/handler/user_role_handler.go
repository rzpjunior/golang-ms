package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type UserRoleHandler struct {
	Option global.HandlerOptions
}

// URLMapping declare endpoint with handler function.
func (h *UserRoleHandler) URLMapping(r *echo.Group) {
	r.GET("/:id", h.UserRoleByUserId, middleware.NewMiddleware().Authorized())
}

func (h UserRoleHandler) UserRoleByUserId(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sUserRole := service.ServiceUserRole()

	var user dto.UserRoleByUserIdResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	user, err = sUserRole.GetByUserID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = user

	return ctx.Serve(err)
}
