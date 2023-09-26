package distribution_network

import (
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("dis_net_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("dis_net_rdd"))
	r.GET("/:id/top_product", h.readMerchantTopProduct, auth.Authorized("dis_net_rdd"))
	r.GET("/:id/order_performance", h.merchantOrderPerformance, auth.Authorized("dis_net_rdd"))
	r.GET("/:id/payment_performance", h.merchantPaymentPerformance, auth.Authorized("dis_net_rdd"))
	r.GET("/:id/order_payment_details", h.merchantOrderPaymentDetails, auth.Authorized("dis_net_rdd"))
	r.POST("", h.create, auth.Authorized("dis_net_crt"))
	r.PUT("/:id", h.update, auth.Authorized("dis_net_upd"))
	r.PUT("/reset/password/:id", h.resetPassword, auth.Authorized("dis_net_rst_psw"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.Merchant
	var total int64

	if data, total, e = repository.GetMerchantsDistributionNetwork(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if ctx.ResponseData, e = repository.GetMerchant("id", id); e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}

// readMerchantTopProduct : function to get list top sales product by merchant id
func (h *Handler) readMerchantTopProduct(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var (
		id          int64
		fromDateStr = ctx.QueryParam("fromdate")
		toDateStr   = ctx.QueryParam("todate")
		layout      = "2006-01-02"
		loc, _      = time.LoadLocation("Asia/Jakarta")
	)
	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	_, e = time.ParseInLocation(layout, fromDateStr, loc)
	if e != nil {
		return ctx.Serve(e)
	}
	_, e = time.ParseInLocation(layout, toDateStr, loc)
	if e != nil {
		return ctx.Serve(e)
	}

	if ctx.ResponseData, e = repository.GetMerchantTopProduct(fromDateStr, toDateStr, "id", id); e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}

// merchantOrderPerformance : function to get detailed of order performance that merchant data by id
func (h *Handler) merchantOrderPerformance(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var topProductId int64
	topProductIdStr := ctx.QueryParam("top_product_id")
	if topProductId, e = common.Decrypt(topProductIdStr); e != nil {
		return ctx.Serve(e)
	}

	fromDateStr := ctx.QueryParam("fromdate")
	toDateStr := ctx.QueryParam("todate")
	loc, _ := time.LoadLocation("Asia/Jakarta")

	layout := "2006-01-02"
	fromDate, e := time.ParseInLocation(layout, fromDateStr, loc)
	if e != nil {
		return ctx.Serve(e)
	}
	toDate, e := time.ParseInLocation(layout, toDateStr, loc)
	if e != nil {
		return ctx.Serve(e)
	}

	var id int64
	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if ctx.ResponseData, e = repository.GetMerchantOrderPerformance("id", int(topProductId), fromDate, toDate, id); e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}

// merchantPaymentPerformance : function to get detailed of payment performance that merchant data by id
func (h *Handler) merchantPaymentPerformance(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if ctx.ResponseData, e = repository.GetMerchantPaymentPerformance("id", id); e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}

// merchantOrderPaymentDetails : function to get detailed of order payment details that merchant data by id
func (h *Handler) merchantOrderPaymentDetails(c echo.Context) (e error) {
	var (
		data  []*model.SalesInvoice
		total int64
	)

	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var id int64
	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if data, total, e = repository.GetMerchantSalesInvoices(rq, id); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

//create : function to create new data based on input
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Save(r)

	return ctx.Serve(e)
}

// update : function to update data
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Update(r)

	return ctx.Serve(e)
}

// resetPassword : function to reset merchant password
func (h *Handler) resetPassword(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r resetPasswordRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	if e = resetPassword(r); e != nil {
		return ctx.Serve(e)
	}

	return ctx.Serve(e)
}
