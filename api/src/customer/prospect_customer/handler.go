package prospect_customer

import (
	"git.edenfarm.id/cuxs/common"
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
	r.GET("", h.read, auth.Authorized("pro_cst_rdl"))
	r.POST("", h.create)
	r.GET("/:id", h.detail)
	r.GET("/filter", h.readFilter)
	r.PUT("/decline/:id", h.decline, auth.Authorized("pro_cst_dec"))
	r.POST("/decline_type", h.declineType, auth.Authorized("pro_cst_dec"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.ProspectCustomer
	var total int64

	if data, total, e = repository.GetProspectiveCustomers(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

//create : function to create new data based on input
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest

	if e = ctx.Bind(&r); e == nil {
		ctx.ResponseData, e = Save(r)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.ProspectCustomer
	var total int64

	if data, total, e = repository.GetProspectiveCustomers(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetProspectiveCustomer("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// decline : function to decline requested data based on parameters
func (h *Handler) decline(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r declineRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Decline(r)
			}
		}
	}

	return ctx.Serve(e)
}

// declineType : function to get requested data based on parameters
func (h *Handler) declineType(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r requestGetListDeclineType
	o := orm.NewOrm()
	o.Using("read_only")

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	o.Raw(
		"SELECT " +
			"g.value_int, g.value_name " +
			"FROM glossary g " +
			"WHERE g.table = 'prospect_customer' AND g.attribute = 'decline_type'").QueryRows(&r.DeclineType)

	for _, v := range r.DeclineType {
		v.ValueIntEnc = common.Encrypt(v.ValueInt)
	}

	ctx.ResponseData = r.DeclineType
	return ctx.Serve(e)
}
