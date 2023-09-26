package fridge

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("/list_branch_fridge", h.readBranchFridge, auth.Authorized("filter_rdl"))
	r.GET("/user_fridge", h.readUserFridge, auth.Authorized("usf_crt"))
	r.POST("/user_fridge", h.createUserFridge, auth.Authorized("usf_rdl"))

}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Fridge
	var total int64

	if data, total, e = repository.GetFridges(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// read : function to get requested data based on parameters
func (h *Handler) readUserFridge(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.UserFridge
	var total int64

	if data, total, e = repository.GetUserFridges(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// read : function to get requested data based on parameters
func (h *Handler) readBranchFridge(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	//statusFilter := ctx.QueryParam("status")
	var data []*model.BranchFridgeListQuery
	var total int64

	if data, total, e = repository.GetBranchFridges(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.City
	var total int64

	if data, total, e = repository.GetFilterCities(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

type ProductFridgeBoxQuery struct {
	BoxId       string `orm:"column(box_id);null" json:"box_id,omitempty"`
	Name        string
	TotalPrice  float64
	TotalWeight float64
	ImageUrl    string
	Uom         string
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	var data []*ProductFridgeBoxQuery
	o := orm.NewOrm()
	o.Using("read_only")
	o.Raw("select p.name ,pb.total_price ,pb.total_weight,pi.image_url,u.name as uom from "+
		"box_fridge bf join branch_fridge bf2 "+
		"on bf.fridge_id =bf2.fridge_id "+
		"join product_box pb "+
		"on bf.box_id =pb.box_id "+
		"join product p "+
		"on pb.product_id =p.id "+
		"join product_image pi on pi.product_id=p.id "+
		"join uom u on p.uom_id=u.id "+
		"where bf.status=1 and bf2.fridge_id =?", id).QueryRows(&data)
	ctx.ResponseData = data

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) boxProductDetail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	id := ctx.Param("id")
	var data *ProductFridgeBoxQuery
	o := orm.NewOrm()
	o.Using("read_only")
	e = o.Raw("select pb.box_id,p.name ,pb.total_price ,pb.total_weight,pi.image_url,u.name as uom from "+
		"product_box pb "+
		"join product p "+
		"on pb.product_id =p.id "+
		"join product_image pi on pi.product_id=p.id "+
		"join box b on pb.box_id=b.id "+
		"join uom u on p.uom_id=u.id "+
		"where pb.status=1 and b.code =?", id).QueryRow(&data)
	if e == nil {
		data.BoxId = common.Encrypt(data.BoxId)
		ctx.ResponseData = data
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) box(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Box
	var total int64

	if data, total, e = repository.GetBoxes(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)

}

// detail : function to get detailed data by id
func (h *Handler) productBox(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.BoxItem
	var total int64

	if data, total, e = repository.GetBoxItems(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)

}

// // detail : function to get detailed data by id
// func (h *Handler) branchFridge(c echo.Context) (e error) {
// 	ctx := c.(*cuxs.Context)
// 	rq := ctx.RequestQuery()

// 	var data []*model.BranchFridge
// 	var total int64

// 	if data, total, e = repository.GetBranchFridges(rq); e == nil {
// 		ctx.Data(data, total)
// 	}

// 	return ctx.Serve(e)

// }

// detail : function to get detailed data by id
func (h *Handler) boxFridge(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.BoxFridge
	var total int64

	if data, total, e = repository.GetBoxFridges(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)

}

// detail : function to get detailed data by id
func (h *Handler) boxDetail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	id := ctx.Param("id")
	var data *ProductFridgeBoxQuery
	o := orm.NewOrm()
	o.Using("read_only")
	e = o.Raw("select pb.box_id,p.name ,pb.total_price ,pb.total_weight,pi.image_url,u.name as uom from "+
		"product_box pb "+
		"join product p "+
		"on pb.product_id =p.id "+
		"join product_image pi on pi.product_id=p.id "+
		"join box b on pb.box_id=b.id "+
		"join uom u on p.uom_id=u.id "+
		"where pb.status=1 and b.code =?", id).QueryRow(&data)
	if e == nil {
		data.BoxId = common.Encrypt(data.BoxId)
		ctx.ResponseData = data
	}

	return ctx.Serve(e)
}

// //create : function to create new data based on input
// func (h *Handler) create(c echo.Context) (e error) {
// 	ctx := c.(*cuxs.Context)
// 	var r createRequest

// 	if r.Session, e = auth.UserSession(ctx); e == nil {
// 		if e = ctx.Bind(&r); e == nil {
// 			ctx.ResponseData, e = Save(r)
// 		} else {
// 			//post error
// 			errLog := util.ErrorLog{
// 				ErrorCode:    422,
// 				Name:         r.Session.Staff.Name,
// 				Email:        r.Session.Staff.User.Email,
// 				ErrorMessage: e.Error(),
// 				Function:     "create_box_fridge",
// 			}
// 			util.PostToServiceErrorLog(errLog)
// 		}
// 	}

// 	return ctx.Serve(e)
// }

// //create : function to create new data based on input
// func (h *Handler) createTransaction(c echo.Context) (e error) {
// 	ctx := c.(*cuxs.Context)
// 	var r createRequestTransaction

// 	if r.Session, e = auth.UserSession(ctx); e == nil {
// 		if e = ctx.Bind(&r); e == nil {
// 			ctx.ResponseData, e = SaveTransaction(r)
// 		} else {
// 			//post error
// 			errLog := util.ErrorLog{
// 				ErrorCode:    422,
// 				Name:         r.Session.Staff.Name,
// 				Email:        r.Session.Staff.User.Email,
// 				ErrorMessage: e.Error(),
// 				Function:     "create_box_fridge",
// 			}
// 			util.PostToServiceErrorLog(errLog)
// 		}
// 	}

// 	return ctx.Serve(e)
// }

// create : function to create a new user fridge based on input
func (h *Handler) createUserFridge(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequestUser

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = SaveUser(r)
	return ctx.Serve(e)
}
