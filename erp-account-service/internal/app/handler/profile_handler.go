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

type ProfileHandler struct {
	Option global.HandlerOptions
}

// URLMapping declare endpoint with handler function.
func (h *ProfileHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	r.GET("", h.Get, middleware.NewMiddleware().Authorized())
	r.PUT("", h.Update, middleware.NewMiddleware().Authorized())
	r.GET("/menu", h.GetMenu, middleware.NewMiddleware().Authorized())
	r.PUT("/reset_password", h.ResetPassword, middleware.NewMiddleware().Authorized())
}

func (h ProfileHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sUser := service.ServiceUser()

	userID := ctx.Request().Context().Value(constants.KeyUserID).(int64)
	ctx.ResponseData, err = sUser.GetByID(ctx.Request().Context(), userID)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h ProfileHandler) GetMenu(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sProfile := service.ServiceProfile()

	userID := ctx.Request().Context().Value(constants.KeyUserID).(int64)
	ctx.ResponseData, err = sProfile.GetMenu(ctx.Request().Context(), userID)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h ProfileHandler) Update(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sUser := service.ServiceUser()

	userID := ctx.Request().Context().Value(constants.KeyUserID).(int64)
	var req dto.ProfileRequestUpdate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = sUser.UpdateProfile(ctx.Request().Context(), req, userID)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h ProfileHandler) ResetPassword(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sUser := service.ServiceUser()

	userID := ctx.Request().Context().Value(constants.KeyUserID).(int64)
	var req dto.UserRequestResetPassword
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	// validate password
	if req.Password != req.PasswordConfirm {
		err = edenlabs.ErrorValidation("password", "The password is not match")
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = sUser.ResetPassword(ctx.Request().Context(), req, userID)

	return ctx.Serve(err)
}
