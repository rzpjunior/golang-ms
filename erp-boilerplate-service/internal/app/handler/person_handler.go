package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-boilerplate-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-boilerplate-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-boilerplate-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type PersonHandler struct {
	Option         global.HandlerOptions
	ServicesPerson service.IPersonService
}

// URLMapping implements router.RouteHandlers
func (h *PersonHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesPerson = service.NewPersonService()

	r.GET("", h.Index)
	r.GET("/:id", h.Detail)
	r.POST("", h.Create)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
}

func (h PersonHandler) Index(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params search
	search := ctx.QueryParam("search")

	var persons []dto.PersonResponseGet
	var total int64
	persons, total, err = h.ServicesPerson.Get(ctx.Request().Context(), page.Start, page.Limit, search)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(persons, total, page)

	return ctx.Serve(err)
}

func (h PersonHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var person dto.PersonResponseGet

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorInvalid("id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	person, err = h.ServicesPerson.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = person

	return ctx.Serve(err)
}

func (h PersonHandler) Create(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.PersonRequestCreate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesPerson.Create(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h PersonHandler) Update(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.PersonRequestUpdate

	if req.ID, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorInvalid("id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesPerson.Update(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h PersonHandler) Delete(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.PersonRequestDelete

	if req.ID, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorInvalid("id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesPerson.Delete(ctx.Request().Context(), req)

	return ctx.Serve(err)
}
