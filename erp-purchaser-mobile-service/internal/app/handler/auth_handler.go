package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	Option       global.HandlerOptions
	ServicesAuth service.IAuthService
}

// URLMapping declare endpoint with handler function.
func (h *AuthHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesAuth = service.NewAuthService()

	cMiddleware := middleware.NewMiddleware()
	r.POST("/login", h.Login)
	r.GET("/session", h.Session, cMiddleware.Authorized("purchaser_app"))
}

func (h AuthHandler) Login(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.LoginRequest

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	timezone := ctx.Request().Header.Get("Timezone")

	if timezone != "" {
		req.Timezone = timezone
	} else {
		req.Timezone = "Asia/Jakarta"
	}

	ctx.ResponseData, err = h.ServicesAuth.Login(ctx.Request().Context(), req)

	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h AuthHandler) Session(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	ctx.ResponseData, err = h.ServicesAuth.Session(ctx.Request().Context())
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}
