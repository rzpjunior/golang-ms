package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	Option global.HandlerOptions
}

// URLMapping declare endpoint with handler function.
func (h *AuthHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	r.POST("/login", h.Login)
}

func (h AuthHandler) Login(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	timezone := ctx.Request().Header.Get("Timezone")

	var req dto.LoginRequest
	login := service.ServiceAuth()

	if timezone != "" {
		req.Timezone = timezone
	} else {
		req.Timezone = "Asia/Jakarta"
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = login.Login(ctx.Request().Context(), req)

	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}
