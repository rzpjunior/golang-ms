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

type UserHandler struct {
	Option global.HandlerOptions
}

// URLMapping declare endpoint with handler function.
func (h *UserHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	r.GET("", h.Get, middleware.NewMiddleware().Authorized())
	r.GET("/:id", h.Detail, middleware.NewMiddleware().Authorized())
	r.POST("", h.Create, middleware.NewMiddleware().Authorized("usr_crt"))
	r.PUT("/:id", h.Update, middleware.NewMiddleware().Authorized("usr_upd"))
	r.PUT("/:id/reset_password", h.ResetPassword, middleware.NewMiddleware().Authorized("usr_rst_pass"))
	r.PUT("/archive/:id", h.Archive, middleware.NewMiddleware().Authorized("usr_arc"))
	r.PUT("/unarchive/:id", h.UnArchive, middleware.NewMiddleware().Authorized("usr_urc"))
}

func (h UserHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sUser := service.ServiceUser()

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	req := &dto.GetUserRequest{
		Offset:     page.Offset,
		Limit:      page.Limit,
		Status:     ctx.GetParamInt("status"),
		Search:     ctx.GetParamString("search"),
		OrderBy:    ctx.GetParamString("order_by"),
		SiteID:     ctx.GetParamString("site_id"),
		RegionID:   ctx.GetParamString("region_id"),
		DivisionID: int64(ctx.GetParamInt("division_id")),
		RoleID:     int64(ctx.GetParamInt("role_id")),
	}

	var users []*dto.UserResponse
	var total int64
	users, total, err = sUser.Get(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(users, total, page)

	return ctx.Serve(err)
}

func (h UserHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sUser := service.ServiceUser()

	var user dto.UserResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	user, err = sUser.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = user

	return ctx.Serve(err)
}

func (h UserHandler) Create(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sUser := service.ServiceUser()

	var req dto.UserRequestCreate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	// validate email
	res, _ := sUser.GetByEmail(ctx.Request().Context(), req.Email)
	if res.ID != 0 {
		err = edenlabs.ErrorValidation("email", "The email already exists")
		return ctx.Serve(err)
	}

	// validate password
	if req.Password != req.PasswordConfirm {
		err = edenlabs.ErrorValidation("password", "The password is not match")
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = sUser.Create(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h UserHandler) Update(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sUser := service.ServiceUser()

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var req dto.UserRequestUpdate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = sUser.Update(ctx.Request().Context(), req, id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h UserHandler) Archive(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sUser := service.ServiceUser()

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = sUser.Archive(ctx.Request().Context(), id)

	return ctx.Serve(err)
}

func (h UserHandler) UnArchive(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sUser := service.ServiceUser()

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = sUser.UnArchive(ctx.Request().Context(), id)

	return ctx.Serve(err)
}

func (h UserHandler) ResetPassword(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sUser := service.ServiceUser()

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

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

	ctx.ResponseData, err = sUser.ResetPassword(ctx.Request().Context(), req, id)

	return ctx.Serve(err)
}
