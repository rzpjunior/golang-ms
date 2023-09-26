package koli

import (
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.AuthorizedMobile())
	r.GET("/combine", h.readCombineWithDeliveryKoli, auth.AuthorizedMobile())
	r.GET("/id", h.detail, auth.AuthorizedMobile())
	r.POST("", h.create, auth.AuthorizedMobile())
}

// create : function to create new data based on input
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest

	if e = ctx.Bind(&r); e == nil {
		ctx.ResponseData, e = Save(r)
	}

	return ctx.Serve(e)
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Koli
	var total int64

	if data, total, e = repository.GetKolis(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readCombineWithDeliveryKoli : function to get requested data based on parameters
func (h *Handler) readCombineWithDeliveryKoli(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	o := orm.NewOrm()
	o.Using("read_only")

	var data []*model.Koli
	var dataDeliveryKolies []*model.DeliveryKoli
	var total int64

	q := "select * from koli k"
	q1 := "select count(*) from koli k"
	o.Raw(q).QueryRows(&data)
	o.Raw(q1).QueryRow(&total)

	if dataDeliveryKolies, _, e = repository.GetDeliveryKolis(rq); e == nil {
		for _, koli := range data {
			for _, dKoli := range dataDeliveryKolies {
				if koli.ID == dKoli.Koli.ID {
					koli.Quantity = dKoli.Quantity
				}
			}
		}
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetKoli("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}
