package farm_stack

import (
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/src/auth"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.POST("", h.read, auth.AuthorizedFarmStack())
}

// read to show detail image and tag
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r requestGet
	o := orm.NewOrm()
	o.Using("read_only")
	if e = ctx.Bind(&r); e == nil {
		o.Raw(
			"SELECT s.name 'farmer_name', s.latitude 'latitude',s.longitude 'longitude',p.`name` 'product_name',poi.unit_price 'product_price',po.eta_date 'purchase_date', po.eta_time 'purchase_time' "+
				"from purchase_order_item poi "+
				"join purchase_order po on po.id = poi.purchase_order_id "+
				"join supplier s on s.id = po.supplier_id "+
				"join product p on p.id = poi.product_id "+
				"where po.`status` not in (3,4) "+
				"and po.supplier_badge_id in (1,8) "+
				"and (s.latitude is not null and s.longitude is not null ) "+
				"and poi.unit_price != 0 "+
				"and  poi.product_id in (40,41,42,43,44,45,46,353,356,400,401,402,409,410,411,412,413,414,458,468,472,486,568,610,761,796,797,824,825,830,831,851,854,856,861,991,1046,1070,1120) "+
				"and po.eta_date = ?", r.EtaDate.Format("2006-01-02")).QueryRows(&r.Data)
		if e == nil {
			ctx.ResponseData = r.Data
		}
	}

	return ctx.Serve(e)
}
